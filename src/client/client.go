package client

import (
	"GoTorrent/bitfield"
	"GoTorrent/handshake"
	"GoTorrent/message"
	"GoTorrent/peers"
	"bytes"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield bitfield.Bitfield
}

func NewClient(peer *peers.Peer, h *handshake.Handshake) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 30*time.Second)
	if err != nil {
		return nil, err
	}

	_, err = completeHandshake(conn, h)
	if err != nil {
		return nil, err
	}

	bf, err := recvBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
	}, nil
}

// Read reads and consumes a message from the connection
func (c *Client) Read() (*message.Message, error) {
	msg, err := message.Read(c.Conn)
	return msg, err
}

func completeHandshake(conn net.Conn, h *handshake.Handshake) (*handshake.Handshake, error) {
	conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	write, err := conn.Write(h.Serialize())
	fmt.Printf("Handshake sent\n %d", write)
	if err != nil {
		return nil, err
	}

	res, err := handshake.Deserialize(conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(res.InfoHash[:], h.InfoHash[:]) {
		return nil, fmt.Errorf("expected infohash %x but got %x", res.InfoHash, res.InfoHash)
	}

	return res, nil
}

func recvBitfield(conn net.Conn) (bitfield.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}

	if msg == nil {
		err := fmt.Errorf("expected bitfield but got %v", msg)
		return nil, err
	}

	if msg.ID != message.Bitfield {
		err := fmt.Errorf("expected bitfield but got ID %d", msg.ID)
		return nil, err
	}

	return msg.Payload, nil
}

// SendInterested sends an Interested message to the peer
func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.Interested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendNotInterested sends a NotInterested message to the peer
func (c *Client) SendNotInterested() error {
	msg := message.Message{ID: message.NotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke sends an Unchoke message to the peer
func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.Unchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendRequest sends a Request message to the peer
func (c *Client) SendRequest(index, begin, length int) error {
	req := message.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}

// SendHave sends a Have message to the peer
func (c *Client) SendHave(index int) error {
	msg := message.FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}
