package stanza

import (
	"bytes"
	"errors"
)

func Encode(stanza Stanza) ([]byte, error) {
	switch stanza.Type() {
	case TypeKickoff:
		return encodeKickoff()

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
	buf := bytes.Buffer{}
	buf.WriteByte(TypeMessage)
	buf.Write(event)

	return buf.Bytes(), nil
}

func encodePing() ([]byte, error) {
	return []byte{TypePing}, nil
}

func encodeKickoff() ([]byte, error) {
	return []byte{TypeKickoff}, nil
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

	case TypeKickoff:
		return Kickoff{}, nil

	case TypeMessage:
		return decodeMessage(buf)

	default:
		return nil, errors.New("invalid stanza type")
	}
}

func decodeMessage(buf *bytes.Buffer) (Message, error) {
	message := make([]byte, buf.Len())
	_, err := buf.Read(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}
