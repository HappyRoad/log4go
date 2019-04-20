package log4go

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type RollingTimeType = int

const (
	ByDay RollingTimeType = iota
	ByHour
	ByMinute
)

type RotatingTimeLogger struct {
	dir       string
	suffix    string
	mu        *sync.RWMutex
	format    string
	file      *os.File
	bufWriter *bufio.Writer
}

func NewRotatingTimeLogger(dir, suffix string, timeType RollingTimeType) (logger *RotatingTimeLogger, err error) {
	logger = &RotatingTimeLogger{
		dir:    dir,
		suffix: suffix,
		mu:     &sync.RWMutex{},
	}
	fi, err := os.Stat(dir)
	if err != nil || os.IsNotExist(err) || !fi.IsDir() {
		err = errors.New(fmt.Sprintf("not exist directory: %s", dir))
	}
	//默认按天切分
	if timeType == ByMinute {
		logger.format = "200601021504"
	} else if timeType == ByHour {
		logger.format = "2006010215"
	} else {
		logger.format = "20060102"
	}
	filename := logger.getFileName(time.Now().Format(logger.format))
	logger.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	logger.bufWriter = bufio.NewWriter(logger.file)

	go logger.monitor()

	return
}

//要求并发安全，并自动切割文件
func (self *RotatingTimeLogger) Printf(format string, v ...interface{}) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.bufWriter.WriteString(fmt.Sprintf(format, v...))
}

func (self *RotatingTimeLogger) rotate(timeflag string) (err error) {
	self.mu.Lock()
	defer self.mu.Unlock()
	//关闭之前的文件
	self.file.Close()
	//重新打开文件
	filename := self.getFileName(timeflag)
	self.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		self.file = os.Stderr
		fmt.Errorf("open file fail: %s, err=%s", filename, err)
	}
	//构造buf writer
	self.bufWriter = bufio.NewWriter(self.file)

	return
}

func (self *RotatingTimeLogger) flush() {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.bufWriter.Flush()
}

func (self *RotatingTimeLogger) getFileName(timeflag string) string {
	return self.dir + "/" + timeflag + self.suffix
}

func (self *RotatingTimeLogger) monitor() {
	timerSecond := time.NewTicker(time.Second)
	timerMinute := time.NewTicker(time.Minute)

	lastTime := time.Now().Format(self.format)
	curTime := lastTime
	for {
		select {
		case <-timerSecond.C:
			//检查是否需要切割文件
			curTime = time.Now().Format(self.format)
			if curTime != lastTime {
				lastTime = curTime
				//开始切割文件
				self.rotate(lastTime)
			}
		case <-timerMinute.C:
			//保证每分钟都会flush，避免部分缓存日志丢失
			self.flush()
		}
	}
}
