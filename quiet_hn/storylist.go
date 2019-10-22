package main

import (
	"fmt"
	"net/url"
	"quiet_hn/hn"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	pending = iota
	invalid
	valid
)

type StoryList struct {
	Items  []Item
	Status []int

	mu         sync.Mutex
	validCount int32
	ptr        int
	doneCh     chan<- struct{}
	done       bool
}

func (i *StoryList) Add(item itemPayload) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Items[item.pos] = item.item
	i.addStatus(item)
	i.isDoneChecker()
}

func (i *StoryList) addStatus(item itemPayload) {
	if isStoryLink(item.item) {
		i.Status[item.pos] = valid
	} else {
		i.Status[item.pos] = invalid
	}
}

func (i *StoryList) isDoneChecker() {
	if i.done {
		return
	}
	for ; i.Status[i.ptr] != pending; i.ptr++ {
		fmt.Println("i.ptr", i.ptr)
		if i.Status[i.ptr] == valid {
			atomic.AddInt32(&i.validCount, 1)
		}
		if i.validCount >= 30 {
		fmt.Println("sending done signal")
			i.doneCh <- struct{}{}
			close(i.doneCh)
			i.done = true
			return
		}
	}
}

func (i *StoryList) GetList() []Item {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.Items
}

// Item is the same as the hn.Item, but adds the Host field
type Item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []Item
	Time    time.Duration
}

func isStoryLink(item Item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) Item {
	ret := Item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}
