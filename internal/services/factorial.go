package services

import (
	"fmt"
	"math/big"
	"simple-go-api/internal/cache"
)

// CalculateFactorial calculates the factorial of n using memoization/cache.
func CalculateFactorial(n int, fc *cache.InMemCacheMap) *big.Int {
	// Lock the mutex for reading
	fc.Mu.RLock()
	result, found := fc.Cache[n]
	fc.Mu.RUnlock()

	if found {
		fmt.Printf("Cache hit for %d! = %s\n", n, result.String())
		// Return a copy of the cached value
		return new(big.Int).Set(result)
	}

	// Calculate the factorial
	fmt.Printf("Calculating %d!\n", n)
	result = big.NewInt(1)
	if n > 1 {
		result.Mul(big.NewInt(int64(n)), CalculateFactorial(n-1, fc))
	}

	// Lock the mutex for writing
	fc.Mu.Lock()
	// Store a copy in the cache
	fc.Cache[n] = new(big.Int).Set(result)
	fc.Mu.Unlock()

	fmt.Printf("Stored %d! = %s in cache\n", n, result.String())
	return result
}
