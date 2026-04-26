// Copyright 2025 Byterio
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memkey

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	updaterOnce sync.Once
	timestamp   uint32
)

type memkey struct {
	cfg  Config
	mu   sync.RWMutex
	db   map[string]entry
	done chan struct{}
}

type entry struct {
	data   []byte
	expiry uint32 // max value is 4294967295 -> Sun Feb 07 2106 06:28:15 GMT+0000
}

// New creates a new memkey.
func New(config ...Config) *memkey {
	cfg := configDefault(config...)
	mk := &memkey{
		db:   make(map[string]entry),
		cfg:  cfg,
		done: make(chan struct{}),
	}
	startTimeStampUpdater()
	go mk.startCleanupLoop()
	return mk
}

// Get value by key.
func (mk *memkey) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	mk.mu.RLock()
	v, ok := mk.db[key]
	mk.mu.RUnlock()
	if !ok || (v.expiry != 0 && v.expiry <= atomic.LoadUint32(&timestamp)) {
		return nil, nil
	}
	return v.data, nil
}

// Has returns true if entry for the given key exists.
func (mk *memkey) Has(key string) bool {
	if len(key) <= 0 {
		return false
	}
	mk.mu.RLock()
	v, ok := mk.db[key]
	mk.mu.RUnlock()
	if !ok || (v.expiry != 0 && v.expiry <= atomic.LoadUint32(&timestamp)) {
		return false
	}
	return true
}

// Set key with value.
func (mk *memkey) Set(key string, val []byte, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}
	var expire uint32
	if exp != 0 {
		expire = uint32(exp.Seconds()) + atomic.LoadUint32(&timestamp)
	}
	e := entry{val, expire}
	mk.mu.Lock()
	mk.db[key] = e
	mk.mu.Unlock()
	return nil
}

// Delete key by key.
func (mk *memkey) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}
	mk.mu.Lock()
	delete(mk.db, key)
	mk.mu.Unlock()
	return nil
}

// Reset all keys.
func (mk *memkey) Reset() error {
	ndb := make(map[string]entry)
	mk.mu.Lock()
	mk.db = ndb
	mk.mu.Unlock()
	return nil
}

// Close the memkey.
func (mk *memkey) Close() error {
	mk.done <- struct{}{}
	return nil
}

// Keys returns all the keys.
func (mk *memkey) Keys() ([][]byte, error) {
	mk.mu.RLock()
	defer mk.mu.RUnlock()
	if len(mk.db) == 0 {
		return nil, nil
	}
	ts := atomic.LoadUint32(&timestamp)
	keys := make([][]byte, 0, len(mk.db))
	for key, v := range mk.db {
		if v.expiry == 0 || v.expiry > ts {
			keys = append(keys, []byte(key))
		}
	}
	if len(keys) == 0 {
		return nil, nil
	}
	return keys, nil
}

// Size returns the number of valid entries.
func (mk *memkey) Size() int {
	mk.mu.RLock()
	defer mk.mu.RUnlock()
	
	ts := atomic.LoadUint32(&timestamp)
	count := 0
	for _, v := range mk.db {
		if v.expiry == 0 || v.expiry > ts {
			count++
		}
	}
	return count
}

func startTimeStampUpdater() {
	updaterOnce.Do(func() {
		atomic.StoreUint32(&timestamp, uint32(time.Now().Unix()))
		go func(sleep time.Duration) {
			ticker := time.NewTicker(sleep)
			defer ticker.Stop()

			for t := range ticker.C {
				atomic.StoreUint32(&timestamp, uint32(t.Unix()))
			}
		}(1 * time.Second)
	})
}

func (mk *memkey) startCleanupLoop() {
	ticker := time.NewTicker(mk.cfg.CleanupInterval)
	defer ticker.Stop()
	var exp []string
	for {
		select {
		case <-mk.done:
			return
		case <-ticker.C:
			ts := atomic.LoadUint32(&timestamp)
			exp = exp[:0]
			mk.mu.RLock()
			for id, v := range mk.db {
				if v.expiry != 0 && v.expiry < ts {
					exp = append(exp, id)
				}
			}
			mk.mu.RUnlock()
			mk.mu.Lock()
			for i := range exp {
				v := mk.db[exp[i]]
				if v.expiry != 0 && v.expiry <= ts {
					delete(mk.db, exp[i])
				}
			}
			mk.mu.Unlock()
		}
	}
}
