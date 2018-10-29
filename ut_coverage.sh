#!/usr/bin/env bash
set -e
WORK_DIR=$(cd `dirname $0`; pwd)

# root path
export ROOTPATH=`pwd`
# go path
export GOPATH=$ROOTPATH:$ROOTPATH/vendor:$GOPATH
# bin path
export PATH=$WORK_DIR/bin:$PATH

workdir=.cover
profile="$workdir/cover.out"
mode=count

function generate_cover_data() {
    for pkg in "$@"; do
        # filter test package
        if [[ $pkg =~ "test/3rd_test" || $pkg =~ "test/benchmark" || $pkg =~ "test/module_test" || $pkg =~ "test/resource_test" || $pkg =~ "test/smoke_test" || $pkg =~ "test/system_test" || $pkg =~ "tutorial" ]]; then
            continue
        fi
        f="$workdir/$(echo $pkg | tr / -).cover"
        go test -covermode="$mode" -coverprofile="$f" "./src/$pkg"
    done

    echo "mode: $mode" > "$profile"
    gocovmerge $(ls -d $workdir/*) > $profile
}

function prepare_tool() {
    mkdir -p $WORK_DIR/bin
    cd $WORK_DIR/bin
    go build github.com/axw/gocov/gocov
    go build github.com/wadey/gocovmerge
    go build github.com/AlekSi/gocov-xml
    cd $WORK_DIR
}

function main() {
    # 1. generate cover data
    case $1 in
        part)
            prepare_tool
            generate_cover_data $2
        ;;
        *)
            # default is for jfactory work
            generate_cover_data $(go list ./...)
            gocov convert $profile | gocov-xml > cover_cobertura_result
            sed -i "s#/ws/cmw_jmiss-jvessel_all_coverage/src/#cmw_jmiss-jvessel_all_coverage/src/#g" cover_cobertura_result
            ;;
    esac

    # 2. final result
    echo "======== final result ========"
    go tool cover -func=$profile
}

# prepare work dir
rm -rf "$workdir"
mkdir "$workdir"

# run
main $@
