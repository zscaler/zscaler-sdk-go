package cache

import (
	"sync"
	"testing"
)

const (
	listKey  = "https://api/zia/api/v1/urlCategories?customOnly=true"
	itemKey  = "https://api/zia/api/v1/urlCategories/CUSTOM_14"
	listPfx  = "https://api/zia/api/v1/urlCategories"
	otherPfx = "https://api/zia/api/v1/firewallFilteringRules"
)

func TestWriteTracker_ReadWithoutWritesIsCacheable(t *testing.T) {
	var tr WriteTracker
	if !tr.BeginRead(listKey).Done() {
		t.Fatal("read without writes must be cacheable")
	}
}

func TestWriteTracker_WriteInFlightWhenReadStarts(t *testing.T) {
	var tr WriteTracker
	end := tr.BeginWrite(itemKey, listPfx)
	r := tr.BeginRead(listKey)
	end()
	if r.Done() {
		t.Fatal("read started while a write was in flight must not be cached")
	}
}

func TestWriteTracker_WriteStillInFlightWhenReadEnds(t *testing.T) {
	var tr WriteTracker
	end := tr.BeginWrite(listPfx)
	defer end()
	if tr.BeginRead(listKey).Done() {
		t.Fatal("read overlapping a write for its whole duration must not be cached")
	}
}

func TestWriteTracker_WriteStartsDuringRead(t *testing.T) {
	var tr WriteTracker
	r := tr.BeginRead(listKey)
	end := tr.BeginWrite(listPfx)
	defer end()
	if r.Done() {
		t.Fatal("read during which a write started must not be cached")
	}
}

func TestWriteTracker_WriteStartsAndFinishesDuringRead(t *testing.T) {
	var tr WriteTracker
	r := tr.BeginRead(listKey)
	tr.BeginWrite(listPfx)()
	if r.Done() {
		t.Fatal("read during which a write ran must not be cached")
	}
}

func TestWriteTracker_ReadAfterWriteFinishedIsCacheable(t *testing.T) {
	var tr WriteTracker
	tr.BeginWrite(listPfx)()
	if !tr.BeginRead(listKey).Done() {
		t.Fatal("read started after the write finished must be cacheable")
	}
}

func TestWriteTracker_UnrelatedWriteDoesNotAffectRead(t *testing.T) {
	var tr WriteTracker
	end := tr.BeginWrite(otherPfx)
	defer end()
	r := tr.BeginRead(listKey)
	tr.BeginWrite(otherPfx + "/1")()
	if !r.Done() {
		t.Fatal("read of another collection must stay cacheable")
	}
}

func TestWriteTracker_OverlappingWritesOfSamePrefix(t *testing.T) {
	var tr WriteTracker
	end1 := tr.BeginWrite(listPfx)
	end2 := tr.BeginWrite(listPfx)
	end1()
	if tr.BeginRead(listKey).Done() {
		t.Fatal("read must not be cached while the second write is in flight")
	}
	end2()
	if !tr.BeginRead(listKey).Done() {
		t.Fatal("read after both writes finished must be cacheable")
	}
}

func TestWriteTracker_EndAndDoneAreIdempotent(t *testing.T) {
	var tr WriteTracker
	end1 := tr.BeginWrite(listPfx)
	end2 := tr.BeginWrite(listPfx)
	end1()
	end1() // must not end the second write
	if tr.BeginRead(listKey).Done() {
		t.Fatal("calling a write's end twice must not end another write")
	}
	end2()

	r := tr.BeginRead(listKey)
	if !r.Done() {
		t.Fatal("first Done must report cacheable")
	}
	if r.Done() {
		t.Fatal("later Done calls must report false")
	}
}

func TestWriteTracker_ReleasesFinishedWritesAndReads(t *testing.T) {
	var tr WriteTracker
	for i := 0; i < 100; i++ {
		end := tr.BeginWrite(listPfx, itemKey)
		r := tr.BeginRead(listKey)
		end()
		r.Done()
	}
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if len(tr.writes) != 0 || len(tr.reads) != 0 {
		t.Fatalf("tracker kept %d writes and %d reads, want none", len(tr.writes), len(tr.reads))
	}
}

func TestWriteTracker_NilTrackerAndRead(t *testing.T) {
	var tr *WriteTracker
	tr.BeginWrite(listPfx)()
	if !tr.BeginRead(listKey).Done() {
		t.Fatal("nil tracker must allow caching")
	}
}

func TestWriteTracker_Concurrent(t *testing.T) {
	var tr WriteTracker
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			tr.BeginWrite(listPfx)()
		}()
		go func() {
			defer wg.Done()
			tr.BeginRead(listKey).Done()
		}()
	}
	wg.Wait()
	if !tr.BeginRead(listKey).Done() {
		t.Fatal("read after all writes finished must be cacheable")
	}
}
