package govocab

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	Dict map[string]int32
)

const (
	DictLength = 2
)

func LoadDict(path string) (err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	newDict := make(map[string]int32)
	for scanner.Scan() {
		line := scanner.Text()
		text := strings.Split(line, "\t")
		if len(text) != DictLength {
			continue
		}
		rank, err := strconv.ParseInt(text[1], 10, 32)
		if err != nil {
			continue
		}
		newDict[text[0]] = rank
	}
	Dict = newDict
}

func Fit(wordList []string) {
	var idx int32 = 1
	for _, word := range wordList {
		for _, vocab := range strings.Split(word, " ") {
			if _, ok := Dict[vocab]; !ok {
				Dict[vocab] = idx
				idx += 1
			}
		}
	}
}

func transform(wordList []string) [][]int32 {
	result := make([][]int32, 0)
	for _, word := range wordList {
		one := make([]int32, 0)
		for _, vocab := range strings.Split(word, " ") {
			one = append(one, Dict[vocab])
		}
		result = append(result, one)
	}
	return result
}
func Save(path string) (err error) {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for k, v := range Dict {
		writer.WriteString(fmt.Sprintf("%s %d", k, v))
	}
}
