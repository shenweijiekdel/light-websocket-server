package stanza

type (
	Stanza interface {
		Type() byte
	}
)

type Message []byte

func (q Message) Type() byte {
	return TypeMessage
}
