package generator

import (
	"crypto/rand"
	"errors"
	"math/big"
)

var (
	ErrLengthOutOfRange   = errors.New("length_out_of_range")
	ErrNoCategorySelected = errors.New("no_category_selected")
)

type FlagsSet struct {
	Upper, Lower, Digits, SpecSymbols, ExcludeSimilar bool
}

// Categories of symbols
const (
	ups       = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lows      = "abcdefghijklmnopqrstuvwxyz"
	digs      = "0123456789"
	specSymbs = "!@#№$;%^:&?*()-_=+[]{}<>.,/|`~"
	similars  = "il1O0"
)

// Generate creates a password length from 4 to 35 characters
// each selected category occurs ≥1 time
// at least 1 flag must be selected (true)
func Generate(lenght int, flags FlagsSet) (string, error) {
	if lenght < 4 || lenght > 35 {
		return "", ErrLengthOutOfRange
	}
	// Building a pool of symbols and categories
	var pool []rune
	var required [][]rune

	if flags.Upper {
		required = append(required, []rune(ups))
		pool = append(pool, []rune(ups)...)
	}
	if flags.Lower {
		required = append(required, []rune(lows))
		pool = append(pool, []rune(lows)...)
	}
	if flags.Digits {
		required = append(required, []rune(digs))
		pool = append(pool, []rune(digs)...)
	}
	if flags.SpecSymbols {
		required = append(required, []rune(specSymbs))
		pool = append(pool, []rune(specSymbs)...)
	}
	if len(required) == 0 {
		return "", ErrNoCategorySelected
	}

	// Remove similar characters if flag is set
	if flags.ExcludeSimilar {
		filtered := pool[:0]
		for _, r := range pool {
			if !containsRune(similars, r) {
				filtered = append(filtered, r)
			}
		}
		pool = filtered
	}

	// Password assembly
	password := make([]rune, lenght)

	// First, one symbol is taken from each category
	// 4 categories = 4 symbols minimum one per category
	for i, cat := range required {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(cat))))
		password[i] = cat[index.Int64()]
	}

	// The remaining symbols are taken from the general pool
	for i := len(required); i < lenght; i++ {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
		password[i] = pool[index.Int64()]
	}

	// Mix the symbols
	mix(password)
	return string(password), nil
}

// containsRune compares symbols to remove similars
func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

// mix shuffles symbols
func mix(runes []rune) {
	n := len(runes)
	for i := n - 1; i > 0; i-- {
		jBig, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(jBig.Int64())
		runes[i], runes[j] = runes[j], runes[i]
	}
}
