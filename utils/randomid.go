package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomID generates a random uint ID within a given range
func GenerateRandomID(min, max uint) uint {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	return min + uint(rand.Uint32())%uint(max-min+1) // Generate random uint within the range [min, max]
}