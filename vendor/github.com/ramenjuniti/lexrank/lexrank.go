package lexrank

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/dcadenas/pagerank"
	"github.com/gaspiman/cosine_similarity"
	"github.com/ikawaha/kagome/tokenizer"
)

// SummaryData contains data for summary
type SummaryData struct {
	originalSentences []string
	wordsPerSentence  [][]string
	tfScores          [][]float64
	idfScores         [][]float64
	tfIdfScores       [][]float64
	similarityMatrix  [][]float64

	LexRankScores []lexRankScore

	maxCharacters int
	threshold     float64
	tolerance     float64
	damping       float64
}

type lexRankScore struct {
	Id       int     `json:"id"`
	Sentence string  `json:"sentence"`
	Score    float64 `json:"score"`
}

const (
	defaultMaxCharacters = 0
	defaultThreshold     = 0.5
	defaultTolerance     = 0.0001
	defaultDamping       = 0.85
)

// New return SummaryData
func New() *SummaryData {
	return &SummaryData{
		maxCharacters: defaultMaxCharacters,
		threshold:     defaultThreshold,
		tolerance:     defaultTolerance,
		damping:       defaultDamping,
	}
}

// Set set SummaryData params
func (s *SummaryData) Set(m int, th, to, d float64) {
	s.maxCharacters = m
	s.threshold = th
	s.tolerance = to
	s.damping = d
}

// Summarize generate summary
func (s *SummaryData) Summarize(text, delimiter string) {
	if len(text) == 0 || len(delimiter) == 0 {
		fmt.Println("input isn't specifyed.")
		return
	}
	s.splitText(text, delimiter)
	s.splitSentence()
	s.calculateTf()
	s.calculateIdf()
	s.calculateTfidf()
	s.createSimilarityMatrix()
	s.calculateLexRank()
}

func (s *SummaryData) splitText(text, delimiter string) {
	sentences := strings.Split(text, delimiter)
	s.originalSentences = sentences[:len(sentences)-1]
}

func (s *SummaryData) splitSentence() {
	s.wordsPerSentence = make([][]string, len(s.originalSentences))
	t := tokenizer.New()
	for i, sentence := range s.originalSentences {
		tokens := t.Tokenize(sentence)[1:]
		s.wordsPerSentence[i] = make([]string, len(tokens)-1)
		for j := 0; j < len(tokens)-1; j++ {
			s.wordsPerSentence[i][j] = tokens[j].Surface
		}
	}
}

func (s *SummaryData) calculateTf() {
	s.tfScores = make([][]float64, len(s.originalSentences))
	var allWordsCount float64
	for _, sentence := range s.wordsPerSentence {
		allWordsCount += float64(len(sentence))
	}
	for i, sentence1 := range s.wordsPerSentence {
		s.tfScores[i] = make([]float64, len(sentence1))
		for j, word1 := range sentence1 {
			var count float64
			for _, sentence2 := range s.wordsPerSentence {
				for _, word2 := range sentence2 {
					if word1 == word2 {
						count++
					}
				}
			}
			s.tfScores[i][j] = count / allWordsCount
		}
	}
}

func (s *SummaryData) calculateIdf() {
	s.idfScores = make([][]float64, len(s.originalSentences))
	n := float64(len(s.originalSentences))
	for i, sentence1 := range s.wordsPerSentence {
		s.idfScores[i] = make([]float64, len(sentence1))
		for j, word1 := range sentence1 {
			var count float64
			for _, sentence2 := range s.wordsPerSentence {
				for _, word2 := range sentence2 {
					if word1 == word2 {
						count++
						break
					}
				}
			}
			s.idfScores[i][j] = math.Log(n/count) + 1
		}
	}
}

func (s *SummaryData) calculateTfidf() {
	s.tfIdfScores = make([][]float64, len(s.originalSentences))
	for i := 0; i < len(s.originalSentences); i++ {
		s.tfIdfScores[i] = make([]float64, len(s.wordsPerSentence[i]))
		for j := 0; j < len(s.wordsPerSentence[i]); j++ {
			s.tfIdfScores[i][j] = s.tfScores[i][j] * s.idfScores[i][j]
		}
	}
}

func (s *SummaryData) createSimilarityMatrix() {
	s.similarityMatrix = make([][]float64, len(s.originalSentences))
	for i := range s.similarityMatrix {
		s.similarityMatrix[i] = make([]float64, len(s.originalSentences))
	}
	for i := 0; i < len(s.similarityMatrix); i++ {
		for j := i; j < len(s.similarityMatrix[i]); j++ {
			if i == j {
				s.similarityMatrix[i][j] = 1
				s.similarityMatrix[j][i] = 1
				continue
			} else {
				s.similarityMatrix[i][j], _ = cosine_similarity.Cosine(s.tfIdfScores[i], s.tfIdfScores[j])
				s.similarityMatrix[j][i] = s.similarityMatrix[i][j]
			}
		}
	}
}

func (s *SummaryData) calculateLexRank() {
	graph := pagerank.New()
	s.LexRankScores = make([]lexRankScore, len(s.originalSentences))
	for i, similarityList := range s.similarityMatrix {
		for j, similarity := range similarityList {
			if similarity >= s.threshold {
				graph.Link(i, j)
			}
		}
	}
	graph.Rank(s.damping, s.tolerance, func(identifier int, rank float64) {
		s.LexRankScores[identifier] = lexRankScore{Id: identifier, Sentence: s.originalSentences[identifier], Score: rank}
	})
	sort.Slice(s.LexRankScores, func(i, j int) bool {
		return s.LexRankScores[i].Score > s.LexRankScores[j].Score
	})
}
