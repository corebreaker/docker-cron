package main

import (
	"log"
	"os"
	"sync"
)

var (
	logMtx sync.Mutex
	exit = func() {
		os.Exit(1)
	}
)

func LogFatal(msg string, args ...interface{}) {
	logMtx.Lock()
	defer logMtx.Unlock()

	log.Printf("[FATL] " + msg, args...)

	exit()
}

func LogError(msg string, args ...interface{}) {
	logMtx.Lock()
	defer logMtx.Unlock()

	log.Printf("[ERRO] " + msg, args...)
}

func LogWarn(msg string, args ...interface{}) {
	logMtx.Lock()
	defer logMtx.Unlock()

	log.Printf("[WARN] " + msg, args...)
}

func LogInfo(msg string, args ...interface{}) {
	logMtx.Lock()
	defer logMtx.Unlock()

	log.Printf("[INFO] " + msg, args...)
}
