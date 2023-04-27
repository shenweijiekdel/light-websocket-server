package stanza

type Ping struct {
}

type Pong struct {
}

func (p Ping) Type() byte {
	return TypePing
}

func (p Pong) Type() byte {
	return TypePong
}
