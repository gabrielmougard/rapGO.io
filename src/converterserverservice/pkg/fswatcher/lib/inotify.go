package lib
 
// The inotify mechanism allows you to position 
// a watch descriptor on a file, which will send notifications 
// to the system when events affect the tracked file. 
// As a reminder, in the UNIX world, a file can be a simple file, 
// a directory, a device, a link, etc.

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct {
	Events   chan Event
	Errors   chan error
	mu       sync.Mutex //Map access
	fd       int //file descriptor
	poller   *fdPoller
	watches  map[string]*watch // Map of notifier watches (key: path)
	paths    map[int]string // Map of watched paths (key: watch descriptor)
	done     chan struct{} // Channel for sending a "quit message" to the reader goroutine
	doneResp chan struct{} // Channel to respond to Close
}

// NewWatcher established a new watcher with the underlying OS(UNIX) and begins waiting for events
func NewWacther() (*Watcher, error) {
	fd, errno := unix.InotifyInit1(unix.IN_CLOEXEC)
	if fd == -1 {
		return nil, errno
	}

	//create epoll
	poller, err := newFdPoller(fd)
	if err != nil {
		unix.Close(fd)
		return nil, err
	}
	w := &Watcher{
		fd:       fd,
		poller:   poller,
		watches:  make(map[string]*watch),
		paths:    make(map[int]string),
		Events:   make(chan Event),
		Errors:   make(chan error),
		done:     make(chan struct{}),
		doneResp: make(chan struct{}),
	}

	go w.readEvents()
	return w, nil
}

func (w *Watcher) isClosed() bool {
	select {
	case <-w.done:
		return true
	default:
		return false
	}
}

// Close removes all watchers and closes the events channel
func (w *Watcher) Close() error {
	if w.isClosed() {
		return nil
	}

	// Send 'close'  signal to goroutine and set the Watcher to closed
	close(w.done)

	// Wake up goroutine
	w.poller.wake()

	// Wait for goroutine to close
	<-w.doneResp

	return nil
}

// Add starts watching the named file or directory (not recursive)
func (w *Watcher) Add(name string) error {
	name = filepath.Clean(name)
	if w.isClosed() {
		return errors.New("inotify instance already closed")
	}

	const agnosticEvents = unix.IN_MOVED_TO | unix.IN_MOVED_FROM |
	unix.IN_CREATE | unix.IN_ATTRIB | unix.IN_MODIFY |
	unix.IN_MOVE_SELF | unix.IN_DELETE | unix.IN_DELETE_SELF

	var flags uint32 = agnosticEvents

	w.mu.Lock()
	defer w.mu.Unlock()
	watchEntry := w.watches[name]
	if watchEntry != nil {
		flags |= watchEntry.flags | unix.IN_MASK_ADD
	}
	wd, errno := unix.InotifyAddWatch(w.fd, name, flags)
	if wd == -1 {
		return errno
	}

	if watchEntry == nil {
		w.watches[name] = &watch{wd: uint32(wd), flags: flags}
		w.paths[wd] = name
	} else {
		watchEntry.wd = uint32(wd)
		watchEntry.flags = flags
	}

	return nil
}

// Remove stops watching the named file or directory (not recursive)
func (w *Watcher) Remove(name string) error {
	
}