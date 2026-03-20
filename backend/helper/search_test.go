// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var texts = []string{
	"go-testify",
	"go-testify",
	"go-testify",
	"0apache",
	"bsd 2.1",
	"bsd2.1",
	"lala",
	"whatever test lala test",
	"(test and apache-2.0)",
	"(test and apache)",
	"test1234",
	"The quick brown fox jumps over the lazy dog",
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit",
}

var (
	once       sync.Once
	seededRand = big.NewInt(int64(len(charset)))
)

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		number, err := rand.Int(rand.Reader, seededRand)
		if err != nil {
			return ""
		}
		b[i] = charset[number.Int64()]
	}
	return string(b)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func search(arr []string, search string, exact bool) []string {
	return Search(arr, search, exact)
}

// finds (test and apache) and (test and apache-2.0)
func TestSearchBrackets(t *testing.T) {
	createIndex(t)
	searchResult := search(texts, "(test and apache)", false)
	assert.Len(t, searchResult, 1)
}

func TestSearchEscapedBrackets(t *testing.T) {
	createIndex(t)
	searchResult := search(texts, "\"(test and apache)\"", false)
	assert.Len(t, searchResult, 0)
}

func TestVersion(t *testing.T) {
	createIndex(t)
	searchResult := search(texts, "2.1", false)
	assert.Len(t, searchResult, 2)
}

func TestBsd(t *testing.T) {
	createIndex(t)
	searchResult := search(texts, "bsd 2.1", false)
	assert.Len(t, searchResult, 1)
}

func TestBsd21Exact(t *testing.T) {
	createIndex(t)
	searchResult := search(texts, "bsd 2.1", true)
	assert.Len(t, searchResult, 1)
}

func TestBsdExact(t *testing.T) {
	createIndex(t)
	searchResult := search(texts, "bsd", true)
	assert.Empty(t, searchResult)
}

func TestBsd21(t *testing.T) {
	createIndex(t)
	searchResult := search(texts, "bsd2.1", false)
	assert.Len(t, searchResult, 1)
}

func TestMiultipleSearch(t *testing.T) {
	createIndex(t)
	searchResult := search(texts, "lala test", false)
	assert.Len(t, searchResult, 1)
}

func TestSimpleAnotherSearch(t *testing.T) {
	createIndex(t)
	searchResult := search(texts, "go-testify", false)
	assert.Len(t, searchResult, 3)
}

func createIndex(t *testing.T) {
	once.Do(func() {
		for i := 0; i < 1000000; i++ {
			randomString := generateRandomString(10)
			texts = append(texts, randomString)
		}
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("Alloc: %d mb\n", bToMb(m.Alloc))
		assert.Less(t, bToMb(m.Alloc), uint64(100))
	})
}
