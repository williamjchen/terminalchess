package main

import (
	"github.com/williamjchen/terminalchess/server"	
)

func main() {
	s, err := server.NewServer(".ssh/term_info_ed25519", "0.0.0.0", 2324)
	if err != nil {
		return
	}
	s.Start()
}

