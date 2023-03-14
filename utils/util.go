package utils

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

func CreateFolder(name string) error {
	err := os.MkdirAll(name, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

// Remove all directories and files
func RemoveFolder(name string) error {
	err := os.RemoveAll(name)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

func RemoveFile(name string) error {
	err := os.Remove(name)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func WriteToFile(pathToFile string, bytes []byte) error {
	f, err := os.OpenFile(pathToFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(bytes); err != nil {
		return err
	}

	return nil
}

// StringToUnicode encode all string's chars
// into their stringified unicode
func StringToUnicode(s string) string {
	var newString string
	for i, l := range []rune(s) {
		newString += strconv.Itoa(int(l))
		if i < len(s)-1 {
			newString += "_"
		}
	}

	return newString
}

// UnicodeSToString decode stringified uncode chars
// which are concatenated with underscore
// into symbols then into string
func UnicodeSToString(s string) string {
	var runes []rune
	s = strings.Trim(s, "_")
	split := strings.Split(s, "_")
	for _, l := range split {
		r, err := strconv.Atoi(l)
		if err != nil {
			log.Fatalf("convert %s into int: %s\n", l, err)
		}
		runes = append(runes, rune(r))
	}

	return string(runes)
}
