package services

import (
	"math/big"
	"simple-go-api/internal/cache"
	"testing"
)

func TestCalculateFactorial(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want *big.Int
	}{
		{
			name: "Factorial of 0",
			n:    0,
			want: big.NewInt(1),
		},
		{
			name: "Factorial of 1",
			n:    1,
			want: big.NewInt(1),
		},
		{
			name: "Factorial of 5",
			n:    5,
			want: big.NewInt(120),
		},
		{
			name: "Factorial of 10",
			n:    10,
			want: big.NewInt(3628800),
		},
		{
			name: "Factorial of 100",
			n:    100,
			want: func() *big.Int {
				val, success := new(big.Int).SetString("93326215443944152681699238856266700490715968264381621468592963895217599993229915608941463976156518286253697920827223758251185210916864000000000000000000000000", 10)
				if !success {
					t.Fatalf("Failed to set big.Int from string")
				}
				return val
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize a new cache for each test
			fc := cache.NewInMemCacheMap()

			// Calculate factorial
			got := CalculateFactorial(tt.n, fc)

			// Compare with expected result using big.Int.Cmp
			if got.Cmp(tt.want) != 0 {
				t.Errorf("CalculateFactorial() = %v, want %v", got, tt.want)
			}
		})
	}
}
