package utils

import (
	"math/rand"
	"time"
)

// Alphabet constant containing all characters that can be used for generating random strings
const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randomGenerator struct encapsulates the random generation logic
type randomGenerator struct {
	rand *rand.Rand
}

// NewRandomGenerator initializes a new instance of randomGenerator with a seeded rand
func NewRandomGenerator() *randomGenerator {
	return &randomGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// RandomInt generates a random integer between min (inclusive) and max (inclusive)
func (rg *randomGenerator) RandomInt(min, max int64) int64 {
	return min + rg.rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n using characters from the alphabet
func (rg *randomGenerator) RandomString(n int) string {
	result := make([]byte, n)
	for i := range result {
		result[i] = alphabet[rg.rand.Intn(len(alphabet))]
	}
	return string(result)
}

// RandomOwner generates a random owner string of length 6
func (rg *randomGenerator) RandomOwner() string {
	return rg.RandomString(6)
}

// RandomMoney generates a random money amount between 0 and 1000
func (rg *randomGenerator) RandomMoney() int64 {
	return rg.RandomInt(0, 1000)
}

// RandomCurrency generates a random currency from a list of provided currencies
func (rg *randomGenerator) RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	if n == 0 {
		return "" // Return empty string if no currencies provided
	}
	return currencies[rg.rand.Intn(n)]
}
