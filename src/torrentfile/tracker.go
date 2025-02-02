package torrentfile

import (
	"GoTorrent/peers"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
)

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (t *TorrentFile) buildAnnounceUrl(peerId [20]byte, port int) (string, error) {
	parsedUrl, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	parsedUrl.RawQuery = params.Encode()
	return parsedUrl.String(), nil
}

func (t *TorrentFile) RequestPeers(peerId [20]byte, port int) ([]peers.Peer, error) {
	url, err := t.buildAnnounceUrl(peerId, port)
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	trackerResp := bencodeTrackerResp{}
	err = bencode.Unmarshal(resp.Body, &trackerResp)
	if err != nil {
		return nil, err
	}

	return peers.Unmarshal([]byte(trackerResp.Peers))
}
