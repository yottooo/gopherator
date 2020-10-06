package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Records map[string]string

// History is a data structure for safely storing translations
type History struct {
	data  Records
	guard sync.RWMutex
}

// Store is a method for saving a translation
func (h *History) Store(en string, gf string) {
	h.guard.Lock()
	defer h.guard.Unlock()
	if h.data == nil {
		h.data = make(Records, 0)
	}
	h.data[en] = gf
}

// Load is a method that returns a translation if stored already or false if it's not
func (h *History) Load(en string) (string, bool) {
	h.guard.RLock()
	defer h.guard.RUnlock()
	str, ok := h.data[en]
	return str, ok
}

// GetData returns a copy of the stored translations
func (h *History) GetData() Records {
	h.guard.RLock()
	defer h.guard.RUnlock()
	cp := make(Records, len(h.data))
	for key, val := range h.data {
		cp[key] = val
	}
	return cp
}

type entry struct {
	English string
	Gopher  string
}

func (e *entry) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("{\"%s\":\"%s\"}", e.English, e.Gopher)
	return []byte(str), nil
}

type byAlphabet []string

func (a byAlphabet) Len() int           { return len(a) }
func (a byAlphabet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAlphabet) Less(i, j int) bool { return strings.Compare(a[i], a[j]) < 0 }

// ToJSON returns json encoding of History object
func (h *History) ToJSON() ([]byte, error) {
	data := h.GetData()

	keys := make([]string, 0)
	for key := range data {
		keys = append(keys, key)
	}

	sort.Sort(byAlphabet(keys))

	sortedData := make([]entry, len(keys))

	for idx, key := range keys {
		sortedData[idx] = entry{English: key, Gopher: data[key]}
	}

	return json.Marshal(&struct {
		History []entry `json:"history"`
	}{
		History: sortedData,
	})
}
