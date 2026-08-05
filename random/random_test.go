package random

import "testing"

func TestSameSeedSameSequence(t *testing.T) {
	a, b := New(12345), New(12345)
	for i := 0; i < 1000; i++ {
		if x, y := a.Next(), b.Next(); x != y {
			t.Fatalf("draw %d diverged: %d vs %d", i, x, y)
		}
	}
}

func TestDifferentSeedsDiverge(t *testing.T) {
	a, b := New(1), New(2)
	same := 0
	for i := 0; i < 100; i++ {
		if a.Next() == b.Next() {
			same++
		}
	}
	if same > 1 {
		t.Errorf("distinct seeds produced %d identical draws in 100", same)
	}
}

func TestZeroSeedDoesNotCollapse(t *testing.T) {
	// xorshift32 is stuck at zero forever if seeded with zero.
	r := New(0)
	for i := 0; i < 10; i++ {
		if r.Next() == 0 {
			t.Fatal("zero seed collapsed the sequence to zero")
		}
	}
}

func TestBelowStaysInRange(t *testing.T) {
	r := New(99)
	for i := 0; i < 5000; i++ {
		if v := r.Below(10); v >= 10 {
			t.Fatalf("Below(10) returned %d", v)
		}
	}
	if v := r.Below(0); v != 0 {
		t.Errorf("Below(0) should be 0, got %d", v)
	}
}

func TestRangeStaysInBounds(t *testing.T) {
	r := New(7)
	for i := 0; i < 5000; i++ {
		if v := r.Range(-5, 5); v < -5 || v >= 5 {
			t.Fatalf("Range(-5, 5) returned %d", v)
		}
	}
	if v := r.Range(3, 3); v != 3 {
		t.Errorf("empty range should return lo, got %d", v)
	}
}

func TestShuffleIsDeterministicAndKeepsElements(t *testing.T) {
	build := func() []int { return []int{0, 1, 2, 3, 4, 5, 6, 7} }

	a, b := build(), build()
	New(42).Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	New(42).Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })

	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("same seed shuffled differently at %d: %v vs %v", i, a, b)
		}
	}

	seen := make(map[int]bool)
	for _, v := range a {
		seen[v] = true
	}
	if len(seen) != 8 {
		t.Errorf("shuffle lost or duplicated elements: %v", a)
	}
}
