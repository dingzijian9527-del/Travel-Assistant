package rag

import (
	"hash/fnv"
	"math"
	"strings"
	"unicode"
)

// Embedder 把文本转换成固定维度向量。
type Embedder interface {
	Embed(text string) []float32
}

// HashEmbedder 是用于资料初始化的确定性本地向量器。
type HashEmbedder struct {
	dim int
}

// NewHashEmbedder 创建指定维度的确定性向量器。
func NewHashEmbedder(dim int) *HashEmbedder {
	if dim <= 0 {
		dim = 1
	}
	return &HashEmbedder{dim: dim}
}

// Embed 把文本词元转换成归一化哈希向量。
func (e *HashEmbedder) Embed(text string) []float32 {
	vector := make([]float32, e.dim)
	tokens := embeddingTokens(text)
	if len(tokens) == 0 {
		return vector
	}
	for _, token := range tokens {
		index, sign := hashToken(token, e.dim)
		vector[index] += sign
	}
	normalize(vector)
	return vector
}

func embeddingTokens(text string) []string {
	var tokens []string
	var current strings.Builder
	flush := func() {
		if current.Len() == 0 {
			return
		}
		tokens = append(tokens, strings.ToLower(current.String()))
		current.Reset()
	}
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
			continue
		}
		flush()
	}
	flush()
	return tokens
}

func hashToken(token string, dim int) (int, float32) {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(token))
	value := hash.Sum64()
	sign := float32(1)
	if value&1 == 0 {
		sign = -1
	}
	return int(value % uint64(dim)), sign
}

func normalize(vector []float32) {
	var sum float64
	for _, value := range vector {
		sum += float64(value * value)
	}
	if sum == 0 {
		return
	}
	length := float32(math.Sqrt(sum))
	for i := range vector {
		vector[i] = vector[i] / length
	}
}
