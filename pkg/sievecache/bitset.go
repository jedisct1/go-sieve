package sievecache

// BitSet provides a memory-efficient way to store boolean values
// using 1 bit per value instead of 1 byte per value.
type BitSet struct {
	bits []uint64
	size int
}

// NewBitSet creates a new bit set with the given initial capacity.
func NewBitSet(capacity int) *BitSet {
	// Calculate how many uint64s we need to store capacity bits
	numWords := (capacity + 63) / 64
	return &BitSet{
		bits: make([]uint64, numWords),
		size: capacity,
	}
}

// Set sets the bit at the given index to the specified value.
func (b *BitSet) Set(index int, value bool) {
	if index >= b.size {
		b.resize(index + 1)
	}

	wordIndex := index / 64
	bitIndex := index % 64

	if value {
		b.bits[wordIndex] |= 1 << bitIndex
	} else {
		b.bits[wordIndex] &= ^(1 << bitIndex)
	}
}

// Get returns the value of the bit at the given index.
func (b *BitSet) Get(index int) bool {
	if index >= b.size {
		return false
	}

	wordIndex := index / 64
	bitIndex := index % 64

	return (b.bits[wordIndex] & (1 << bitIndex)) != 0
}

// Resize increases the capacity of the bit set to at least the specified size.
func (b *BitSet) resize(newSize int) {
	if newSize <= b.size {
		return
	}

	// Calculate new number of words needed
	numWords := (newSize + 63) / 64

	// If we need more words, extend the slice
	if numWords > len(b.bits) {
		newBits := make([]uint64, numWords)
		copy(newBits, b.bits)
		b.bits = newBits
	}

	b.size = newSize
}

// Append adds a new bit to the end of the set.
func (b *BitSet) Append(value bool) {
	b.Set(b.size, value)
}

// Truncate reduces the size of the bit set to the specified size.
func (b *BitSet) Truncate(newSize int) {
	if newSize >= b.size {
		return
	}

	// Calculate new number of words needed
	numWords := (newSize + 63) / 64

	// Clear any bits in the last word that are beyond the new size
	if numWords > 0 {
		lastWordBits := newSize % 64
		if lastWordBits > 0 {
			// Create a mask for the bits we want to keep
			mask := (uint64(1) << lastWordBits) - 1
			// Apply the mask to the last word
			b.bits[numWords-1] &= mask
		}
	}

	// If we need fewer words, truncate the slice
	if numWords < len(b.bits) {
		b.bits = b.bits[:numWords]
	}

	b.size = newSize
}

// Size returns the number of bits in the set.
func (b *BitSet) Size() int {
	return b.size
}

// CountSetBits returns the number of bits that are set to true.
func (b *BitSet) CountSetBits() int {
	count := 0
	for _, word := range b.bits {
		// Count the bits in this word using Hamming weight
		for word != 0 {
			count++
			word &= word - 1 // Clear the least significant set bit
		}
	}
	return count
}
