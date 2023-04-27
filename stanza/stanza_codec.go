package stanza

import (
	"bytes"
	"errors"
)

const (
	TypeMessage byte = 0x04
	TypePing    byte = 0x05
	TypePong    byte = 0x06
)

func Encode(stanza Stanza) ([]byte, error) {
	switch stanza.Type() {
	case TypePing:
		return encodePing()

	case TypePong:
		return encodePong()

	case TypeMessage:
		return encodeMessage(stanza.(Message))

	default:
		return nil, errors.New("invalid stanza type")
	}
}

func encodeMessage(event Message) ([]byte, error) {
	var b = []byte{TypeMessage}

	buf := bytes.Buffer{}
	buf.WriteByte(TypeMessage)
	buf.Write(event)

	return b, nil
}

func encodePing() ([]byte, error) {
	return []byte{TypePing}, nil
}

func encodePong() ([]byte, error) {
	return []byte{TypePong}, nil
}

func Decode(b []byte) (Stanza, error) {
	buf := bytes.NewBuffer(b)
	stanzaType, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	switch stanzaType {
	case TypePing:
		return Ping{}, nil

	case TypePong:
		return Pong{}, nil

	case TypeMessage:
		return decodeEvent(buf)

	default:
		return nil, errors.New("invalid stanza type")
	}
}

func decodeEvent(buf *bytes.Buffer) (*Message, error) {
	var message Message
	_, err := buf.Read(message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
