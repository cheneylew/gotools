package tool

import (
	"runtime"
	"path"
	"fmt"
	"github.com/pkg/errors"
)

func SprintfLine(format string, a ...interface{}) string {
	_, file, line, _ := runtime.Caller(1)
	_, file2, line2, _ := runtime.Caller(2)
	fileName := path.Base(file)
	fileName2 := path.Base(file2)
	callStack := fmt.Sprintf("[%s:%d=>%s:%d]", fileName2, line2, fileName, line)
	return fmt.Sprintf(callStack+format, a...)
}


func ErrorLine(err error) error {
	return errors.New(SprintfLine("%v", err))
}

func ErrorLineWithMSG(msg string, err error) error {
	return errors.New(SprintfLine(msg+":%v", err))
}
