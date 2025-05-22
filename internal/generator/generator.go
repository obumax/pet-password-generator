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

// HasAny returns true if at least one of the flags is selected (except ExcludeSimilar)
func (f FlagsSet) HasAny() bool {
	return f.Upper || f.Lower || f.Digits || f.SpecSymbols
}

// Categories of symbols
const (
	ups       = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lows      = "abcdefghijklmnopqrstuvwxyz"
	digs      = "0123456789"
	specSymbs = "!@#№$;%^:&?*()-_=+[]{}<>.,/|`~"
	similars  = "il1O0"
)

// Generate creates a password length from 4 to 35 symbols
// each selected category occurs ≥1 time
// at least 1 flag must be selected (true)
func Generate(length int, flags FlagsSet) (string, error) {
	if length < 4 || length > 35 {
		return "", ErrLengthOutOfRange
	}
	// Building a pool of symbols and categories
	var pool []rune
	var required [][]rune

	if flags.Upper {
		runes := []rune(ups)
		required = append(required, runes)
		pool = append(pool, runes...)
	}
	if flags.Lower {
		runes := []rune(lows)
		required = append(required, runes)
		pool = append(pool, runes...)
	}
	if flags.Digits {
		runes := []rune(digs)
		required = append(required, runes)
		pool = append(pool, runes...)
	}
	if flags.SpecSymbols {
		runes := []rune(specSymbs)
		required = append(required, runes)
		pool = append(pool, runes...)
	}
	if len(required) == 0 {
		return "", ErrNoCategorySelected
	}
	if len(required) > length {
		return "", ErrLengthOutOfRange
	}

	// Remove similar characters if flag is set
	if flags.ExcludeSimilar {
		fp := make([]rune, 0, len(pool))
		for _, r := range pool {
			if !containsRune(similars, r) {
				fp = append(fp, r)
			}
		}
		pool = fp
		for i, cat := range required {
			fc := make([]rune, 0, len(cat))
			for _, r := range cat {
				if !containsRune(similars, r) {
					fc = append(fc, r)
				}
			}
			required[i] = fc
		}
	}

	// Password assembly
	password := make([]rune, length)

	// First, one symbol is taken from each category
	// 4 categories = 4 symbols minimum one per category
	for i, cat := range required {
		idxBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(cat))))
		if err != nil {
			return "", err
		}
		password[i] = cat[idxBig.Int64()]
	}

	// The remaining symbols are taken from the general pool
	for i := len(required); i < length; i++ {
		idxBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
		if err != nil {
			return "", err
		}
		password[i] = pool[idxBig.Int64()]
	}

	// shuffle the symbols
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
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue
		}
		j := int(jBig.Int64())
		runes[i], runes[j] = runes[j], runes[i]
	}
}
