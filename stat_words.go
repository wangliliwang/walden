package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func readFile(path string) string {
	bs, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

type Word struct {
	Word    string
	Count   int
	Chinese string
}

func (w Word) String() string {
	return fmt.Sprintf("%s %d %s", w.Word, w.Count, w.Chinese)
}

type Words []Word

func (ws Words) Len() int {
	return len(ws)
}
func (ws Words) Less(i, j int) bool {
	return ws[i].Count > ws[j].Count
}
func (ws Words) Swap(i, j int) {
	ws[i], ws[j] = ws[j], ws[i]
}

// 没有找到好用的翻译字典。下午找找看。
func statWords(text string) Words {
	counter := map[string]int{}
	// 找出word
	for _, preWord := range strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || (unicode.IsPunct(r) && r != '-')
	}) {
		// 变成小写
		preWord = strings.ToLower(preWord)
		counter[preWord]++
	}

	// sort
	words := make(Words, 0)
	dict := NewECDict()
	for word, count := range counter {
		w := Word{
			Word:  word,
			Count: count,
		}

		// 查dic
		dictRecord, ok := dict.Match(word)
		if ok {
			w.Chinese = dictRecord.Translate
		}
		words = append(words, w)
	}
	sort.Sort(words)
	return words
}

func exportToCSV(path string, words Words) {
	fl, createErr := os.Create(path)
	if createErr != nil {
		panic(createErr)
	}
	csvWriter := csv.NewWriter(fl)
	// transform
	records := make([][]string, len(words)+1)
	records[0] = []string{"word", "count", "translation"}
	for index, word := range words {
		records[index+1] = []string{word.Word, strconv.Itoa(word.Count), word.Chinese}
	}

	// write
	writeErr := csvWriter.WriteAll(records)
	if writeErr != nil {
		panic(writeErr)
	}
	csvWriter.Flush()
	logrus.Infof("[export] export to csv done, path = [%s], total lines = [%d].", path, len(words))
}

func initLog() {
	customFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	logrus.SetFormatter(customFormatter)
	logrus.SetOutput(os.Stdout)
}

func initDir() {
	mkdirErr := os.MkdirAll("dist/", os.ModePerm)
	if mkdirErr != nil {
		panic(mkdirErr)
	}
}

func main() {
	// init log
	initLog()
	// init dir
	initDir()

	// parse args
	rawFile := flag.String("raw", "", "raw file")
	flag.Parse()
	if rawFile == nil || *rawFile == "" {
		panic("raw参数必填")
	}

	words := statWords(readFile(*rawFile))
	parts := strings.Split(*rawFile, "/")
	exportFile := fmt.Sprintf("dist/%s.csv", parts[len(parts)-1])
	exportToCSV(exportFile, words)
}
