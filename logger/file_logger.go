package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

type FileLogger struct {
	Level       LogLevel
	filePath    string //文件路径
	fileName    string //文件名
	errFileName string //错误日志文件名
	maxFileSize int64
	fileObj     *os.File
	errFileObj  *os.File
}

func NewFileLogger(levelStr, fp, fn string, maxSize int64) *FileLogger {
	logLevel, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	fl := &FileLogger{
		Level:       logLevel,
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
		// mutex:       sync.Mutex{},
	}
	fl.initFile()
	return fl
}

func (fl *FileLogger) initFile() error {
	fullPathName := path.Join(fl.filePath, fl.fileName)
	fileObj, err := os.OpenFile(fullPathName, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("open file failed,err:%v\n", err)
		return err
	}
	errfileObj, err := os.OpenFile("error."+fullPathName, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("open file failed,err:%v\n", err)
		return err
	}
	fl.fileObj = fileObj
	fl.errFileObj = errfileObj
	return nil

}

func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed,err:%v\n", err)
		return false
	}
	return fileInfo.Size() >= f.maxFileSize
}

// 切割文件
func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {
	//切割文件
	nowStr := time.Now().Format("20060102150405000")
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed,err:%v\n", err)
		return nil, err
	}
	logName := path.Join(f.filePath, fileInfo.Name())
	newLogName := fmt.Sprintf("%s.%s", logName, nowStr)
	file.Close()
	err = os.Rename(logName, newLogName)
	if err != nil {
		return nil, err
	}
	//打开一个新的日志文件
	fileObj, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open new log file failed, err:%v\n", err)
		return nil, err
	}
	return fileObj, nil

}

func (log *FileLogger) logPrint(format, level string, args ...interface{}) {
	l, err := parseLogLevel(level)
	if err != nil {
		panic(err)
	}
	if log.Level > l {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05.000")
	funcName, fileName, lineNo := getInfo(3)
	//检查文件大小是否超过最大文件大小
	if log.checkSize(log.fileObj) {
		//打开一个新的日志文件
		fileObj, err := log.splitFile(log.fileObj)
		if err != nil {
			panic(err)
		}
		log.fileObj = fileObj
	}

	s := fmt.Sprintf(format, args...)
	fmt.Fprintf(log.fileObj, "[%s] [%s] [%s:%s:%d] %s\n", now, level, fileName, funcName, lineNo, s)
	if l >= ERROR {
		//检查文件大小是否超过最大文件大小
		if log.checkSize(log.errFileObj) {
			//打开一个新的日志文件
			erorFileObj, err := log.splitFile(log.errFileObj)
			if err != nil {
				panic(err)
			}
			log.errFileObj = erorFileObj
		}
		fmt.Fprintf(log.errFileObj, "[%s] [%s] [%s:%s:%d] %s\n", now, level, fileName, funcName, lineNo, s)
	}
}

func (log *FileLogger) Trace(format string, args ...interface{}) {
	log.logPrint(format, "TRACE", args...)
}

func (log *FileLogger) Debug(format string, args ...interface{}) {
	log.logPrint(format, "DEBUG", args...)
}

func (log *FileLogger) Info(format string, args ...interface{}) {
	log.logPrint(format, "INFO", args...)
}

func (log *FileLogger) Warn(format string, args ...interface{}) {
	log.logPrint(format, "WARN", args...)
}

func (log *FileLogger) Error(format string, args ...interface{}) {
	log.logPrint(format, "ERROR", args...)
}

func (log *FileLogger) Fatal(format string, args ...interface{}) {
	log.logPrint(format, "FATAL", args...)
}

func (log *FileLogger) Close() {
	log.fileObj.Close()
	log.errFileObj.Close()
}
