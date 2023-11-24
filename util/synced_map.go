package util

import (
	"log"
	"sync"
)

type SyncedMap[K comparable, V any] struct {
	ObservableMap[K, V]
	register   func(m ObservableMap[K, V])
	unregister func(m ObservableMap[K, V])
	registered bool
	mutex      sync.Mutex
}

func NewSyncedMap[K comparable, V any](register func(m ObservableMap[K, V]), unregister func(m ObservableMap[K, V])) *SyncedMap[K, V] {
	log.Printf("new synced map")

	m := NewObservableMap[K, V]()
	sm := &SyncedMap[K, V]{
		ObservableMap: m,
		register:      register,
		unregister:    unregister,
		mutex:         sync.Mutex{},
	}

	m.ObserverCount().Subscribe(func(i int) {
		log.Printf("observer count: %d", i)
		sm.mutex.Lock()
		defer sm.mutex.Unlock()
		if i > 0 {
			if !sm.registered {
				log.Printf("registering synced map")
				sm.register(m)
				sm.registered = true
			}
		} else {
			if sm.registered {
				sm.unregister(m)
				sm.registered = false
			}
		}
	})

	return sm
}
