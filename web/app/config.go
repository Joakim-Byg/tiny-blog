package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

type ConfigScope struct {
	configFilePath string
}

func (cs ConfigScope) registerEventHandler(op fsnotify.Op, f func()) {

}

type ConfigMap struct {
	valueMap   map[string]any
	sourceFile string
}

func getConfigMap(configFilePath string) *ConfigMap {
	cfg := new(ConfigMap)
	cfg.sourceFile = configFilePath

	return cfg
}

type FileWatcher struct {
	files   []string
	watcher fsnotify.Watcher
}

type FileHandler struct {
	fileName      string
	eventHandlers map[fsnotify.Op]map[string]func()
}

func registerHandler(fh *FileHandler, operation fsnotify.Op, handlerName string, handler func()) {
	if fh.eventHandlers == nil {
		fh.eventHandlers = make(map[fsnotify.Op]map[string]func())
	} else if fh.eventHandlers[operation] == nil {
		fh.eventHandlers[operation] = make(map[string]func())
	}
	fh.eventHandlers[operation][handlerName] = handler
}
func initFileWatcher() *FileHandler {
	fileHandler := new(FileHandler)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	return fileHandler
}

func watchEvent(watcher fsnotify.Watcher) {
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
}
