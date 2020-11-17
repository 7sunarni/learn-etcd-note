package main

import (
	"log"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	ticker := NewSandboxTTLTicker()
	// add
	ticker.InsertOrUpdate("1111", 1111)
	ticker.InsertOrUpdate("2222", 2222)
	ticker.InsertOrUpdate("3333", 3333)

	// update
	ticker.InsertOrUpdate("2222", 4444)
	ticker.InsertOrUpdate("2222", 6666)

	for i := 0; i < 3; i++ {
		log.Println(ticker.Pop())
	}

}

func TestPeek(t *testing.T) {
	ticker := NewSandboxTTLTicker()
	// add
	ticker.InsertOrUpdate("1111", 1111)
	ticker.InsertOrUpdate("2222", 2222)
	ticker.InsertOrUpdate("3333", 3333)

	for i := 0; i < 2; i++ {
		log.Println(ticker.Peek())
	}
}

func TestDelete(t *testing.T) {
	ticker := NewSandboxTTLTicker()
	// add
	ticker.InsertOrUpdate("1111", time.Now().Add(time.Second*5).UnixNano())
	ticker.InsertOrUpdate("2222", time.Now().Add(time.Second*7).UnixNano())
	ticker.InsertOrUpdate("3333", time.Now().Add(time.Second*9).UnixNano())
	ticker.InsertOrUpdate("4444", time.Now().Add(time.Second*15).UnixNano())

	ticker.InsertOrUpdate("2222", time.Now().Add(time.Second*3).UnixNano())

	for {

	}
}
