package snowflake

type Option struct {
	// node id allocator
	allocator NodeIdAllocator
	// time synchronizer
	synchronizer TimeSynchronizer
}

// OptionFn option function
type OptionFn func(op *Option)

// WithNodeIdAllocator sets node id allocator
func WithNodeIdAllocator(allocator NodeIdAllocator) func(op *Option) {
	return func(op *Option) {
		op.allocator = allocator
	}
}

// WithTimeSynchronizer sets time synchronizer
func WithTimeSynchronizer(synchronizer TimeSynchronizer) func(op *Option) {
	return func(op *Option) {
		op.synchronizer = synchronizer
	}
}
