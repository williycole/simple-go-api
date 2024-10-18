package cache

import (
	"fmt"
	"math/big"
	"sync"
)

// InMemCacheMap stores factorial results as big.Ints.
// It uses a map to cache calculated factorials and a RWMutex for concurrent access.
type InMemCacheMap struct {
	Mu    sync.RWMutex
	Cache map[int]*big.Int
}

// NewInMemCacheMap initializes a new cache with the base case 0! = 1.
// It returns a pointer to an InMemCacheMap with the initial cache setup.
func NewInMemCacheMap() *InMemCacheMap {
	c := map[int]*big.Int{0: big.NewInt(1)}
	fmt.Println("Cache initialized with base case 0! = 1")
	return &InMemCacheMap{
		Cache: c,
	}
}

// Set stores a factorial value in the cache for a given key.
// It locks the cache for writing to ensure thread safety.
func (s *InMemCacheMap) Set(key int, value *big.Int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Cache[key] = value
	fmt.Printf("Stored %d! = %s in cache\n", key, value.String())
}

// Get retrieves a factorial value from the cache for a given key.
// It locks the cache for reading to ensure thread safety and returns the value
// and a boolean indicating if the value was found.
func (s *InMemCacheMap) Get(key int) (*big.Int, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	val, foundFactorial := s.Cache[key]
	if foundFactorial {
		fmt.Printf("Retrieved %d! = %s from cache\n", key, val.String())
	}
	return val, foundFactorial
}
