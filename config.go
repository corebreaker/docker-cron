package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	gerr "github.com/corebreaker/goerrors"
	"gopkg.in/yaml.v3"
)

var (
	GetJobs func() []*Job
	setJobs func(jobs []*Job)
)

func init() {
	var mtx sync.Mutex
	var jobs []*Job

	GetJobs = func() []*Job {
		mtx.Lock()
		defer mtx.Unlock()

		return jobs
	}

	setJobs = func(j []*Job) {
		mtx.Lock()
		defer mtx.Unlock()

		jobs = j
	}
}

func readConfig(configPath string) ([]*Job, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, gerr.DecorateError(err)
	}

	defer func() { _ = f.Close() }()

	var jobs []*Job

	ext := filepath.Ext(configPath)
	switch ext {
	case ".json":
		if err := json.NewDecoder(f).Decode(&jobs); err != nil {
			return nil, gerr.DecorateError(err)
		}

	case ".yml", ".yaml":
		if err := yaml.NewDecoder(f).Decode(&jobs); err != nil {
			return nil, gerr.DecorateError(err)
		}
	}

	return jobs, nil
}

func UpdateConfig(configPath string) error {
	LogInfo("Updating config")

	configJobs, err := readConfig(configPath)
	if err != nil {
		return err
	}

	var jobs []*Job

	for i, job := range configJobs {
		expr, err := ParseSpec(job.At)
		if err != nil {
			LogWarn("Bad cron expression for the entry %d: %s", i, err)

			continue
		}

		jobs = append(jobs, &Job{
			At:        job.At,
			Container: job.Container,
			Command:   job.Command,
			expr:      expr,
		})
	}

	setJobs(jobs)

	LogInfo("Config updated")

	return nil
}
