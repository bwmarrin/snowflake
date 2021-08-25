package snowflake

// The default epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
// You may use Config to customize this to set a different epoch for your application.
const defaultEpoch = 1288834974657

// Config can be used to specify parameters for a Node.
type Config struct {
	// Epoch from which the time-part of the ID is offset
	Epoch int64

	// NodeBits holds the number of bits to use for Node
	// Remember, you have a total 22 bits to share between Node/Step
	NodeBits uint8

	// StepBits holds the number of bits to use for Step
	// Remember, you have a total 22 bits to share between Node/Step
	StepBits uint8

	// MaxOverflowMs specifies the number of milliseconds a Node is allowed to
	// go over by
	MaxOverflowMs int64
}

var defaultConfig = Config{
	Epoch:    defaultEpoch,
	NodeBits: 10,
	StepBits: 12,
}
