package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"bnzq.com/common/panic-catch/cpanic"
)

var (
	gPainicFile   *os.File
	gPainicLogger *log.Logger
)

func main() {

	var (
		err error
	)

	lNow := time.Now().Unix()

	sPFileName := fmt.Sprintf("%v_panic.log", lNow)

	// 将 stderr 重定向到 f
	gPainicFile, gPainicLogger, err = cpanic.NewPanicFile(sPFileName)
	if nil != err {
		log.Printf("cpanic.NewPanicFile %v, ", time.Now())
		return
	}

	var (
		chExit chan int = make(chan int)
	)

	go procPan()

	select {
	case <-chExit:
		return
	}
}
