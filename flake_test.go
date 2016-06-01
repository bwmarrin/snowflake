package flake

import (
	"testing"
)

func BenchmarkGenerate(b *testing.B) {

	node, _ := NewNode(1)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, _ = node.GenerateNoSleep()
	}

}
