package main

import (
	"GoTorrent/torrent"
	"GoTorrent/torrentfile"
	"log"
	"os"
)

func main() {
	inPath := os.Args[1]
	outPath := os.Args[2]

	tf, err := torrentfile.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}

	t, err := torrent.NewTorrent(tf)
	t.Download(outPath)
}
