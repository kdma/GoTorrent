package main

import (
	"GoTorrent/torrent"
	"GoTorrent/torrentfile"
	"log"
)

func main() {
	inPath := "C:\\Users\\franc\\Desktop\\debian-12.9.0-amd64-netinst.iso.torrent"
	outPath := "debian.iso"

	tf, err := torrentfile.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}

	t, err := torrent.NewTorrent(tf)
	t.Download(outPath)
}
