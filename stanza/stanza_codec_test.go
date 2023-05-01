package stanza

import "testing"

func TestStanzaCodec(t *testing.T) {
	var inputs = []Stanza{
		Ping{}, Pong{}, Message("哈哈"),
	}

	for _, input := range inputs {
		t.Logf("Test input %v\n", inputs)
		encode, err := Encode(input)
		if err != nil {
			t.Fatalf("encode error: %v\n", err)
			return
		}

		decode, err := Decode(encode)
		if err != nil {
			t.Fatalf("decode error: %v\n", err)
			return
		}

		if !stanzaEquals(input, decode) {
			t.Fatal("stanzaEquals: not match\n", err)
			return
		}
	}
}

func stanzaEquals(s1 Stanza, s2 Stanza) bool {
	if _, ok := s1.(Ping); ok {
		return s1 == s2
	}

	if _, ok := s1.(Pong); ok {
		return s1 == s2
	}

	if message, ok := s1.(Message); ok {
		return byteSliceEquals(message, s2.(Message))
	}

	return false
}

func byteSliceEquals(b1 []byte, b2 []byte) bool {
	if len(b1) != len(b2) {
		return false
	}

	for i := range b1 {
		if b1[i] != b2[i] {
			return false
		}
	}

	return true
}
