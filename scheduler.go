package main

import (
	"sync"
	"time"
)

var GetLastTick func() time.Time

func StartScheduler() {
	var tickMtx  sync.Mutex
	var lastTick time.Time

	setLastTick := func(t *time.Time) time.Time {
		tickMtx.Lock()
		defer tickMtx.Unlock()

		var old = lastTick

		lastTick = *t

		return old
	}

	GetLastTick = func() time.Time {
		tickMtx.Lock()
		defer tickMtx.Unlock()

		return lastTick
	}

	now := time.Now()
	start := now.Round(time.Minute)
	if start.Before(now) {
		start = start.Add(time.Minute)
	}

	if start.Sub(time.Now()).Seconds() < 5 {
		start = start.Add(time.Minute)
	}

	setLastTick(&start)

	time.AfterFunc(start.Sub(time.Now()), func() {
		for t := range time.Tick(time.Minute) {
			old := setLastTick(&t)
			for _, job := range GetJobs() {
				if job.next.IsZero() {
					job.next = job.expr.Next(old)
				}

				if (!job.next.IsZero()) && t.After(job.next) {
					job.next = job.expr.Next(t)
					go job.Run()
				}
			}
		}
	})
}
