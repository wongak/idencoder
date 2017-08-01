package idencoder

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"strings"
)

// Obfuscator can hide and show integer IDs
//
// The idea is using a modular multiplicative inverse of the ID
// using a secret prime.
//
// Additionally a simple XOR is applied with another chosen secret.
type Obfuscator struct {
	max   *big.Int
	prime *big.Int
	inv   *big.Int
	xor   *big.Int
}

// NewObfuscator creates a new obfuscator
func NewObfuscator(p, xor int64) (*Obfuscator, error) {
	o := &Obfuscator{}
	o.max = big.NewInt(math.MaxInt64)
	o.prime = big.NewInt(p)
	o.inv = new(big.Int)
	o.xor = big.NewInt(xor)

	g := new(big.Int)
	g.GCD(o.inv, nil, o.prime, new(big.Int).Add(o.max, big.NewInt(1)))
	if g.Int64() != 1 {
		return nil, fmt.Errorf("invalid p")
	}
	o.inv.Mod(o.inv, new(big.Int).Add(o.max, big.NewInt(1)))
	return o, nil
}

// Hide returns the hidden representation of i
func (o *Obfuscator) Hide(i int64) int64 {
	hidden := big.NewInt(i)
	hidden.Mul(hidden, o.prime)
	hidden.And(hidden, o.max)
	hidden.Xor(hidden, o.xor)
	return hidden.Int64()
}

// Show reveals the hidden ID
func (o *Obfuscator) Show(i int64) int64 {
	shown := big.NewInt(i)
	shown.Xor(shown, o.xor)
	shown.Mul(shown, o.inv)
	shown.And(shown, o.max)
	return shown.Int64()
}

// AlphabetEncoder takes an alphabet and encodes/decodes any integer into the
// representation in the given alphabet
type AlphabetEncoder struct {
	alphabet string
	base     int64
}

// NewAlphabetEncoder creates a new alphabet encoder in base len(alphabet)
//
// base must be at least 2
func NewAlphabetEncoder(a string) (*AlphabetEncoder, error) {
	if len(a) <= 1 {
		return nil, fmt.Errorf("invalid alphabet. require at least base 2")
	}
	return &AlphabetEncoder{
		alphabet: a,
		base:     int64(len(a)),
	}, nil
}

// FromBase10 encodes a base 10 integer to its alphabet representation
func (a *AlphabetEncoder) FromBase10(i int64) string {
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	if i == 0 {
		return string(a.alphabet[0])
	}
	for i > 0 {
		rem := i % a.base
		buf.WriteByte(a.alphabet[rem])
		i = i / a.base
	}
	// reverse
	bs := buf.Bytes()
	a.reverse(bs)
	return string(bs)
}

func (a *AlphabetEncoder) reverse(bs []byte) {
	for i := len(bs)/2 - 1; i >= 0; i-- {
		opp := len(bs) - 1 - i
		bs[i], bs[opp] = bs[opp], bs[i]
	}
}

// ToBase10 reverses the encoding
func (a *AlphabetEncoder) ToBase10(s string) (int64, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("invalid empty string")
	}
	bs := []byte(s)
	a.reverse(bs)
	var n int64
	for i, c := range bs {
		ix := strings.IndexByte(a.alphabet, c)
		if ix == -1 {
			return 0, fmt.Errorf("invalid non-alphabet char: %c", c)
		}
		n += (int64(ix) * pow(a.base, int64(i)))
	}
	return n, nil
}

func pow(a, b int64) int64 {
	p := int64(1)
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}
