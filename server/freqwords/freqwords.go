package freqwords

import (
	"bufio"
	"container/heap"
	"fmt"
	"hash/crc32"
	"os"
	"strings"
)

const nmappers = 12
const nreducers = 8

var Order string

// records the number of occurrences of a particular word
type Wc struct {
	Word  string
	Count int
}

// This heap (priority queue) structure is used by the reducers
// to send their most common words to back to the main thread.
type wcHeap []*Wc

func (wch wcHeap) Len() int { return len(wch) }
func (wch wcHeap) Less(i, j int) bool {
	// we want the *most* commonly occurring words (largest counts)
	if Order == "asc" {
		return wch[i].Count < wch[j].Count
	}
	return wch[i].Count > wch[j].Count
}
func (wch wcHeap) Swap(i, j int) {
	wch[i], wch[j] = wch[j], wch[i]
}
func (wch *wcHeap) Push(x interface{}) {
	*wch = append(*wch, x.(*Wc))
}
func (wch *wcHeap) Pop() interface{} {
	old := *wch
	n := len(old)
	wcItem := old[n-1]
	*wch = old[0 : n-1]
	return wcItem
}

// End heap definitions

// Read the (already-opened) directory, open each item,
// send the non-director to the mappers, recurse on
// the directories.
func walkRecursive(files chan *os.File, fname string) {
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println("bad", err)
		os.Exit(1)
	}
	if strings.HasSuffix(fname, ".txt") {
		files <- file
		return
	}
	dirContents, err := file.Readdirnames(0)
	if err != nil {
		// this is a non-txt file, non-directory, ignore
		file.Close()
		return
	}
	// fname is a directory, its contents are in 'dirContents'
	{
		err := os.Chdir(fname)
		if err != nil {
			fmt.Println("bad", err)
		}
	}
	for _, name := range dirContents {
		walkRecursive(files, name)
	}
	os.Chdir("..")
	file.Close()
}

func walk(files chan *os.File, dir string) {
	walkRecursive(files, dir)
	close(files)
}

func CommonWords(startDir string, numWords int) []Wc {
	// for hashing the words
	crc32q := crc32.MakeTable(0xD5828281)

	// each reducer will send its most common words to here
	results := make(chan *Wc)

	// start recursively reading directories and
	// writing file descriptors to channel 'files'
	files := make(chan *os.File, 4)
	go walk(files, startDir)

	// *** Reducer threads:
	var reduceCh [nreducers]chan string
	for i := range reduceCh {
		noDups := make(map[string]int)
		reduceCh[i] = make(chan string, 8)
		go func(i int) {
			for {
				word, ok := <-reduceCh[i]
				if !ok {
					break
				}
				// creates the map entry (initializes count to 0) if needed
				noDups[word]++
			}
			wch := make(wcHeap, 0)
			heap.Init(&wch)
			for word, count := range noDups {
				wcItem := &Wc{word, count}
				heap.Push(&wch, wcItem)
			}
			for i := 0; i < numWords && wch.Len() > 0; i++ {
				wcItem := heap.Pop(&wch).(*Wc)
				results <- wcItem
			}
			// this tells the receiving thread that we're done
			results <- &Wc{"", 0}
		}(i)
	}

	// *** Mapper threads:
	mappersdone := make(chan struct{})
	for i := 0; i < nmappers; i++ {
		go func(i int) {
			for { // each file
				file, ok := <-files
				if !ok {
					// no more files to process
					break
				}
				input := bufio.NewScanner(file)
				input.Split(bufio.ScanWords)
				for input.Scan() { // each word
					// send each word to a pseudo-random reducer,
					// duplicate words go to the same reducer
					text := input.Text()
					ri := crc32.Checksum([]byte(text), crc32q) % nreducers
					reduceCh[ri] <- text
				}
				file.Close()
			}
			mappersdone <- struct{}{}
		}(i)
	}

	// wait for all the mappers to finish
	for i := 0; i < nmappers; i++ {
		<-mappersdone
	}
	// since the mappers are all done, we can close their
	// pipes to the reducers, so they can finish
	for i := 0; i < nreducers; i++ {
		close(reduceCh[i])
	}
	// read the results from the reducers, put into one final heap
	wch := make(wcHeap, 0)
	heap.Init(&wch)
	remaining := nreducers
	for {
		wcp := <-results
		if wcp.Count == 0 {
			// one of the reducers is done
			remaining--
			if remaining == 0 {
				// all the reducers are done
				break
			}
			continue
		}
		heap.Push(&wch, wcp)
	}
	// last step: create a slice of the most common words overall
	most := make([]Wc, 0)
	for i := 0; i < numWords && wch.Len() > 0; i++ {
		wcItem := heap.Pop(&wch).(*Wc)
		most = append(most, *wcItem)
	}
	return most
}
