package main

import "runtime"

func main() {
	StartScheduler()

	if err := InitAutoUpdate(); err != nil {
		LogError("%s", err)

		return
	}

	runtime.Goexit()
}
