package worder

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"log"
)

type Worder struct {
	Text    string
	Workers int
	Path    string
}

type word struct {
	word   string
	number int
}

const punct = ".,!?"

// Run blocks until completion or first error.
func (w *Worder) Run() {
	if w.Workers == 0 {
		w.Workers = 1
	}
	if w.Path == "" {
		w.Path = "."
	}

	s := w.Text
	for _, v := range punct {
		s = strings.ReplaceAll(s, string([]rune{v}), "")
	}
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	words := strings.Split(s, " ")

	wCh := make(chan word, 1000)
	errCh := make(chan error)

	var wg sync.WaitGroup
	wg.Add(w.Workers)
	for i := 0; i < w.Workers; i++ {
		go func() {
			w.worker(wCh, errCh)
			wg.Done()
		}()
	}
	go logErrors(errCh)

	for i, w := range words {
		wCh <- word{
			word:   w,
			number: i + 1,
		}
	}

	close(wCh)
	wg.Wait()
	close(errCh)
}

func (w *Worder) worker(wCh <-chan word, errCh chan<- error) {
	for word := range wCh {
		fileName := fmt.Sprintf("%v/%v-%v", strings.TrimSuffix(w.Path, "/"), word.number, word.word)
		data := strings.Repeat(word.word, 1000)
		err := ioutil.WriteFile(fileName, []byte(data), 0644)
		if err != nil {
			errCh <- err
		}
	}
}

func logErrors(errCh <-chan error) {
	for e := range errCh {
		log.Printf("Encountered error: %v", e)
	}
}
