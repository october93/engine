// Code generated by github.com/vektah/dataloaden, DO NOT EDIT.

package dataloader

import (
	"sync"
	"time"

	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/model"
)

// ReactionSliceLoader batches and caches requests
type ReactionSliceLoader struct {
	// this method provides the data for the loader
	fetch func(keys []globalid.ID) ([][]model.Reaction, []error)

	// how long to done before sending a batch
	wait time.Duration

	// this will limit the maximum number of keys to send in one batch, 0 = no limit
	maxBatch int

	// INTERNAL

	// lazily created cache
	cache map[globalid.ID][]model.Reaction

	// the current batch. keys will continue to be collected until timeout is hit,
	// then everything will be sent to the fetch method and out to the listeners
	batch *reactionSliceBatch

	// mutex to prevent races
	mu sync.Mutex
}

type reactionSliceBatch struct {
	keys    []globalid.ID
	data    [][]model.Reaction
	error   []error
	closing bool
	done    chan struct{}
}

// Load a reaction by key, batching and caching will be applied automatically
func (l *ReactionSliceLoader) Load(key globalid.ID) ([]model.Reaction, error) {
	return l.LoadThunk(key)()
}

// LoadThunk returns a function that when called will block waiting for a reaction.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *ReactionSliceLoader) LoadThunk(key globalid.ID) func() ([]model.Reaction, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() ([]model.Reaction, error) {
			return it, nil
		}
	}
	if l.batch == nil {
		l.batch = &reactionSliceBatch{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()

	return func() ([]model.Reaction, error) {
		<-batch.done

		var data []model.Reaction
		if pos < len(batch.data) {
			data = batch.data[pos]
		}

		var err error
		// its convenient to be able to return a single error for everything
		if len(batch.error) == 1 {
			err = batch.error[0]
		} else if batch.error != nil {
			err = batch.error[pos]
		}

		if err == nil {
			l.mu.Lock()
			l.unsafeSet(key, data)
			l.mu.Unlock()
		}

		return data, err
	}
}

// LoadAll fetches many keys at once. It will be broken into appropriate sized
// sub batches depending on how the loader is configured
func (l *ReactionSliceLoader) LoadAll(keys []globalid.ID) ([][]model.Reaction, []error) {
	results := make([]func() ([]model.Reaction, error), len(keys))

	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}

	reactions := make([][]model.Reaction, len(keys))
	errors := make([]error, len(keys))
	for i, thunk := range results {
		reactions[i], errors[i] = thunk()
	}
	return reactions, errors
}

// Prime the cache with the provided key and value. If the key already exists, no change is made
// and false is returned.
// (To forcefully prime the cache, clear the key first with loader.clear(key).prime(key, value).)
func (l *ReactionSliceLoader) Prime(key globalid.ID, value []model.Reaction) bool {
	l.mu.Lock()
	var found bool
	if _, found = l.cache[key]; !found {
		l.unsafeSet(key, value)
	}
	l.mu.Unlock()
	return !found
}

// Clear the value at key from the cache, if it exists
func (l *ReactionSliceLoader) Clear(key globalid.ID) {
	l.mu.Lock()
	delete(l.cache, key)
	l.mu.Unlock()
}

func (l *ReactionSliceLoader) unsafeSet(key globalid.ID, value []model.Reaction) {
	if l.cache == nil {
		l.cache = map[globalid.ID][]model.Reaction{}
	}
	l.cache[key] = value
}

// keyIndex will return the location of the key in the batch, if its not found
// it will add the key to the batch
func (b *reactionSliceBatch) keyIndex(l *ReactionSliceLoader, key globalid.ID) int {
	for i, existingKey := range b.keys {
		if key == existingKey {
			return i
		}
	}

	pos := len(b.keys)
	b.keys = append(b.keys, key)
	if pos == 0 {
		go b.startTimer(l)
	}

	if l.maxBatch != 0 && pos >= l.maxBatch-1 {
		if !b.closing {
			b.closing = true
			l.batch = nil
			go b.end(l)
		}
	}

	return pos
}

func (b *reactionSliceBatch) startTimer(l *ReactionSliceLoader) {
	time.Sleep(l.wait)
	l.mu.Lock()

	// we must have hit a batch limit and are already finalizing this batch
	if b.closing {
		l.mu.Unlock()
		return
	}

	l.batch = nil
	l.mu.Unlock()

	b.end(l)
}

func (b *reactionSliceBatch) end(l *ReactionSliceLoader) {
	b.data, b.error = l.fetch(b.keys)
	close(b.done)
}
