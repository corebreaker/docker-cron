package main

import (
	"sync"
	"time"
)

func fire(newTick, oldTick time.Time) {
	LogInfo("Tick at %s, old at %s", newTick, oldTick)
	for _, job := range GetJobs() {
		if job.next.IsZero() {
			job.next = job.expr.Next(oldTick)
		}

		LogInfo("Next at: %s", job.next)
		if (!job.next.IsZero()) && newTick.After(job.next) {
			job.next = job.expr.Next(newTick)
			go job.Run()
		}
	}
}

func StartScheduler() {
	var tickMtx  sync.Mutex
	var lastTick time.Time

	nextTick := func() time.Time {
		tickMtx.Lock()
		defer tickMtx.Unlock()

		old := lastTick
		lastTick = lastTick.Add(time.Minute)

		return old
	}

	now := time.Now()
	start := now.Round(time.Minute)

	dist := start.Sub(time.Now()).Seconds()
	LogInfo("First Tick: %s; %.0fs from now", start, dist)

	if dist < 3 {
		start = start.Add(time.Minute)
		LogInfo("First Tick corrected to %s", start)
	}

	lastTick = start.Add(-1 * time.Minute)

	time.AfterFunc(start.Sub(time.Now()), func() {
		LogInfo("Started at %s", time.Now())
		go fire(time.Now(), nextTick())

		for t := range time.Tick(time.Minute) {
			go fire(t, nextTick())
		}
	})
}
