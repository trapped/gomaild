package locker

import (
	"errors"
	"log"
	. "sync"
)

var Locks map[string]bool = make(map[string]bool, 0)

var Waits map[string]*Mutex = make(map[string]*Mutex, 0)

func Lock(s string) error {
	if Locks[s] != false {
		return errors.New("already locked")
	} else {
		Locks[s] = true
		log.Println("Locker: Locking `" + s + "`")
		return nil
	}
}

func Unlock(s string) error {
	if Locks[s] != true {
		return errors.New("already unlocked")
	} else {
		Locks[s] = false
		log.Println("Locker: Unlocking `" + s + "`")
		return nil
	}
}

func MLock(s string) {
	if Waits[s] == nil {
		Waits[s] = &Mutex{}
	}
	Waits[s].Lock()
}

func MUnlock(s string) {
	if Waits[s] != nil {
		Waits[s].Unlock()
	}
}
