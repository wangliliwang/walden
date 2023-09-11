package main

import (
	"encoding/csv"
	"os"
	"sync"
)

type Record struct {
	Word      string
	Translate string
	Phonetic  string
}

type ECDict struct {
	once sync.Once
	data map[string]*Record
	path string
}

func NewECDict() *ECDict {
	return &ECDict{
		once: sync.Once{},
		data: nil,
		path: "dict/ecdict.csv",
	}
}

func (e *ECDict) load() {
	fl, openErr := os.Open(e.path)
	if openErr != nil {
		panic(openErr)
	}
	csvReader := csv.NewReader(fl)
	records, readErr := csvReader.ReadAll()
	if readErr != nil {
		panic(readErr)
	}
	e.data = make(map[string]*Record)
	for _, record := range records {
		word, phonetic, translate := record[0], record[1], record[3]
		e.data[word] = &Record{
			Word:      word,
			Translate: translate,
			Phonetic:  phonetic,
		}
	}
}

func (e *ECDict) lazyLoad() {
	// 查询第一个单词的时候，才会load
	e.once.Do(e.load)
}

func (e *ECDict) Match(key string) (*Record, bool) {
	e.lazyLoad()
	r, ok := e.data[key]
	return r, ok
}
