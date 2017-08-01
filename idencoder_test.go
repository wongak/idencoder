package idencoder_test

import (
	"testing"

	"github.com/wongak/idencoder"
)

func TestObfuscatorWithInvalidP(t *testing.T) {
	_, err := idencoder.NewObfuscator(24, 231)
	if err == nil {
		t.Error("expect error, got nil")
	}
}

func TestObfuscator(t *testing.T) {
	p := int64(32416189079)
	xor := int64(123)
	id := int64(12341441)

	o, err := idencoder.NewObfuscator(p, xor)
	if err != nil {
		t.Fatal(err)
	}
	hidden := o.Hide(id)
	t.Logf("hidden: %d", hidden)
	if o.Show(hidden) != id {
		t.Errorf("expect revealed %d to be original %d", o.Show(hidden), id)
	}
}

func TestAlphabetBinary(t *testing.T) {
	n := int64(123)
	expect := "yyyyxyy"
	enc, err := idencoder.NewAlphabetEncoder("xy")
	if err != nil {
		t.Fatal(err)
	}
	if enc.FromBase10(n) != expect {
		t.Errorf("expect encoded %s, got %s", expect, enc.FromBase10(n))
	}
	dec, err := enc.ToBase10(enc.FromBase10(n))
	if err != nil {
		t.Fatal(err)
	}
	if dec != n {
		t.Errorf("expect decoded to be %d, got %d", n, dec)
	}
}

func TestAlphabet(t *testing.T) {
	enc, err := idencoder.NewAlphabetEncoder("abcdefghkmpqrswxyzABCDEFGHKMPQRSWXYZ23456789")
	if err != nil {
		t.Fatal(err)
	}
	n := int64(131414341)
	encoded := enc.FromBase10(n)
	t.Logf("%d = %s", n, encoded)
	dec, err := enc.ToBase10(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if dec != n {
		t.Errorf("expect decoded to be %d, got %d", n, dec)
	}
}
