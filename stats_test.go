package flashtext

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestStats_Init(t *testing.T) {
	ctx := context.Background()
	s := newStats(ctx, 0.2, 512)
	defer s.close()

	if s == nil {
		t.Fatal("stats should not be nil")
	}
	if s.alpha != 0.2 {
		t.Errorf("expected alpha 0.2, got %v", s.alpha)
	}
	if density := s.getDensity(); density != 0.01 {
		t.Errorf("expected density 0.01, got %v", density)
	}
}

func TestStats_Integration_With_Processor(t *testing.T) {
	kp := NewKeywordProcessor()
	defer kp.Close()

	// 1. Initial density should be 0.01
	if density := kp.stats.getDensity(); density != 0.01 {
		t.Errorf("expected initial density 0.01, got %v", density)
	}

	// 2. Add keywords and process text
	kp.AddKeyWord("apple")
	kp.AddKeyWord("banana")
	kp.Build()

	// Text with 50% density (10 chars total, "apple" is 5 chars)
	text := "apple....." // "apple" (5) + "....." (5) = 10 chars
	matches := kp.ExtractKeywords(text)
	if len(matches) != 1 {
		t.Errorf("expected 1 match, got %d", len(matches))
	}

	// Allow time for async update
	time.Sleep(50 * time.Millisecond)

	// 3. Verify density updated
	// New density = 1 match / 10 runes = 0.1
	// EWMA = 0.2 * 0.1 + 0.8 * 0.0 = 0.02
	currentDensity := kp.stats.getDensity()
	if currentDensity >= 0.1 {
		t.Errorf("expected density < 0.1, got %v", currentDensity)
	}
	if math.Abs(currentDensity-0.02) > 0.01 {
		t.Errorf("expected density ~0.02, got %v", currentDensity)
	}
}

func TestStats_EWMA_Logic(t *testing.T) {
	ctx := context.Background()
	// Use alpha = 0.5 for easier calculation
	s := newStats(ctx, 0.5, 512)
	defer s.close()

	// Update 1: matches=10, runes=100 -> density=0.1
	// EWMA = 0.5 * 0.1 + 0.5 * 0.0 = 0.05
	s.add(10, 100)
	time.Sleep(50 * time.Millisecond)
	if density := s.getDensity(); math.Abs(density-0.05) > 0.01 {
		t.Errorf("expected density ~0.05, got %v", density)
	}

	// Update 2: matches=20, runes=100 -> density=0.2
	// EWMA = 0.5 * 0.2 + 0.5 * 0.05 = 0.1 + 0.025 = 0.125
	s.add(20, 100)
	time.Sleep(50 * time.Millisecond)
	if density := s.getDensity(); math.Abs(density-0.125) > 0.01 {
		t.Errorf("expected density ~0.125, got %v", density)
	}
}

func TestStats_Concurrency(t *testing.T) {
	ctx := context.Background()
	s := newStats(ctx, 0.2, 512)
	defer s.close()

	var wg sync.WaitGroup
	count := 100

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			s.add(1, 10)
		}()
	}
	wg.Wait()

	// Give processor time to consume channel
	time.Sleep(100 * time.Millisecond)

	// Simply verify it didn't panic and density increased
	if density := s.getDensity(); density <= 0.0 {
		t.Errorf("expected density > 0.0, got %v", density)
	}
}

func TestStats_ZeroRunes(t *testing.T) {
	ctx := context.Background()
	s := newStats(ctx, 0.2, 512)
	defer s.close()

	// Should not panic or update (division by zero protection)
	initialDensity := s.getDensity()
	s.add(1, 0)
	time.Sleep(20 * time.Millisecond)
	if density := s.getDensity(); density != initialDensity {
		t.Errorf("expected density to remain %v, got %v", initialDensity, density)
	}
}

func TestFloat64BitsConversion(t *testing.T) {
	val := 0.123456
	bits := float64ToBits(val)
	val2 := float64FromBits(bits)
	if val != val2 {
		t.Errorf("expected %v, got %v", val, val2)
	}
}

func TestAtomicLogicProtection(t *testing.T) {
	// This test verifies that we're using math.Float64bits correctly
	// because atomic functions work on uint64, not float64.
	// We want to ensure the conversion doesn't lose precision or scramble data.

	testVal := 0.5
	testBits := math.Float64bits(testVal)

	var atomicItem uint64
	// Simulate CAS
	swapped := atomic.CompareAndSwapUint64(&atomicItem, 0, testBits)
	if !swapped {
		t.Fatal("expected swap to succeed")
	}

	loadedBits := atomic.LoadUint64(&atomicItem)
	loadedVal := math.Float64frombits(loadedBits)

	if loadedVal != testVal {
		t.Errorf("expected %v, got %v", testVal, loadedVal)
	}
}
