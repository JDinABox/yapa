package ilog

import (
	"runtime/debug"

	"github.com/golang/glog"
)

func Error(err error) {
	glog.Errorf("%v\n%v\n", err, string(debug.Stack()))
}

func Exit(err error) {
	glog.Exitf("%v\n%v\n", err, string(debug.Stack()))
}
