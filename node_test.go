package snowflake

import (
	"testing"
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
