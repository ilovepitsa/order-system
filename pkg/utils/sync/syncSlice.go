package utils

import (
	"errors"
	"sync"
)

var (
	ErrEmpty = errors.New("empty slice")
)

type SyncSlice[T any] struct {
	storage []T
	m       *sync.RWMutex
}

func NewSyncSlice[T any](len, cap int) *SyncSlice[T] {
	return &SyncSlice[T]{
		storage: make([]T, len, cap),
		m:       &sync.RWMutex{},
	}
}

func (sc *SyncSlice[T]) Push(value T) {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.storage = append(sc.storage, value)
}

func (sc *SyncSlice[T]) Top() (value T, err error) {
	sc.m.RLock()
	defer sc.m.RUnlock()
	if len(sc.storage) > 0 {
		return sc.storage[0], nil
	}
	return value, ErrEmpty
}

func (sc *SyncSlice[T]) Pop() error {
	sc.m.Lock()
	defer sc.m.Unlock()
	if len(sc.storage) > 0 {
		sc.storage = sc.storage[:0]
		return nil
	}
	return ErrEmpty
}

func (sc *SyncSlice[T]) Set(i int, val T) {
	sc.m.Lock()
	defer sc.m.Unlock()
	if i >= len(sc.storage) {
		return
	}
	sc.storage[i] = val
}

func (sc *SyncSlice[T]) Get(i int) (val T) {
	sc.m.RLock()
	defer sc.m.RUnlock()
	if i >= len(sc.storage) {
		return
	}
	return sc.storage[i]
}

func (sc *SyncSlice[T]) Len() int {
	sc.m.Lock()
	defer sc.m.Unlock()
	return len(sc.storage)
}
