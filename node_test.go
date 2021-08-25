package snowflake

import (
	"math"
	"testing"
	"time"
)

func TestNewNodeWithConfig(t *testing.T) {
	type args struct {
		node int64
		c    Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid config",
			args: args{
				node: 0,
				c: Config{
					NodeBits: 10,
					StepBits: 12,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid node",
			args: args{
				node: -1,
				c: Config{
					NodeBits: 10,
					StepBits: 12,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid node-bits; NodeBits should be >0",
			args: args{
				node: -1,
				c: Config{
					NodeBits: 0,
					StepBits: 12,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid step-bits; StepBits should be >0",
			args: args{
				node: -1,
				c: Config{
					NodeBits: 10,
					StepBits: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid node & step bits; NodeBits+StepBits cannot me more than 22",
			args: args{
				node: -1,
				c: Config{
					NodeBits: 10,
					StepBits: 20,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid node & step bits; NodeBits+StepBits cannot be less than 22",
			args: args{
				node: -1,
				c: Config{
					NodeBits: 13,
					StepBits: 2,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid max-overflow; it should be 0 or more",
			args: args{
				node: -1,
				c: Config{
					NodeBits:      12,
					StepBits:      10,
					MaxOverflowMs: -1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewNodeWithConfig(tt.args.node, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNodeWithConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func isTimeEqualWithTolerance(t1, t2 time.Time, tolerance time.Duration) bool {
	d := t1.Sub(t2)
	return math.Abs(float64(d)) <= math.Abs(float64(tolerance))
}

func TestNode_GenerateN(t *testing.T) {
	epoch := time.Now().UnixNano() / 1000000
	node, err := NewNodeWithConfig(0, Config{
		Epoch:         epoch,
		NodeBits:      10,
		StepBits:      12,
		MaxOverflowMs: 1,
	})

	if err != nil {
		t.Fatalf("error initialising Node; err:%v", err)
	}

	id := node.Generate()

	if id != 0 {
		t.Fatalf("ID expected to be 0")
	}

	ids, of := node.GenerateN(10000)

	if of.IsZero() {
		t.Fatal("an overflow was expected")
	}

	if !isTimeEqualWithTolerance(time.Now(), of.Time, 1*time.Millisecond) {
		t.Fatal("unexpected timestamp in overflow")
	}

	of.Time = time.Time{} // zero-out the time for comparison
	expectedOf := Overflow{
		Ms:   1,
		Step: 0xfff,
	}
	if of != expectedOf {
		t.Fatal("unexpected overflow; expected:", expectedOf, " got:", of)
	}

	expectedBlock := Block{
		First:    1,    // the next ID after the one we already generated
		N:        8191, // one less than the 2xStepSize (one step from current ms and one from the allowed overflow ms)
		NodeMask: 0x3ff << 12,
		StepMask: 0xfff,
	}
	if ids != expectedBlock {
		t.Fatal("unexpected block; expected:", expectedBlock, " got:", ids)
	}

	ids, of = node.GenerateN(1)
	if !ids.IsZero() {
		t.Fatal("no more ids expected as we are fully overflown")
	}

	<-of.AfterCleared()
	ids, of = node.GenerateN(1)
	if ids.IsZero() {
		t.Fatal("id was expected as we not overflown anymore after sleep")
	}
}

func TestNode_getMaxBlockSize(t *testing.T) {
	tests := []struct {
		name string
		node *Node
		now  int64
		want int64
	}{
		// maxOverflowMs: 0 --------------------------------
		{
			name: "New Node @0",
			node: &Node{
				time:          -1,
				step:          0,
				stepMask:      0xfff,
				maxOverflowMs: 0,
			},
			now:  0,
			want: 0xfff + 1,
		},
		{
			name: "Node{t:0 s:0} @0",
			node: &Node{
				time:          0,
				step:          0,
				stepMask:      0xfff,
				maxOverflowMs: 0,
			},
			now:  0,
			want: (0xfff + 1) - 1,
		},
		{
			name: "Node{t:0 s:101} @0",
			node: &Node{
				time:          0,
				step:          101,
				stepMask:      0xfff,
				maxOverflowMs: 0,
			},
			now:  0,
			want: 0xfff - 101,
		},
		{
			name: "Node{t:0 s:4095 (0xfff)} @0",
			node: &Node{
				time:          0,
				step:          0xfff,
				stepMask:      0xfff,
				maxOverflowMs: 0,
			},
			now:  0,
			want: 0,
		},
		{
			name: "Node{t:0 s:4095 (0xfff)} @1",
			node: &Node{
				time:          0,
				step:          0xfff,
				stepMask:      0xfff,
				maxOverflowMs: 0,
			},
			now:  1,
			want: 0xfff + 1,
		},

		// maxOverflowMs: 1 --------------------------------
		{
			name: "New Node @0",
			node: &Node{
				time:          -1,
				step:          0,
				stepMask:      0xfff,
				maxOverflowMs: 1,
			},
			now:  0,
			want: 2 * (0xfff + 1),
		},
		{
			name: "Node{t:0 s:0} @0",
			node: &Node{
				time:          0,
				step:          0,
				stepMask:      0xfff,
				maxOverflowMs: 1,
			},
			now:  0,
			want: 2*(0xfff+1) - 1,
		},
		{
			name: "Node{t:0 s:101} @0",
			node: &Node{
				time:          0,
				step:          101,
				stepMask:      0xfff,
				maxOverflowMs: 1,
			},
			now:  0,
			want: 2*(0xfff+1) - 102,
		},
		{
			name: "Node{t:0 s:4095 (0xfff)} @0",
			node: &Node{
				time:          0,
				step:          0xfff,
				stepMask:      0xfff,
				maxOverflowMs: 1,
			},
			now:  0,
			want: 0xfff + 1,
		},
		{
			name: "Node{t:0 s:4095 (0xfff)} @1",
			node: &Node{
				time:          0,
				step:          0xfff,
				stepMask:      0xfff,
				maxOverflowMs: 1,
			},
			now:  1,
			want: 2 * (0xfff + 1),
		},
		{
			name: "Node{t:1 s:101} @1",
			node: &Node{
				time:          1,
				step:          101,
				stepMask:      0xfff,
				maxOverflowMs: 1,
			},
			now:  0,
			want: (0xfff + 1) - 102,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.getMaxBlockSize(tt.now); got != tt.want {
				t.Errorf("getMaxBlockSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
