package lib

import (
	"errors"

	"golang.org/x/sys/unix"
)

// epoll is a Linux kernel system call for a scalable I/O event notification mechanism.
// Its function is to monitor multiple file descriptors to see whether I/O is possible on any of them
// epoll operates in O(1) time !
// epoll uses a red-black tree (RB-tree) data structure 
// to keep track of all file descriptors that are currently being monitored

type fdPoller struct {
	fd   int // File descriptor
	epfd int // Epoll file descriptor
	pipe [2]int //Pipe for waking up
}

func emptyPoller(fd int) *fdPoller {
	poller := new(fdPoller)
	poller.fd = fd
	poller.epfd = -1
	poller.pipe[0] = -1
	poller.pipe[1] = -1
	return poller
}

// Create a new inotify poller.
// This creates an inotify handler, and an epoll handler.
func newFdPoller(fd int) (*fdPoller, error) {
	var errno error
	poller := emptyPoller(fd)
	defer func() {
		if errno != nil {
			poller.close()
		}
	}()
	poller.fd = fd
	
	// Create epoll fd
	poller.epfd, errno = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if poller.epfd == -1 {
		return nil, errno
	}

	// Create pipe; pipe[0] is the read end, pipe[1] is the write end.
	errno = unix.Pipe2(poller.pipe[:], unix.O_NONBLOCK|unix.O_CLOEXEC)
	if errno != nil {
		return nil, errno
	}

	// Register inotify fd with epoll
	event := unix.EpollEvent{
		Fd:     int32(poller.fd),
		Events: unix.EPOLLIN,
	}

	errno = unix.EpollCtl(poller.epfd, unix.EPOLL_CTL_ADD, poller.fd, &event)
	if errno != nil {
		return nil, errno
	}

	// Register pipe fd with epoll
	event = unix.EpollEvent{
		Fd:     int32(poller.pipe[0]),
		Events: unix.EPOLLIN,
	}
	errno = unix.EpollCtl(poller.epfd, unix.EPOLL_CTL_ADD, poller.pipe[0], &event)
	if errno != nil {
		return nil, errno
	}

	return poller, nil
}

// Wait using epoll
// Returns true if sth is ready to be read
// false if there is not
func (poller *fdPoller) wait()(bool, error) {
	// 3 possible events per fd, and 2 fds, makes a maximum of 6 events.
	// I don't know whether epoll_wait returns the number of events returned,
	// or the total number of events ready.
	// I decided to catch both by making the buffer one larger than the maximum.
	events := make([]unix.EpollEvent, 7)
	for {
		n, errno := unix.EpollWait(poller.epfd, events, -1)
		if n == -1 {
			if errno == unix.EINTR {
				continue
			}
			return false, errno
		}
		if n == 0 {
			// If there are no events, try again.
			continue
		}
		if n > 6 {
			// This should never happen. More events were returned than should be possible.
			return false, errors.New("epoll_wait returned more events than I know what to do with")
		}
		ready := events[:n]
		epollhup := false
		epollerr := false
		epollin := false
		for _, event := range ready {
			if event.Fd == int32(poller.fd) {
				if event.Events&unix.EPOLLHUP != 0 {
					// This should not happen, but if it does, treat it as a wakeup.
					epollhup = true
				}
				if event.Events&unix.EPOLLERR != 0 {
					// If an error is waiting on the file descriptor, we should pretend
					// something is ready to read, and let unix.Read pick up the error.
					epollerr = true
				}
				if event.Events&unix.EPOLLIN != 0 {
					// There is data to read.
					epollin = true
				}
			}
			if event.Fd == int32(poller.pipe[0]) {
				if event.Events&unix.EPOLLHUP != 0 {
					// Write pipe descriptor was closed, by us. This means we're closing down the
					// watcher, and we should wake up.
				}
				if event.Events&unix.EPOLLERR != 0 {
					// If an error is waiting on the pipe file descriptor.
					// This is an absolute mystery, and should never ever happen.
					return false, errors.New("Error on the pipe descriptor.")
				}
				if event.Events&unix.EPOLLIN != 0 {
					// This is a regular wakeup, so we have to clear the buffer.
					err := poller.clearWake()
					if err != nil {
						return false, err
					}
				}
			}
		}

		if epollhup || epollerr || epollin {
			return true, nil
		}
		return false, nil
	}
}

// Close the write end of the poller
func (poller *fdPoller) wake() error {
	buf := make([]byte, 1)
	n, errno := unix.Write(poller.pipe[1], buf)
	if n == -1 {
		if errno == unix.EAGAIN {
			// Buffer is full, poller will wake.
			return nil
		}
		return errno
	}
	return nil
}

func (poller *fdPoller) clearWake() error {
	// You have to be woken up a LOT in order to get to 100!
	buf := make([]byte, 100)
	n, errno := unix.Read(poller.pipe[0], buf)
	if n == -1 {
		if errno == unix.EAGAIN {
			// Buffer is empty, someone else cleared our wake.
			return nil
		}
		return errno
	}
	return nil
}

// Close all poller file descriptors, but not the one passed to it.
func (poller *fdPoller) close() {
	if poller.pipe[1] != -1 {
		unix.Close(poller.pipe[1])
	}
	if poller.pipe[0] != -1 {
		unix.Close(poller.pipe[0])
	}
	if poller.epfd != -1 {
		unix.Close(poller.epfd)
	}
}