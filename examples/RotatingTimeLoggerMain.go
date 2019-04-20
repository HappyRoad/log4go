package main

import (
	"log4go"
	"time"
	"os"
)

func main(){
	l, err := log4go.NewRotatingTimeLogger("./log", ".log", log4go.ByDay)
	if err != nil {
		log4go.Error(err.Error())
		os.Exit(1)
	}

	log4go.Std = log4go.NewLogger(l, log4go.DEBUG, "test")

	for {
		log4go.Debug("hello world")
		log4go.Info("hello world")
		time.Sleep(time.Millisecond)
	}
}