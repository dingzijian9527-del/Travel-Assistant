package rag

// Document describes one travel knowledge snippet.
type Document struct {
	ID      string
	Title   string
	City    string
	Tags    []string
	Content string
}

// Result is one retrieval hit.
type Result struct {
	Document Document
	Score    float64
}
