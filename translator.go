package main

import (
	"fmt"
	"strings"
)

// Error codes for translation errors
const (
	ErrConfusedGopher = 100 // word for translation is shortend
	ErrInvalidWord    = 101
)

// TransError is a transalation error stucture
type TransError struct {
	Err  string
	Code int
}

func (e *TransError) Error() string {
	return e.Err
}

// TranslateWord translates a single
func TranslateWord(word string) (string, error) {
	if len(word) == 0 {
		return "", &TransError{Code: ErrInvalidWord, Err: "No word was provided"}
	}

	if strings.ContainsAny(word, "’'") {
		return "", &TransError{Code: ErrConfusedGopher, Err: "Gophers can not understand shortened words"}
	}

	word = strings.ToLower(word)

	prefix := "g"
	vowelIdx := strings.Index(word, "xr")
	if vowelIdx == 0 {
		prefix = "ge"
	} else {

		vowelIdx = strings.IndexAny(word, "aeiou")
		if vowelIdx == -1 {
			vowelIdx = strings.Index(word, "y")
		}

		if vowelIdx >= 2 && word[vowelIdx-1:vowelIdx+1] == "qu" {
			vowelIdx++
		}
	}

	if vowelIdx == -1 {
		return "", &TransError{Code: ErrInvalidWord, Err: fmt.Sprintf("'%s' has no vowels", word)}
	}

	var builder strings.Builder

	if vowelIdx == 0 {
		builder.WriteString(prefix)
	}
	builder.WriteString(word[vowelIdx:len(word)])
	builder.WriteString(word[0:vowelIdx])
	if vowelIdx != 0 {
		builder.WriteString("ogo")
	}
	return builder.String(), nil
}

func extractSign(word string) (string, string) {
	var sign string
	if strings.LastIndexAny(word, ",.?!") == len(word)-1 {
		ln := len(word)
		sign = word[ln-1 : ln]
		word = word[:ln-1]
	}
	return word, sign
}

// TranslateSentence translates a whole sentence in gopher
func TranslateSentence(sentence string) (string, error) {
	english := strings.Split(sentence, " ")
	var gopher []string

	for _, word := range english {
		word, sign := extractSign(word)
		translated, err := TranslateWord(word)
		if e, ok := err.(*TransError); ok {
			if e.Code == ErrConfusedGopher {
				continue
			} else {
				return "", err
			}
		}

		gopher = append(gopher, translated+sign)
	}

	return strings.Join(gopher, " "), nil
}
