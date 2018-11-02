package app

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

func Generate(text, delimiter string) {
	sentences := splitText(text, delimiter)
	textData := make([][]string, len(sentences))
	for i, sentence := range sentences {
		textData[i] = splitSentence(sentence)
	}
	fmt.Println(textData)
}

func splitText(text, delimiter string) []string {
	sentences := strings.Split(text, delimiter)
	return sentences[:len(sentences)-1]
}

func splitSentence(sentence string) []string {
	t := tokenizer.New()
	tokens := t.Tokenize(sentence)[1:]
	words := make([]string, len(tokens)-1)
	for i := 0; i < len(tokens)-1; i++ {
		words[i] = tokens[i].Surface
	}
	return words
}
