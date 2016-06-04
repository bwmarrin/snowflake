package snowflake

import "testing"

func BenchmarkGenerate(b *testing.B) {

	node, _ := NewNode(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = node.Generate()
	}
}

func BenchmarkUnmarshal(b *testing.B) {

	node, _ := NewNode(1)
	id, _ := node.Generate()
	var id2 ID

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = id2.UnmarshalJSON(id.Bytes())
	}
}

func BenchmarkMarshal(b *testing.B) {

	node, _ := NewNode(1)
	id, _ := node.Generate()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = id.MarshalJSON()
	}
}
