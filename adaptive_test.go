package flashtext

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestAdaptiveCapacity(t *testing.T) {

	base := "apple word word word word "
	var sb strings.Builder
	for i := 0; i < 200; i++ {
		sb.WriteString(base)
	}
	text := sb.String()
	// Text length approx 24 * 200 = 4800 bytes
	// Matches = 200 "apple"s
	// Density = 200 / 4800 ≈ 0.0416

	kp := NewKeywordProcessor()
	kp.AddKeyWord("apple")
	kp.Build()
	defer kp.Close()

	fmt.Printf("\n=== Starting Adaptive Capacity Test (Runes: %d, Actual Matches: 200) ===\n", len([]rune(text)))
	fmt.Println("Iter\tCapacity\tMatches\tAllocated/Actual")
	fmt.Println("----\t--------\t-------\t----------------")

	for i := 1; i <= 1000; i++ {
		matches := kp.ExtractKeywords(text)

		if i == 1 || i == 5 || i == 10 || i == 50 || i == 100 || i == 500 || i == 1000 {
			cap := cap(matches)
			ratio := float64(cap) / float64(len(matches))
			fmt.Printf("%d\t%d\t\t%d\t%.2fx\n", i, cap, len(matches), ratio)
		}

		time.Sleep(1 * time.Millisecond)
	}
	fmt.Println("================================================================")
}
