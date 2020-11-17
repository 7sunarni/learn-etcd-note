package main

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"sync"

	"math/rand"
	"time"
)

func main() {
	ticker := NewSandboxTTLTicker()
	count := 10
	for i := 0; i < count; i++ {
		t := rand.Int63n(500)
		ticker.InsertOrUpdate(fmt.Sprintf("sandbox-%d", t), t)
	}

	fmt.Println(ticker)

}

type SandboxTTLTicker struct {
	*queue
	m         map[string]*TTL
	ticker    *time.Timer
	recentTTL *TTL // unix_nano
	lock      sync.Mutex
}

func NewSandboxTTLTicker() *SandboxTTLTicker {
	st := SandboxTTLTicker{
		queue:     &queue{},
		m:         make(map[string]*TTL),
		lock:      sync.Mutex{},
		recentTTL: &TTL{time: math.MaxInt64},
		ticker:    time.NewTimer(time.Duration(math.MaxInt64)),
	}

	heap.Init(st.queue)
	go st.tickerLoop()
	return &st
}

func (s *SandboxTTLTicker) tickerLoop() {
	for {
		select {
		case <-s.ticker.C:
			log.Printf("%s trigger", s.recentTTL.name)
			s.Remove(s.recentTTL.name)
			recent := s.Peek()
			if recent == nil {
				continue
			}
			s.recentTTL = recent
			s.ticker.Reset(time.Nanosecond * time.Duration(s.recentTTL.time-time.Now().UnixNano()))
		}
	}
}

func (s *SandboxTTLTicker) InsertOrUpdate(sbName string, ttl /*unixnano*/ int64) {

	// 判断 ttl 是否大于当前时间
	if time.Now().UnixNano() > ttl {
		log.Printf("%d is small than now %d", ttl, time.Now().UnixNano())
		return
	}

	s.lock.Lock()
	if old, ok := s.m[sbName]; ok {
		old.time = ttl
		heap.Fix(s.queue, old.index)
	} else {
		item := &TTL{
			name: sbName,
			time: ttl,
		}
		heap.Push(s.queue, item)
		s.m[sbName] = item
	}
	s.lock.Unlock()

	recent := s.Peek()

	if recent.time < s.recentTTL.time {
		s.recentTTL = recent
		s.ticker.Reset(time.Nanosecond * time.Duration(s.recentTTL.time-time.Now().UnixNano()))
	}
}

func (s *SandboxTTLTicker) Peek() *TTL {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.queue.Len() > 0 {
		return s.queue.First()
	}
	return nil
}

func (s *SandboxTTLTicker) Remove(sandboxID string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if old, ok := s.m[sandboxID]; ok {
		heap.Remove(s.queue, old.index)
		delete(s.m, sandboxID)
	}
}

func (s *SandboxTTLTicker) Pop() *TTL {
	// s.lock.Lock()
	// defer s.lock.Unlock()
	item := heap.Pop(s.queue).(*TTL)
	v := s.m[item.name]
	delete(s.m, item.name)
	return v
}

type TTL struct {
	name  string
	time  int64 // unix_nano
	index int
}

type queue []*TTL

func (pq queue) Len() int { return len(pq) }

func (pq queue) Less(i, j int) bool {
	return pq[i].time < pq[j].time
}

func (pq queue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *queue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*TTL)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *queue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq queue) First() *TTL {
	return pq[0]
}
