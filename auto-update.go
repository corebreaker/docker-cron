package main

import (
	"os"
	"path/filepath"

	gerr "github.com/corebreaker/goerrors"
	"github.com/siddontang/go/ioutil2"
	"github.com/tywkeene/go-fsevents"
)

type tFileEventHandler string

func (h tFileEventHandler) Handle(w *fsevents.Watcher, event *fsevents.FsEvent) error {
	return UpdateConfig(event.Path)
}

func (h tFileEventHandler) Check(event *fsevents.FsEvent) bool {
	return event.Path == string(h)
}

func (h tFileEventHandler) GetMask() uint32 {
	return fsevents.FileChangedEvent | fsevents.FileCreatedEvent
}

type tFallbackHandler struct {}

func (tFallbackHandler) Handle(w *fsevents.Watcher, event *fsevents.FsEvent) error { return nil }
func (tFallbackHandler) Check(event *fsevents.FsEvent) bool                        { return true }
func (tFallbackHandler) GetMask() uint32                                           { return fsevents.AllEvents }

func InitAutoUpdate() error {
	confPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		return gerr.DecorateError(err)
	}

	ext := filepath.Ext(confPath)
	switch ext {
	case ".json", ".yml", ".yaml":
	default:
		return gerr.MakeError("The config file must be in JSON or YAML format with the good file extension")
	}

	if ioutil2.FileExists(confPath) {
		if err := UpdateConfig(confPath); err != nil {
			return err
		}
	}

	mask := fsevents.FileChangedEvent | fsevents.FileCreatedEvent

	confDir := filepath.Dir(confPath)
	if !ioutil2.FileExists(confDir) {
		mask |= fsevents.DirCreatedEvent
	}

	w, err := fsevents.NewWatcher(confDir, mask)
	if err != nil {
		return gerr.DecorateError(err)
	}

	if err := w.RegisterEventHandler(tFileEventHandler(confPath)); err != nil {
		return gerr.DecorateError(err)
	}

	if err := w.RegisterEventHandler(tFallbackHandler{}); err != nil {
		return gerr.DecorateError(err)
	}

	go func() {
		for err := range w.Errors {
			LogError("%s", gerr.DecorateError(err))
		}
	}()

	if err := w.StartAll(); err != nil {
		func() {
			defer gerr.DiscardPanic()

			close(w.Errors)
		}()

		return gerr.DecorateError(err)
	}

	go w.WatchAndHandle()

	return nil
}
