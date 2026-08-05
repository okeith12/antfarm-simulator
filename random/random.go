// Package random provides the simulation's deterministic pseudo-random source.
//
// The algorithm is xorshift32, chosen because it is trivial to reimplement
// identically in C++ for the firmware port. The same seed must always produce
// the same colony in both languages, which is what makes the Go-to-C++ parity
// gate possible: run both builds with one seed, diff the tick dumps, and expect
// no differences.
//
// Nothing here uses math/rand. The simulation must never call the global
// generator, because an unseeded run cannot be reproduced or compared.
package random

// Generator is a deterministic xorshift32 generator.
// It is not safe for concurrent use; the simulation is single-threaded.
type Generator struct {
	state uint32
}

// New creates a generator from the given seed.
// A zero seed would make xorshift32 emit only zeroes, so it is replaced with
// the golden-ratio constant that the firmware spec uses for the same reason.
func New(seed uint32) *Generator {
	if seed == 0 {
		seed = 0x9E3779B9
	}
	return &Generator{state: seed}
}

// Next returns the next raw 32-bit value in the sequence.
func (r *Generator) Next() uint32 {
	x := r.state
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	r.state = x
	return x
}

// Below returns a value in [0, n). It returns 0 when n is 0.
//
// This uses plain modulo, which is very slightly biased toward smaller values.
// That is deliberate: the reduction has to be reproducible bit-for-bit in C++,
// and a rejection-sampling scheme would consume a variable number of draws and
// desynchronise the two sequences.
func (r *Generator) Below(n uint32) uint32 {
	if n == 0 {
		return 0
	}
	return r.Next() % n
}

// Range returns a value in [lo, hi). It returns lo when hi is not above lo.
func (r *Generator) Range(lo, hi int32) int32 {
	if hi <= lo {
		return lo
	}
	return lo + int32(r.Below(uint32(hi-lo)))
}

// Chance reports whether a percent-in-100 roll succeeds.
// Chance(10) is true one time in ten.
func (r *Generator) Chance(percent uint32) bool {
	return r.Below(100) < percent
}

// Shuffle permutes n elements using the supplied swap function.
//
// This is Fisher-Yates walking downward, matching the order the C++ port will
// use so that a shuffled slice comes out identical in both implementations.
func (r *Generator) Shuffle(n int, swap func(i, j int)) {
	for i := n - 1; i > 0; i-- {
		j := int(r.Below(uint32(i + 1)))
		swap(i, j)
	}
}
