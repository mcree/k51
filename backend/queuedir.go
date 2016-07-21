package backend

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"time"
	"os"
	"github.com/jehiah/go-strftime"
	"math/rand"
	"fmt"
)

type QueueDirReader struct {
	path            string
	watcher         *fsnotify.Watcher
	items           chan *QueueItem
	handler         ItemHandler
	waitAfterCreate time.Duration
	closed          chan bool
}

type QueueDirWriter struct {
	path string
	namePrefix string
	nameSuffix string
	fileMode os.FileMode
}

type QueueItem struct {
	Name    string
	Payload []byte
	Err     error
}

type ItemHandler func(QueueItem)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewQueueDirReader starts monitoring a given directory for incoming files.
// Item handler is notified on new files.
func NewQueueDirReader(path string, handler ItemHandler) (*QueueDirReader, error) {

	qd := QueueDirReader{
		path:            path,
		items:           make(chan *QueueItem),
		handler:         handler,
		waitAfterCreate: time.Millisecond * 500,
		closed:          make(chan bool),
	}
	var err error

	qd.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// process items
	go func() {
	WatchLoop:
		for {
			select {
			case event := <-qd.watcher.Events:
				//log.Println("event:", event)
				if s, err := os.Stat(event.Name) ; s != nil && err == nil {
					//log.Println("stat:", s)
					if ! s.IsDir() {
						if event.Op & fsnotify.Create == fsnotify.Create {
							go func() {
								time.Sleep(qd.waitAfterCreate)
								//log.Println("incoming file in queue: " + event.Name)
								i := QueueItem{Name: event.Name}
								i.Payload, i.Err = ioutil.ReadFile(event.Name)
								qd.handler(i)
							}()
						}
					}
				}
			case err := <-qd.watcher.Errors:
				log.Println("error:", err)
			case <-qd.closed:
				break WatchLoop
			}
		}
	}()

	// start watching dir
	err = qd.watcher.Add(path)
	if err != nil {
		return nil, err
	}

	// process files that are already present
	var files []os.FileInfo
	files, err = ioutil.ReadDir(qd.path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		// emulate fsnotify event for existing files
		qd.watcher.Events <- fsnotify.Event{Name: qd.path + string(os.PathSeparator) + f.Name(), Op: fsnotify.Create }
	}

	return &qd, err
}

// Close queue dir reader - subsequent attempts to use the dir result in a panic
func (qd *QueueDirReader) Close() error {
	var err error
	err = qd.watcher.Close()
	qd.closed <- true
	return err
}

// NewQueueDirWriter provides utilities for writing files to an outgoing queue directory.
// The filenames are created with the specified name prefix and suffix
func NewQueueDirWriter(path string, namePrefix string, nameSuffix string) (*QueueDirWriter, error) {
	var err error

	qd := QueueDirWriter{
		path: path,
		namePrefix: namePrefix,
		nameSuffix: nameSuffix,
		fileMode: os.FileMode(int(0755)),
	}

	return &qd, err
}

// Write puts the payload into the queue dir with a new unique filename
func (qd *QueueDirWriter) Write(payload []byte) (string, error) {
	hash := strftime.Format("%Y%m%d_%H%M%S", time.Now())
	hash += "_" + fmt.Sprintf("%x_%x", time.Now().UnixNano()&0xffffff, rand.Uint32())
	fileName := qd.path + string(os.PathSeparator) + qd.namePrefix + hash + qd.nameSuffix
	err := ioutil.WriteFile(fileName, payload, qd.fileMode)
	return fileName, err
}

// Close queue dir writer - subsequent attempts to use the dir result in a panic
func (qd *QueueDirWriter) Close() error {
	var err error
	return err
}

