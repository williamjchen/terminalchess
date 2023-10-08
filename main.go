package main

import (
	"github.com/williamjchen/terminalchess/server"
	"fmt"
	
)

func main() {
	fmt.Println("Hello!")

	server.NewServer(".ssh/term_info_ed25519", "localhost", 2324)
}

