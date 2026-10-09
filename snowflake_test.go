package snowflake

import (
	"bytes"
	"reflect"
	"testing"
)

//******************************************************************************
// General Test funcs

func TestNewNode(t *testing.T) {

	_, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	_, err = NewNode(5000)
	if err == nil {
		t.Fatalf("no error creating NewNode, %s", err)
	}

}

// lazy check if Generate will create duplicate IDs
// would be good to later enhance this with more smarts
func TestGenerateDuplicateID(t *testing.T) {

	node, _ := NewNode(1)

	var x, y ID
	for i := 0; i < 1000000; i++ {
		y = node.Generate()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

// I feel like there's probably a better way
func TestRace(t *testing.T) {

	node, _ := NewNode(1)

	go func() {
		for i := 0; i < 1000000000; i++ {

			NewNode(1)
		}
	}()

	for i := 0; i < 4000; i++ {

		node.Generate()
	}

}

//******************************************************************************
// Converters/Parsers Test funcs
// We should have funcs here to test conversion both ways for everything

func TestPrintAll(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	id := node.Generate()

	t.Logf("Int64    : %#v", id.Int64())
	t.Logf("String   : %#v", id.String())
	t.Logf("Base2    : %#v", id.Base2())
	t.Logf("Base32   : %#v", id.Base32())
	t.Logf("Base36   : %#v", id.Base36())
	t.Logf("Base58   : %#v", id.Base58())
	t.Logf("Base64   : %#v", id.Base64())
	t.Logf("IntBase64: %#v", id.IntBase64())
	t.Logf("Bytes    : %#v", id.Bytes())
	t.Logf("IntBytes : %#v", id.IntBytes())

}

func TestInt64(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.Int64()

	pID := ParseInt64(i)
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	mi := int64(1116766490855473152)
	pID = ParseInt64(mi)
	if pID.Int64() != mi {
		t.Fatalf("pID %v != mi %v", pID.Int64(), mi)
	}

}

func TestString(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	si := oID.String()

	pID, err := ParseString(si)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `1116766490855473152`
	_, err = ParseString(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `1112316766490855473152`
	_, err = ParseString(ms)
	if err == nil {
		t.Fatalf("no error parsing %s", ms)
	}
}

func TestBase2(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.Base2()

	pID, err := ParseBase2(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `111101111111101110110101100101001000000000000000000000000000`
	_, err = ParseBase2(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `1112316766490855473152`
	_, err = ParseBase2(ms)
	if err == nil {
		t.Fatalf("no error parsing %s", ms)
	}
}

func TestBase32(t *testing.T) {

	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	for i := 0; i < 100; i++ {

		sf := node.Generate()
		b32i := sf.Base32()
		psf, err := ParseBase32([]byte(b32i))
		if err != nil {
			t.Fatal(err)
		}
		if sf != psf {
			t.Fatal("Parsed does not match String.")
		}
	}
}

func TestBase36(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.Base36()

	pID, err := ParseBase36(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `8hgmw4blvlkw`
	_, err = ParseBase36(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `68h5gmw443blv2lk1w`
	_, err = ParseBase36(ms)
	if err == nil {
		t.Fatalf("no error parsing, %s", err)
	}
}

func TestBase58(t *testing.T) {

	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	for i := 0; i < 10; i++ {

		sf := node.Generate()
		b58 := sf.Base58()
		psf, err := ParseBase58([]byte(b58))
		if err != nil {
			t.Fatal(err)
		}
		if sf != psf {
			t.Fatal("Parsed does not match String.")
		}
	}
}

func TestBase64(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	for i := 0; i < 100; i++ {
		oID := node.Generate()
		b64 := oID.Base64()

		if len(b64) != 12 {
			t.Fatalf("expected Base64 length of 12 for 8 bytes, got %d (%s)", len(b64), b64)
		}

		pID, err := ParseBase64(b64)
		if err != nil {
			t.Fatalf("error parsing, %s", err)
		}
		if pID != oID {
			t.Fatalf("pID %v != oID %v", pID, oID)
		}
	}

	// Legacy string-based Base64 encoding of "1116819494660997120"
	ms := `MTExNjgxOTQ5NDY2MDk5NzEyMA==`
	pID, err := ParseBase64(ms)
	if err != nil {
		t.Fatalf("error parsing legacy base64, %s", err)
	}
	if pID.Int64() != 1116819494660997120 {
		t.Fatalf("legacy pID %v != expected 1116819494660997120", pID.Int64())
	}

	// Specific known 8-byte binary Base64 value
	knownID := ID(1116823421972381696)
	expectedB64 := "D3/A/C+AAAA="
	if knownID.Base64() != expectedB64 {
		t.Fatalf("knownID.Base64() %s != %s", knownID.Base64(), expectedB64)
	}
	pID, err = ParseBase64(expectedB64)
	if err != nil {
		t.Fatalf("error parsing knownID base64, %s", err)
	}
	if pID != knownID {
		t.Fatalf("pID %v != knownID %v", pID, knownID)
	}

	// Zero ID
	zeroID := ID(0)
	zeroB64 := "AAAAAAAAAAA="
	if zeroID.Base64() != zeroB64 {
		t.Fatalf("zeroID.Base64() %s != %s", zeroID.Base64(), zeroB64)
	}
	pID, err = ParseBase64(zeroB64)
	if err != nil {
		t.Fatalf("error parsing zeroID base64, %s", err)
	}
	if pID != zeroID {
		t.Fatalf("pID %v != zeroID %v", pID, zeroID)
	}

	// Invalid Base64 syntax (missing padding)
	ms = `MTExNjgxOTQ5NDY2MDk5NzEyMA`
	_, err = ParseBase64(ms)
	if err == nil {
		t.Fatalf("no error parsing invalid base64 %s", ms)
	}

	// Invalid Base64 payload (not 8 bytes and not a decimal number)
	// "bm90X2FfbnVtYmVy" decodes to "not_a_number" (12 bytes)
	_, err = ParseBase64("bm90X2FfbnVtYmVy")
	if err == nil {
		t.Fatalf("expected error parsing non-numeric legacy base64 payload")
	}
}

func TestIntBase64(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	for i := 0; i < 100; i++ {
		oID := node.Generate()
		b64 := oID.IntBase64()

		if len(b64) != 12 {
			t.Fatalf("expected IntBase64 length of 12 for 8 bytes, got %d (%s)", len(b64), b64)
		}
		if b64 != oID.Base64() {
			t.Fatalf("IntBase64 %s != Base64 %s", b64, oID.Base64())
		}

		pID, err := ParseIntBase64(b64)
		if err != nil {
			t.Fatalf("error parsing, %s", err)
		}
		if pID != oID {
			t.Fatalf("pID %v != oID %v", pID, oID)
		}
	}

	knownID := ID(1116823421972381696)
	expectedB64 := "D3/A/C+AAAA="
	if knownID.IntBase64() != expectedB64 {
		t.Fatalf("knownID.IntBase64() %s != %s", knownID.IntBase64(), expectedB64)
	}
	pID, err := ParseIntBase64(expectedB64)
	if err != nil {
		t.Fatalf("error parsing knownID, %s", err)
	}
	if pID != knownID {
		t.Fatalf("pID %v != knownID %v", pID, knownID)
	}

	// Zero ID
	zeroID := ID(0)
	zeroB64 := "AAAAAAAAAAA="
	if zeroID.IntBase64() != zeroB64 {
		t.Fatalf("zeroID.IntBase64() %s != %s", zeroID.IntBase64(), zeroB64)
	}
	pID, err = ParseIntBase64(zeroB64)
	if err != nil {
		t.Fatalf("error parsing zeroID int base64, %s", err)
	}
	if pID != zeroID {
		t.Fatalf("pID %v != zeroID %v", pID, zeroID)
	}

	// Invalid base64 syntax
	_, err = ParseIntBase64("invalid_base64!!!")
	if err == nil {
		t.Fatalf("expected error parsing corrupt base64")
	}

	// Valid base64 but invalid length (not 8 bytes)
	// Legacy string base64 is 19 bytes decoded
	_, err = ParseIntBase64("MTExNjgxOTQ5NDY2MDk5NzEyMA==")
	if err != ErrInvalidBase64 {
		t.Fatalf("expected ErrInvalidBase64 for 19-byte decoded base64, got %v", err)
	}

	// Base64 of 4 bytes: "AAAAAA=="
	_, err = ParseIntBase64("AAAAAA==")
	if err != ErrInvalidBase64 {
		t.Fatalf("expected ErrInvalidBase64 for 4-byte decoded base64, got %v", err)
	}
}

func TestBytes(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.Bytes()

	pID, err := ParseBytes(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := []byte{0x31, 0x31, 0x31, 0x36, 0x38, 0x32, 0x31, 0x36, 0x37, 0x39, 0x35, 0x37, 0x30, 0x34, 0x31, 0x39, 0x37, 0x31, 0x32}
	_, err = ParseBytes(ms)
	if err != nil {
		t.Fatalf("error parsing, %#v", err)
	}

	ms = []byte{0xFF, 0xFF, 0xFF, 0x31, 0x31, 0x31, 0x36, 0x38, 0x32, 0x31, 0x36, 0x37, 0x39, 0x35, 0x37, 0x30, 0x34, 0x31, 0x39, 0x37, 0x31, 0x32}
	_, err = ParseBytes(ms)
	if err == nil {
		t.Fatalf("no error parsing, %#v", err)
	}
}

func TestIntBytes(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.IntBytes()

	pID := ParseIntBytes(i)
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := [8]uint8{0xf, 0x7f, 0xc0, 0xfc, 0x2f, 0x80, 0x0, 0x0}
	mi := int64(1116823421972381696)
	pID = ParseIntBytes(ms)
	if pID.Int64() != mi {
		t.Fatalf("pID %v != mi %v", pID.Int64(), mi)
	}

}

//******************************************************************************
// Marshall Test Methods

func TestMarshalJSON(t *testing.T) {
	id := ID(13587)
	expected := "\"13587\""

	bytes, err := id.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected error during MarshalJSON")
	}

	if string(bytes) != expected {
		t.Fatalf("Got %s, expected %s", string(bytes), expected)
	}
}

func TestMarshalsIntBytes(t *testing.T) {
	id := ID(13587).IntBytes()
	expected := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x35, 0x13}
	if !bytes.Equal(id[:], expected) {
		t.Fatalf("Expected ID to be encoded as %v, got %v", expected, id)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tt := []struct {
		json        string
		expectedID  ID
		expectedErr error
	}{
		{`"13587"`, 13587, nil},
		{`1`, 0, JSONSyntaxError{[]byte(`1`)}},
		{`"invalid`, 0, JSONSyntaxError{[]byte(`"invalid`)}},
	}

	for _, tc := range tt {
		var id ID
		err := id.UnmarshalJSON([]byte(tc.json))
		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Fatalf("Expected to get error '%s' decoding JSON, but got '%s'", tc.expectedErr, err)
		}

		if id != tc.expectedID {
			t.Fatalf("Expected to get ID '%s' decoding JSON, but got '%s'", tc.expectedID, id)
		}
	}
}

// ****************************************************************************
// Benchmark Methods

func BenchmarkParseBase32(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()
	b32i := sf.Base32()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ParseBase32([]byte(b32i))
	}
}
func BenchmarkBase32(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sf.Base32()
	}
}
func BenchmarkParseBase58(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()
	b58 := sf.Base58()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ParseBase58([]byte(b58))
	}
}
func BenchmarkBase58(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sf.Base58()
	}
}

func BenchmarkParseBase64(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()
	b64 := sf.Base64()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = ParseBase64(b64)
	}
}

func BenchmarkBase64(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sf.Base64()
	}
}

func BenchmarkParseIntBase64(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()
	b64 := sf.IntBase64()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = ParseIntBase64(b64)
	}
}

func BenchmarkIntBase64(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sf.IntBase64()
	}
}
func BenchmarkGenerate(b *testing.B) {

	node, _ := NewNode(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Generate()
	}
}

func BenchmarkGenerateMaxSequence(b *testing.B) {

	NodeBits = 1
	StepBits = 21
	node, _ := NewNode(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Generate()
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	// Generate the ID to unmarshal
	node, _ := NewNode(1)
	id := node.Generate()
	bytes, _ := id.MarshalJSON()

	var id2 ID

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = id2.UnmarshalJSON(bytes)
	}
}

func BenchmarkMarshal(b *testing.B) {
	// Generate the ID to marshal
	node, _ := NewNode(1)
	id := node.Generate()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = id.MarshalJSON()
	}
}

func TestParseBase32(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    ID
		wantErr bool
	}{
		{
			name:    "ok",
			arg:     "b8wjm1zroyyyy",
			want:    1427970479175499776,
			wantErr: false,
		},
		{
			name:    "capital case is invalid encoding",
			arg:     "B8WJM1ZROYYYY",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "l is not allowed",
			arg:     "b8wjm1zroyyyl",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "v is not allowed",
			arg:     "b8wjm1zroyyyv",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "2 is not allowed",
			arg:     "b8wjm1zroyyy2",
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBase32([]byte(tt.arg))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBase32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseBase32() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBase58(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    ID
		wantErr bool
	}{
		{
			name:    "ok",
			arg:     "4jgmnx8Js8A",
			want:    1428076403798048768,
			wantErr: false,
		},
		{
			name:    "0 not allowed",
			arg:     "0jgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "I not allowed",
			arg:     "Ijgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "O not allowed",
			arg:     "Ojgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "l not allowed",
			arg:     "ljgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBase58([]byte(tt.arg))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBase58() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseBase58() got = %v, want %v", got, tt.want)
			}
		})
	}
}
