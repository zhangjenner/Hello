package utils

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"runtime"
	"strings"
)

//=============================================================================
var skipLevel int = 2	//忽略的调用层次

//新建错误信息
func Error(err error) error {
	if strings.HasPrefix(err.Error(), "[ERROR]:") {
		return err
	}
	pc := make([]uintptr, 20)
	n := runtime.Callers(skipLevel, pc)
	errmsg := "[ERROR]:" + err.Error() + "\n"
	frames := runtime.CallersFrames(pc[:n-2])
	for {
		frame, more := frames.Next()
		errmsg += fmt.Sprintf("%s[%d]:%s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}
	skipLevel = 2
	return errors.New(errmsg)
}

//新建错误信息
func Errorf(format string, args ...interface{}) error {
	skipLevel = 3
	if len(args) > 0 {
		return Error(errors.Errorf(format, args))
	} else {
		return Error(errors.New(format))
	}
}

//新建错误回应
func ErrorRsp(err error) peer.Response {
	skipLevel = 3
	return shim.Error(Error(err).Error())
}

//新建错误回应
func ErrorRspf(format string, args ...interface{}) peer.Response {
	skipLevel = 3
	var err error
	if len(args) > 0 {
		err = errors.Errorf(format, args)
	} else {
		err = errors.New(format)
	}
	return shim.Error(Error(err).Error())
}

