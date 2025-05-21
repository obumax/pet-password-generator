package generator

import "testing"

func TestGenerate_Basic(t *testing.T) {
	flags := FlagsSet{Upper: true, Lower: true, Digits: true, SpecSymbols: true}
	password, err := Generate(12, flags)
	if err != nil {
		t.Fatal(err)
	}
	if len(password) != 12 {
		t.Fatalf("Ожидаемая длина 12 символов, полученная %d символов", len(password))
	}
}
