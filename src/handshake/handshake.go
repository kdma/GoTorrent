package handshake

import (
	"fmt"
	"io"
)

// A Handshake is a special message that a peer uses to identify itself
type Handshake struct {
	Prefix   byte
	Pstr     string
	Reserved [8]byte
	InfoHash [20]byte
	PeerID   [20]byte
}

func NewHandshake(infoHash [20]byte, peerId [20]byte) *Handshake {
	return &Handshake{
		Prefix:   19,
		Pstr:     "BitTorrent protocol",
		Reserved: [8]byte{},
		InfoHash: infoHash,
		PeerID:   peerId,
	}
}

func (h *Handshake) Serialize() []byte {
	buf := make([]byte, 1+len(h.Pstr)+len(h.Reserved)+len(h.InfoHash)+len(h.PeerID))
	buf[0] = h.Prefix
	curr := 1
	curr += copy(buf[curr:], h.Pstr)
	curr += copy(buf[curr:], h.Reserved[:]) // 8 reserved bytes
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

func Deserialize(r io.Reader) (*Handshake, error) {
	dummy := NewHandshake([20]byte{}, [20]byte{})
	buf := make([]byte, 1+len(dummy.Pstr)+len(dummy.Reserved)+len(dummy.InfoHash)+len(dummy.PeerID))
	readBytes, error := r.Read(buf)
	if error != nil {
		return nil, error
	}

	if readBytes < 1 {
		return nil, fmt.Errorf("received 0 bytes")
	}

	reader := reader()

	return &Handshake{
		Prefix:   reader(buf, 1)[0],
		Pstr:     string(reader(buf, len(dummy.Pstr))),
		Reserved: [8]byte(reader(buf, len(dummy.Reserved))),
		InfoHash: [20]byte(reader(buf, len(dummy.InfoHash))),
		PeerID:   [20]byte(reader(buf, len(dummy.PeerID))),
	}, nil
}

func reader() func(buffer []byte, toRead int) (b []byte) {
	var count = 0
	return func(buffer []byte, toRead int) (b []byte) {
		parsed := make([]byte, toRead)
		copied := copy(parsed[:], buffer[count:count+toRead])
		count += copied
		return parsed[:]
	}
}
