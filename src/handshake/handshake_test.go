package handshake

import (
	"GoTorrent/utils"
	"bytes"
	"io"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestDeserialize(t *testing.T) {
	infoHash := utils.GenerateRandomBytes()
	peerId := utils.GenerateRandomBytes()

	handshake := NewHandshake(infoHash, peerId)
	serialized := handshake.Serialize()

	var r io.Reader = bytes.NewReader(serialized)
	deserialized, err := Deserialize(r)
	assert.NoError(t, err, "Error deserializing handshake")

	assert.Equal(t, handshake.Prefix, deserialized.Prefix, "Prefix does not match")
	assert.Equal(t, handshake.Pstr, deserialized.Pstr, "Pstr does not match")
	assert.Equal(t, handshake.Reserved, deserialized.Reserved, "Reserved bytes do not match")
	assert.Equal(t, handshake.InfoHash, deserialized.InfoHash, "InfoHash does not match")
	assert.Equal(t, handshake.PeerID, deserialized.PeerID, "PeerID does not match")
}
