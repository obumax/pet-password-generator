package generator

import (
	"testing"
	"unicode"
)

func TestGenerate_Basic(t *testing.T) {
	flags := FlagsSet{Upper: true, Lower: true, Digits: true, SpecSymbols: true}
	pwd, err := Generate(12, flags)
	if err != nil {
		t.Fatal(err)
	}
	if len(pwd) != 12 {
		t.Fatalf("ожидаемая длина 12, полученная длина %d", len(pwd))
	}
	var hasU, hasL, hasD, hasS bool
	for _, r := range pwd {
		switch {
		case unicode.IsUpper(r):
			hasU = true
		case unicode.IsLower(r):
			hasL = true
		case unicode.IsDigit(r):
			hasD = true
		default:
			hasS = true
		}
	}
	if !hasU || !hasL || !hasD || !hasS {
		t.Fatalf("недостаточно категорий: U=%v L=%v D=%v S=%v", hasU, hasL, hasD, hasS)
	}
}

func TestGenerate_ExcludeSimilar(t *testing.T) {
	flags := FlagsSet{Upper: true, Lower: true, Digits: true, SpecSymbols: true, ExcludeSimilar: true}
	pwd, err := Generate(20, flags)
	if err != nil {
		t.Fatal(err)
	}
	bad := map[rune]bool{'i': true, 'l': true, '1': true, 'O': true, '0': true}
	for _, r := range pwd {
		if bad[r] {
			t.Fatalf("символ %q не должен присутствовать", r)
		}
	}
}

func TestGenerate_Errors(t *testing.T) {
	if _, err := Generate(2, FlagsSet{Upper: true}); err != ErrLengthOutOfRange {
		t.Fatalf("ожидаемая ошибка ErrLengthOutOfRange, получено %v", err)
	}
	if _, err := Generate(5, FlagsSet{}); err != ErrNoCategorySelected {
		t.Fatalf("ожидаемая ошибка ErrNoCategorySelected, получено %v", err)
	}
}
