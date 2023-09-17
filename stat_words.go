package main

import (
	"encoding/csv"
	"fmt"
	"github.com/rodaine/table"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
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
	logrus.Infof("[read] read file done, [%s].", path)
	return string(bs)
}

type Word struct {
	Word     string
	Count    int
	Chinese  string
	Phonetic string
}

func (w Word) String() string {
	return fmt.Sprintf("%s %d %s", w.Word, w.Count, w.Chinese)
}

type Words []Word

func (ws Words) Len() int {
	return len(ws)
}
func (ws Words) Less(i, j int) bool {
	// 词频倒序，字典序正序
	if ws[i].Count > ws[j].Count {
		return true
	} else if ws[i].Count < ws[j].Count {
		return false
	} else {
		return ws[i].Word < ws[j].Word
	}
}
func (ws Words) Swap(i, j int) {
	ws[i], ws[j] = ws[j], ws[i]
}

type WordSet map[string]struct{}

var placeholder struct{}

// return: 包含本章新单词的全部单词，本章新单词
func statWords(text string, everWordSet WordSet) (WordSet, Words) {
	counter := map[string]int{}
	// 找出word
	for _, preWord := range strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || (unicode.IsPunct(r) && r != '-')
	}) {
		// 变成小写
		preWord = strings.ToLower(preWord)
		// 略过存量单词
		if _, ok := everWordSet[preWord]; ok {
			continue
		}
		counter[preWord]++
	}

	// 将新单词，合并进入wordSet
	for word := range counter {
		everWordSet[word] = placeholder
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
			w.Phonetic = dictRecord.Phonetic
		}
		words = append(words, w)
	}
	sort.Sort(words)
	logrus.Infof("[stat] stat words done, all words [%d], new words [%d].", len(everWordSet), len(words))
	return everWordSet, words
}

func exportToCSV(path string, words Words) {
	fl, createErr := os.Create(path)
	if createErr != nil {
		panic(createErr)
	}
	csvWriter := csv.NewWriter(fl)
	// transform
	records := make([][]string, len(words)+1)
	records[0] = []string{"word", "phonetic", "count", "translation"}
	for index, word := range words {
		records[index+1] = []string{word.Word, word.Phonetic, strconv.Itoa(word.Count), word.Chinese}
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

// 统计每章的单词
func statAndWriteWords(everWords WordSet, rawFile string) (WordSet, Words) {
	wordSet, words := statWords(readFile(rawFile), everWords)
	parts := strings.Split(rawFile, "/")
	exportFile := fmt.Sprintf("dist/%s.csv", parts[len(parts)-1])
	exportToCSV(exportFile, words)
	return wordSet, words
}

func parseFileNo(path string) int {
	// 现根据 / 分隔，取最后一部分

	ss := strings.Split(path, "/")
	if len(ss) == 0 {
		panic("invalid path")
	}
	path = ss[len(ss)-1]

	ss = strings.Split(path, "-")
	if len(ss) <= 1 {
		panic("invalid path")
	}
	r, e := strconv.Atoi(ss[0])
	if e != nil {
		panic(e)
	}
	return r
}

type Files []string

func (x Files) Len() int { return len(x) }
func (x Files) Less(i, j int) bool {
	return parseFileNo(x[i]) < parseFileNo(x[j])
}
func (x Files) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

type StatInfo struct {
	Name         string
	AllWordCount int
	NewWordCount int
}

type StatInfos []StatInfo

func (s StatInfos) Print() {
	tbl := table.New("RawFile", "AllWordCount", "NewWordCount")
	for _, item := range s {
		tbl.AddRow(item.Name, item.AllWordCount, item.NewWordCount)
	}
	tbl.Print()
}

func getRawFiles() Files {
	// 读取 raw 文件夹的文件，排序
	entries, readDirErr := os.ReadDir("raw/")
	if readDirErr != nil {
		panic(readDirErr)
	}
	files := make(Files, 0)
	fileNameRegex := regexp.MustCompile(`^\d+-.*$`)
	for _, entry := range entries {
		fileName := entry.Name()
		if !fileNameRegex.MatchString(fileName) {
			logrus.Errorf("[parse] invalid filename [%s]", fileName)
		}
		files = append(files, fmt.Sprintf("raw/%s", fileName))
	}
	// sort files
	sort.Sort(files)
	logrus.Infof("[parse] [%d] files to be stat, they are [%s]", len(files), strings.Join(files, ","))
	return files
}

func main() {
	// init log
	initLog()
	// init dir
	initDir()

	files := getRawFiles()

	// stat
	wordSet := make(WordSet)
	var words Words
	statInfos := make(StatInfos, 0)
	for _, rawFile := range files {
		wordSet, words = statAndWriteWords(wordSet, rawFile)
		statInfos = append(statInfos, StatInfo{
			Name:         rawFile,
			AllWordCount: len(wordSet),
			NewWordCount: len(words),
		})
	}

	// print stat info
	statInfos.Print()
}
