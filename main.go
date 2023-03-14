package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/yeralin-munar/words-count/cache"
	"github.com/yeralin-munar/words-count/utils"
)

var (
	n = flag.Int("n", 10, "Unique number of words")
	f = flag.String("f", "", "Path to input file")
)

const SEARCH_DIRECTORY = "searches/"
const OUTPUT_FILENAME = "output.txt"

func main() {
	flag.Parse()

	// Create LRU cache
	lruCache := cache.NewLRUCache(*n, removeCallback)

	// Open input file
	file, err := os.Open(*f)
	if err != nil {
		log.Fatalf("read file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Start scan line by line
	for scanner.Scan() {
		searches := scanner.Text()
		count := lruCache.Get(searches) + 1
		lruCache.Store(searches, count)
	}

	// Clear cache
	lruCache.Clear()

	err = scanner.Err()
	if err != nil {
		log.Fatalf("scanner error: %s", err)
	}

	// Save result into output file
	err = saveOutput(OUTPUT_FILENAME)
	if err != nil {
		log.Fatalf("save output: %s", err)
		return
	}

	// Remove folder with search files
	err = utils.RemoveFolder(SEARCH_DIRECTORY)
	if err != nil {
		log.Fatalf("remove folder and everythin inside: %s", err)
	}
}

func init() {
	// Create folder for searches word files
	err := utils.CreateFolder(SEARCH_DIRECTORY)
	if err != nil {
		log.Fatalf("create folder: %s", err)
	}

	err = utils.RemoveFile(OUTPUT_FILENAME)
	if err != nil {
		log.Fatalf("remove output file %s: %s", OUTPUT_FILENAME, err)
	}
}

func removeCallback(word string, count int) {
	file := utils.StringToUnicode(word)
	bytes := make([]byte, count)
	err := utils.WriteToFile(SEARCH_DIRECTORY+file, bytes)
	if err != nil {
		log.Fatalf("increment file: %s", err)
	}
}

func saveOutput(filename string) error {
	searchesDir, err := os.Open(SEARCH_DIRECTORY)
	if err != nil {
		return errors.New(fmt.Sprintf("open %s: %s", SEARCH_DIRECTORY, err))
	}
	defer searchesDir.Close()

	var (
		buf   bytes.Buffer
		files []fs.DirEntry
	)

	// Read files from dir every n
	for err != io.EOF {
		// Read first n files in dir
		files, err = searchesDir.ReadDir(*n)
		if err != nil && err != io.EOF {
			return err
		}

		for _, f := range files {
			// Transform file name into word
			word := utils.UnicodeSToString(f.Name())
			info, err := f.Info()
			if err != nil {
				return errors.New(fmt.Sprintf("get info from file %s: %s", f.Name(), err))
			}
			count := info.Size()
			buf.WriteString(fmt.Sprintf("%s - %d\n", word, count))

			// Remove file to read new
			err = utils.RemoveFile(SEARCH_DIRECTORY + f.Name())
			if err != nil {
				return errors.New(fmt.Sprintf("remove file %s: %s", f.Name(), err))
			}
		}
	}

	err = utils.WriteToFile(filename, buf.Bytes())
	if err != nil {
		return errors.New(fmt.Sprintf("write to file %s: %s", filename, err))
	}
	buf.Reset()

	return nil
}
