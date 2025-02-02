package torrentfile

import (
	"crypto/rand"
	"testing"

	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	torrent, err := Open("testdata/archlinux-2019.12.01-x86_64.iso.torrent")
	require.Nil(t, err)

	expected := &TorrentFile{
		Announce:    "http://tracker.archlinux.org:6969/announce",
		InfoHash:    [20]byte{222, 232, 106, 127, 166, 242, 134, 169, 215, 76, 54, 32, 20, 97, 106, 15, 245, 228, 132, 61},
		PieceLength: 524288,
		Length:      670040064,
		Name:        "archlinux-2019.12.01-x86_64.iso",
	}

	assert.Equal(t, expected.Announce, torrent.Announce)
	assert.Equal(t, expected.InfoHash, torrent.InfoHash)
	assert.Equal(t, expected.PieceLength, torrent.PieceLength)
	assert.Equal(t, expected.Length, torrent.Length)
	assert.Equal(t, expected.Name, torrent.Name)
}

func TestPeersResponse(t *testing.T) {
	tFile := &TorrentFile{
		Announce:    "http://www.torrentsnipe.info:2701/announce",
		InfoHash:    [20]byte{222, 232, 106, 127, 166, 242, 134, 169, 215, 76, 54, 32, 20, 97, 106, 15, 245, 228, 132, 61},
		PieceLength: 524288,
		Length:      670040064,
		Name:        "archlinux-2019.12.01-x86_64.iso",
	}

	peerId := [20]byte{}
	rand.Read(peerId[:])

	peers, err := tFile.RequestPeers(peerId, 2701)
	assert.Equal(t, err, nil)
	assert.Greater(t, len(peers), 0)
	assert.NotEqual(t, peers[0].IP, nil)
}
