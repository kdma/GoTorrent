package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"io"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	PieceLength int
	InfoHash    [20]byte
	PieceHashes [][20]byte
	Length      int
	Name        string
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func Open(path string) (*TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bto := bencodeTorrent{}
	var r io.Reader
	r = file
	err = bencode.Unmarshal(r, &bto)
	if err != nil {
		return nil, err
	}

	infoHash, err := bto.Info.hash()
	if err != nil {
		return nil, err
	}

	hashes, err := splitHashes(&bto.Info)
	if err != nil {
		return nil, err
	}

	return &TorrentFile{
		Announce:    bto.Announce,
		PieceLength: bto.Info.PieceLength,
		PieceHashes: hashes,
		InfoHash:    infoHash,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}, nil
}

func splitHashes(bto *bencodeInfo) ([][20]byte, error) {

	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(bto.Pieces)

	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}
