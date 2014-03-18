package locker

import (
	"errors"
	"log"
)

var Locks map[string]bool = make(map[string]bool, 0)

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
