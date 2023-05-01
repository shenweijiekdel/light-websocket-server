package stanza

const (
	TypeKickoff byte = 0x03
	TypeMessage byte = 0x04
	TypePing    byte = 0x05
	TypePong    byte = 0x06
)

type (
	Stanza interface {
		Type() byte
	}
)

type Message []byte

func (q Message) Type() byte {
	return TypeMessage
}

type Ping struct {
}

type Pong struct {
}

type Kickoff struct {
}

func (p Ping) Type() byte {
	return TypePing
}

func (p Pong) Type() byte {
	return TypePong
}

func (p Kickoff) Type() byte {
	return TypeKickoff
}
