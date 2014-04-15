//Package locker provides basic software locking of string values, such as filepaths or usernames.
package locker

import (
	"errors"
	"log"
	. "sync"
)

var Locks map[string]bool = make(map[string]bool, 0) //Contains the active locks.

var Waits map[string]*Mutex = make(map[string]*Mutex, 0) //Contains the various mutex waits.

//Locks a value.
func Lock(s string) error {
	if Locks[s] != false {
		return errors.New("already locked")
	} else {
		Locks[s] = true
		log.Println("Locker: Locking `" + s + "`")
		return nil
	}
}

//Unlocks a value.
func Unlock(s string) error {
	if Locks[s] != true {
		return errors.New("already unlocked")
	} else {
		Locks[s] = false
		log.Println("Locker: Unlocking `" + s + "`")
		return nil
	}
}

//Locks a value. If it's already locked, blocks until it unlocks.
func MLock(s string) {
	if Waits[s] == nil {
		Waits[s] = &Mutex{}
	}
	Waits[s].Lock()
}

//Unlocks a value.
func MUnlock(s string) {
	if Waits[s] != nil {
		Waits[s].Unlock()
	}
}
