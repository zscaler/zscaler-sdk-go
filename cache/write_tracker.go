package cache

import (
	"strings"
	"sync"
)

// WriteTracker keeps a read that overlaps a write from caching stale data.
//
// Clients invalidate the cache entries a write affects before sending the
// write. A read of one of those entries that is in flight at the same time can
// still receive the data from before the write and store it, and later reads
// are then served that stale response until it expires. This happens when a
// list is read while an object of that list is being updated, for example
// when Terraform updates several objects of the same type in parallel.
//
// A client registers each write with BeginWrite, passing the cache key
// prefixes the write invalidates, and each cacheable read with BeginRead. A
// read whose key starts with one of the prefixes of a write that was in flight
// when the read started, or that started or finished while the read was in
// flight, is marked as overlapping, and its response must not be cached.
//
// The zero value is ready to use. A nil *WriteTracker tracks nothing and
// allows caching every read.
type WriteTracker struct {
	mu     sync.Mutex
	writes map[string]int          // prefix -> number of writes in flight
	reads  map[*TrackedRead]string // read in flight -> its cache key
}

// TrackedRead is a read registered with BeginRead.
type TrackedRead struct {
	t           *WriteTracker
	overlapping bool
	done        bool
}

// BeginWrite registers a write that invalidates every cache key starting with
// one of prefixes. The returned function must be called once the write has
// finished, successfully or not. Calling it more than once has no effect.
func (t *WriteTracker) BeginWrite(prefixes ...string) (end func()) {
	if t == nil || len(prefixes) == 0 {
		return func() {}
	}
	t.mu.Lock()
	if t.writes == nil {
		t.writes = map[string]int{}
	}
	for _, p := range prefixes {
		t.writes[p]++
	}
	t.markOverlappingLocked(prefixes)
	t.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			t.mu.Lock()
			defer t.mu.Unlock()
			for _, p := range prefixes {
				if t.writes[p] <= 1 {
					delete(t.writes, p)
				} else {
					t.writes[p]--
				}
			}
			t.markOverlappingLocked(prefixes)
		})
	}
}

// BeginRead registers a read of key. Call Done on the result once the read has
// finished, successfully or not.
func (t *WriteTracker) BeginRead(key string) *TrackedRead {
	if t == nil {
		return nil
	}
	r := &TrackedRead{t: t}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.reads == nil {
		t.reads = map[*TrackedRead]string{}
	}
	for p := range t.writes {
		if strings.HasPrefix(key, p) {
			r.overlapping = true
			break
		}
	}
	t.reads[r] = key
	return r
}

// Done unregisters the read and reports whether its response may be cached:
// no write affecting its key was in flight when it started, and none started
// or finished while it was in flight. Only the first call reports this; later
// calls return false.
func (r *TrackedRead) Done() bool {
	if r == nil {
		return true
	}
	r.t.mu.Lock()
	defer r.t.mu.Unlock()
	if r.done {
		return false
	}
	r.done = true
	delete(r.t.reads, r)
	return !r.overlapping
}

func (t *WriteTracker) markOverlappingLocked(prefixes []string) {
	for r, key := range t.reads {
		if r.overlapping {
			continue
		}
		for _, p := range prefixes {
			if strings.HasPrefix(key, p) {
				r.overlapping = true
				break
			}
		}
	}
}
