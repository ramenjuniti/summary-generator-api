package lexrank

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/dcadenas/pagerank"
	"github.com/gaspiman/cosine_similarity"
	"github.com/ikawaha/kagome/tokenizer"
)

// SummaryData contains data for summary
type SummaryData struct {
	characters        int
	originalText      string
	originalSentences []string
	wordsPerSentence  [][]string
	tfScores          [][]float64
	idfScores         [][]float64
	tfIdfScores       [][]float64
	similarityMatrix  [][]float64

	lexRankScores []lexRankScore
	reRanking     []lexRankScore

	LineLimitedSummary      []lexRankScore
	CharacterLimitedSummary []lexRankScore

	maxLines      int
	maxCharacters int
	threshold     float64
	tolerance     float64
	damping       float64
	lambda        float64
}

type lexRankScore struct {
	Id       int     `json:"id"`
	Sentence string  `json:"sentence"`
	Score    float64 `json:"score"`
}

type Option func(*SummaryData)

const (
	delimiter            = "."
	defaultMaxLines      = 0
	defaultMaxCharacters = 0
	defaultThreshold     = 0.001
	defaultTolerance     = 0.0001
	defaultDamping       = 0.85
	defaultLambda        = 1
)

// MaxLines set SummaryData.maxLines
func MaxLines(maxLines int) Option {
	return func(args *SummaryData) {
		if maxLines < 0 {
			return
		}
		args.maxLines = maxLines
	}
}

// MaxCharacters set SummaryData.maxCharacters
func MaxCharacters(maxCharacters int) Option {
	return func(args *SummaryData) {
		if maxCharacters < 0 {
			return
		}
		args.maxCharacters = maxCharacters
	}
}

// Threshold set SummaryData.threshold
func Threshold(threshold float64) Option {
	return func(args *SummaryData) {
		args.threshold = threshold
	}
}

// Tolerance set SummaryData.tolerance
func Tolerance(tolerance float64) Option {
	return func(args *SummaryData) {
		args.tolerance = tolerance
	}
}

// Damping set SummaryData.damping
func Damping(damping float64) Option {
	return func(args *SummaryData) {
		args.damping = damping
	}
}

// Lambda set SummaryData.lambda
func Lambda(lambda float64) Option {
	return func(args *SummaryData) {
		args.lambda = lambda
	}
}

// New return SummaryData
func New(options ...Option) *SummaryData {
	summaryData := &SummaryData{
		maxLines:      defaultMaxLines,
		maxCharacters: defaultMaxCharacters,
		threshold:     defaultThreshold,
		tolerance:     defaultTolerance,
		damping:       defaultDamping,
		lambda:        defaultLambda,
	}
	for _, option := range options {
		option(summaryData)
	}
	return summaryData
}

// Summarize generate summary
func (s *SummaryData) Summarize(text string) {
	if len(text) == 0 {
		fmt.Println("input isn't specifyed.")
		return
	}
	s.changeSentenceEnd(text)
	s.countCharacter()
	s.splitText()
	s.splitSentence()
	s.calculateTf()
	s.calculateIdf()
	s.calculateTfidf()
	s.createSimilarityMatrix()
	s.calculateLexRank()
	s.calculateMmr()
	s.createLineLimitedSummary()
	sort.Slice(s.LineLimitedSummary, func(i, j int) bool {
		return s.LineLimitedSummary[i].Id < s.LineLimitedSummary[j].Id
	})
	s.createCharacterLimitedSummary()
	sort.Slice(s.CharacterLimitedSummary, func(i, j int) bool {
		return s.CharacterLimitedSummary[i].Id < s.CharacterLimitedSummary[j].Id
	})
}

func (s *SummaryData) changeSentenceEnd(text string) {
	if strings.Contains(text, "。") {
		s.originalText = strings.Replace(text, "。", delimiter, -1)
	}
	if strings.Contains(text, "！") {
		s.originalText = strings.Replace(text, "！", delimiter, -1)
	}
	if strings.Contains(text, "？") {
		s.originalText = strings.Replace(text, "？", delimiter, -1)
	}
	if strings.Contains(text, "!") {
		s.originalText = strings.Replace(text, "!", delimiter, -1)
	}
	if strings.Contains(text, "?") {
		s.originalText = strings.Replace(text, "?", delimiter, -1)
	}
}

func (s *SummaryData) countCharacter() {
	s.characters = utf8.RuneCountInString(s.originalText)
}

func (s *SummaryData) splitText() {
	sentences := strings.Split(s.originalText, delimiter)
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
	s.lexRankScores = make([]lexRankScore, len(s.originalSentences))
	for i, similarityList := range s.similarityMatrix {
		for j, similarity := range similarityList {
			if similarity >= s.threshold {
				graph.Link(i, j)
			}
		}
	}
	graph.Rank(s.damping, s.tolerance, func(identifier int, rank float64) {
		s.lexRankScores[identifier] = lexRankScore{Id: identifier, Sentence: s.originalSentences[identifier], Score: rank}
	})
	sort.Slice(s.lexRankScores, func(i, j int) bool {
		return s.lexRankScores[i].Score > s.lexRankScores[j].Score
	})
}

func (s *SummaryData) calculateMmr() {
	if len(s.lexRankScores) == 0 {
		return
	}
	s.reRanking = []lexRankScore{s.lexRankScores[0]}
	for len(s.lexRankScores) > len(s.reRanking) {
		var maxMmr float64
		var maxMmrId int
	L:
		for i, unselected := range s.lexRankScores {
			var maxSim float64
			for _, selected := range s.reRanking {
				if unselected.Id == selected.Id {
					continue L
				}
				if currentSim, _ := cosine_similarity.Cosine(s.tfIdfScores[unselected.Id], s.tfIdfScores[selected.Id]); currentSim > maxSim {
					maxSim = currentSim
				}
			}
			if currentMmr := s.lambda*unselected.Score - (1-s.lambda)*maxSim + 1; currentMmr > maxMmr {
				maxMmr = currentMmr
				maxMmrId = i
				fmt.Println(i)
			}
		}
		s.reRanking = append(s.reRanking, s.lexRankScores[maxMmrId])
	}
}

func (s *SummaryData) createLineLimitedSummary() {
	s.LineLimitedSummary = []lexRankScore{}
	if s.maxLines >= len(s.originalSentences) {
		s.LineLimitedSummary = append(s.LineLimitedSummary, s.reRanking...)
		return
	}
	s.LineLimitedSummary = append(s.LineLimitedSummary, s.reRanking[:s.maxLines]...)
}

func (s *SummaryData) createCharacterLimitedSummary() {
	s.CharacterLimitedSummary = []lexRankScore{}
	if s.maxCharacters >= s.characters {
		s.CharacterLimitedSummary = append(s.CharacterLimitedSummary, s.lexRankScores...)
		return
	}
	n := len(s.originalSentences)
	value := make([]float64, n)
	weight := make([]int, n)
	dp := make([][]float64, n+1)
	use := make([][]bool, n+1)
	for i, v := range s.lexRankScores {
		value[i] = v.Score
		weight[i] = utf8.RuneCountInString(v.Sentence)
	}
	for i := 0; i < n+1; i++ {
		dp[i] = make([]float64, s.maxCharacters+1)
		use[i] = make([]bool, s.maxCharacters+1)
	}
	for i := 1; i < n+1; i++ {
		for j := 1; j < s.maxCharacters+1; j++ {
			if j > weight[i-1] {
				dp[i][j] = math.Max((dp[i-1][j-weight[i-1]] + value[i-1]), dp[i-1][j])
				if dp[i-1][j-weight[i-1]]+value[i-1] > dp[i-1][j] {
					use[i][j] = true
				}
			} else {
				dp[i][j] = dp[i-1][j]
			}
		}
	}
	i := n
	j := s.maxCharacters
	for i > 0 {
		if use[i][j] {
			s.CharacterLimitedSummary = append(s.CharacterLimitedSummary, s.lexRankScores[i-1])
			j -= utf8.RuneCountInString(s.lexRankScores[i-1].Sentence)
		}
		i--
	}
}
