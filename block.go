package snowflake

// Block represents a continuous block of IDs
type Block struct {
	First    ID
	N        int64
	NodeMask int64
	StepMask int64
}

// IsZero evaluates to true if the block is empty or invalid
func (b Block) IsZero() bool {
	return b.N == 0
}

// BlockIterator is a utility to iterate over a Block
type BlockIterator struct {
	Block
	n    int64
	time int64
	node int64
	step int64
}

// NewBlockIterator creates a BlockIterator for the given Block
func NewBlockIterator(b Block) *BlockIterator {
	stepBits := 0
	for i := b.StepMask; i > 0; i >>= 1 {
		stepBits++
	}

	nodeBits := 0
	for i := b.NodeMask >> stepBits; i > 0; i >>= 1 {
		nodeBits++
	}

	return &BlockIterator{
		Block: b,
		n:     -1,
		time:  int64(b.First) & ^(b.NodeMask | b.StepMask),
		node:  int64(b.First) & b.NodeMask,
		step:  int64(b.First) & b.StepMask,
	}
}

// Next returns the next ID from a block being iterated over.
// When iteration reaches the end of the block, ok will be false.
func (i *BlockIterator) Next() (id ID, ok bool) {
	if i.Done() {
		return -1, false
	}

	if i.n < 0 {
		i.n = 1
		return i.First, true
	}

	i.n++
	i.step = (i.step + 1) & i.StepMask
	if i.step == 0 {
		i.time += (i.NodeMask << 1) & ^(i.NodeMask | i.StepMask)
	}

	return ID(i.time | i.node | i.step), true
}

// Done returns true if iteration has reached the end of the block.
func (i BlockIterator) Done() bool {
	return i.n >= i.N
}
