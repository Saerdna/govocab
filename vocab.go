package govocab

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type VocabDict struct{
	Dict map[string]int32
	Lock sync.Mutex
	rank int32
}
const (
	DictLength = 2
	SepField = " "
)

func NewVocabDict()(*VocabDict){
	return &VocabDict{
		Dict:make(map[string]int32),
		Lock:sync.Mutex{},
		rank:0,
	}
}

func (s *VocabDict) LoadDict(path string) (err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	newDict := make(map[string]int32)
	var newRank int32 = 0
	for scanner.Scan() {
		line := scanner.Text()
		text := strings.Split(line, SepField)
		if len(text) != DictLength {
			continue
		}
		rank, err := strconv.ParseInt(text[1], 10, 32)
		if err != nil {
			continue
		}
		newDict[text[0]] = int32(rank)
		if newRank < int32(rank){
			newRank = int32(rank)
		}
	}
	s.Lock.Lock()
	s.Dict = newDict
	s.rank = newRank
	s.Lock.Unlock()
	return nil
}

func (s *VocabDict)Fit(words string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	for _, word := range []rune(words) {
		if string(word) == SepField{
			continue
		}
		if _, ok := s.Dict[string(word)]; !ok {
			s.rank += 1
			s.Dict[string(word)] = s.rank
		}
	}
}

func (s *VocabDict) Transform(words string, max_length int32) []int32 {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	result := make([]int32, max_length)
	for idx, word := range[]rune(words) {
		if idx >= int(max_length) {
			return result
		}
		result[idx] = s.Dict[string(word)]
	}
	return result
}
func (s *VocabDict) Save(path string) (err error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for k, v := range s.Dict {
		writer.WriteString(fmt.Sprintf("%s%s%d\n", k, SepField,v))
	}
	writer.Flush()
	return
}
