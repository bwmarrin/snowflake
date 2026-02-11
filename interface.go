package snowflake

// NodeIdAllocator node id allocator Issue 1: NodeId collision prevention
type NodeIdAllocator interface {
	// Alloc allocates node ID
	Alloc() (nodeId int64, err error)
	// Drift node ID
	Drift(nodeId int64) (newNodeId int64, err error)
}

// TimeSynchronizer time synchronizer Issue 2: Prevent large clock rollback when restarting snowflake ID
type TimeSynchronizer interface {
	// Async synchronizes time This method does not report errors, reducing external component dependencies and increasing robustness
	Async(time int64)
}
