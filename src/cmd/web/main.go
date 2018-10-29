package main

import (
	"flag"

	. "common/clog"
	"net/http"
	"web/conf"
	"web/router"

	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	f := flag.String("conf", "./conf/web.yaml", "web server's configure file path")
	flag.Parse()
	err := conf.Init(*f)
	if err != nil {
		Blog.Fatalf("Init configure error: %s", err.Error())
		return
	}

	Blog.Infof("Start web server on port %s", conf.C.Port)
	r := router.NewRouter()

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.C.HTMLPath)))

	server := graceful.Server{
		Server: &http.Server{
			Addr:    conf.C.Port,
			Handler: r,
		},
		BeforeShutdown: func() bool {
			Blog.Info("Cleanup self-defined resources ...")
			Flush()
			// server will not shutdown if return false
			return true
		},
	}
	err = server.ListenAndServe()
	if err != nil {
		Blog.Fatalf("Start http server error: %s", err.Error())
	}
}
