package main

import (
	"log"
	"os"
	"sync"

	gerr "github.com/corebreaker/goerrors"
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

	log.Printf("[FATAL] " + msg, args...)

	exit()
}

func LogError(msg string, args ...interface{}) {
	logMtx.Lock()
	defer logMtx.Unlock()

	log.Printf("[ERROR] " + msg, args...)
}

func LogWarn(msg string, args ...interface{}) {
	logMtx.Lock()
	defer logMtx.Unlock()

	log.Printf("[WARN]  " + msg, args...)
}

func LogInfo(msg string, args ...interface{}) {
	logMtx.Lock()
	defer logMtx.Unlock()

	log.Printf("[INFO]  " + msg, args...)
}

func LogDebug(msg string, args ...interface{}) {
	if !gerr.GetDebug() {
		return
	}

	logMtx.Lock()
	defer logMtx.Unlock()

	log.Printf("[DEBUG] " + msg, args...)
}
