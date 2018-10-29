package clog

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"time"
)

/*
	log dir
	log file
	log buff
*/

// MaxSize is the maximum size of a log file in bytes.
var MaxSize uint64 = 1024 * 1024 * 500 // 500 M

// MaxFileCount is the maximum of log file with same tag.
var MaxFileCount = 5

var logDir = os.TempDir()

func SetLogDir(dir string) {
	os.MkdirAll(dir, os.ModeDir)
	logDir = dir
}

func getLogDir() string {
	return logDir
}

// log file name example:
//		 jvessel-worker.ERROR.2018.01.02-18.12.12
func logName(tag string, t time.Time) (name, link string) {
	name = fmt.Sprintf("%s.%s.%04d.%02d.%02d-%02d.%02d.%02d",
		program,
		tag,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second())
	return name, program + "." + tag
}

func isLogName(program, tag, filename string) bool {
	pattern := fmt.Sprintf("%s\\.%s\\.[0-9]{4}\\.[0-9]{2}\\.[0-9]{2}-[0-9]{2}\\.[0-9]{2}\\.[0-9]{2}", program, tag)
	matched, err := regexp.MatchString(pattern, filename)
	if err != nil {
		return false
	}
	return matched

}

func getTimestamp(program, tag, filename string) (time.Time, error) {
	if !isLogName(program, tag, filename) {
		return time.Time{}, errors.New("this is not a log file name. " + filename)
	}
	layout := "2006.01.02-15.04.05"
	return time.ParseInLocation(layout, filename[len(filename)-len(layout):], time.Local)

}

// create creates a new log file and returns the file and its filename, which
// contains tag ("INFO", "FATAL", etc.) and t.  If the file is created
// successfully, create also attempts to update the symlink for that tag, ignoring
// errors.
func create(tag string, t time.Time) (f *os.File, filename string, err error) {
	name, link := logName(tag, t)
	dir := logDir
	filename = filepath.Join(dir, name)

	f, err = os.Create(filename)
	if err == nil {
		symlink := filepath.Join(dir, link)
		os.Remove(symlink)        // ignore err
		os.Symlink(name, symlink) // ignore err
		return f, filename, nil
	} else {
		return nil, "", fmt.Errorf("log: cannot create log: %v", err)
	}
}

func rotate(program, tag, dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stdout, "list all log file error: %s", err.Error())
		return err
	}

	type meta struct {
		name string
		ts   time.Time
	}
	metas := make([]*meta, 0)

	for _, v := range files {
		if isLogName(program, tag, v.Name()) {
			ts, err := getTimestamp(program, tag, v.Name())
			if err != nil {
				fmt.Fprintf(os.Stdout, "get file[%s] timestamp error: %s", v.Name(), err.Error())
				continue
			}
			metas = append(metas, &meta{name: v.Name(), ts: ts})
		}
	}

	if len(metas) <= MaxFileCount {
		return nil
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].ts.After(metas[j].ts)
	})
	for _, m := range metas {
		fmt.Println(m.name)
	}
	for i := MaxFileCount; i < len(metas); i++ {
		fmt.Println(time.Now(), "rotate log file", metas[i].name)
		os.Remove(path.Join(dir, metas[i].name))
	}
	return nil

}

// flushSyncWriter is the interface satisfied by logging destinations.
type flushSyncWriter interface {
	Flush() error
	Sync() error
	io.Writer
	Check() error
}

// bufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
const bufferSize = 256 * 1024

// syncBuffer joins a bufio.Writer to its underlying file, providing access to the
// file's Sync method and providing a wrapper for the Write method that provides log
// file rotation. There are conflicting methods, so the file cannot be embedded.
// l.mu is held for all its methods.
type syncBuffer struct {
	logger *loggingT
	*bufio.Writer
	file   *os.File
	sev    severity
	nbytes uint64 // The number of bytes written to this file
}

func (sb *syncBuffer) Sync() error {
	return sb.file.Sync()
}

func (sb *syncBuffer) Check() error {
	files, err := ioutil.ReadDir(getLogDir())
	if err != nil {
		return err
	}
	for _, v := range files {
		// log file exist
		if v.Name() == filepath.Base(sb.file.Name()) {
			return nil
		}
	}
	// switch log file
	err = sb.rotateFile(time.Now())
	if err != nil {
		return err
	}
	return err
}

func (sb *syncBuffer) Write(p []byte) (n int, err error) {
	if sb.nbytes+uint64(len(p)) >= MaxSize {
		if err := sb.rotateFile(time.Now()); err != nil {
			sb.logger.exit(err)
		}
	}
	n, err = sb.Writer.Write(p)
	sb.nbytes += uint64(n)
	if err != nil {
		sb.logger.exit(err)
	}
	return
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (sb *syncBuffer) rotateFile(now time.Time) error {
	if sb.file != nil {
		sb.Flush()
		sb.file.Close()
	}
	var err error
	sb.file, _, err = create(severityName[sb.sev], now)
	sb.nbytes = 0
	if err != nil {
		return err
	}

	sb.Writer = bufio.NewWriterSize(sb.file, bufferSize)

	if sb.sev == debugLog {
		// Write header.
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Log file created at: %s\n", now.Format("2006/01/02 15:04:05"))
		fmt.Fprintf(&buf, "Running on machine: %s\n", host)
		fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		fmt.Fprint(&buf, "Log line format: [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg\n")
		n, err1 := sb.file.Write(buf.Bytes())
		err = err1
		sb.nbytes += uint64(n)
	}
	// try to delete some old log file,
	go rotate(program, severityName[sb.sev], getLogDir())

	return err
}
