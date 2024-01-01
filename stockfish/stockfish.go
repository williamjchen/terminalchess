package stockfish

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"bufio"
)

func Move(moveHistory string, depth int) string {	
	cmd := exec.Command("./stockfish/stockfish-ubuntu-x86-64")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		slog.Error("unable to create stdin pipe for stockfish", err)
		return ""
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("unable to create stdoutpipe for stockfish", err)
		return ""
	}

	scanner := bufio.NewScanner(stdout)

	err = cmd.Start()
	if err != nil {
		slog.Error("Unable to start stockfish", err)
		return ""
	}

	fmt.Fprintln(stdin, "uci")
	fmt.Fprintln(stdin, "position startpos moves ", moveHistory)
	fmt.Fprintf(stdin, "go depth %d\n", depth)

	var s string
	for scanner.Scan() {
		t := scanner.Text()
		s = getBestMove(t)
		if s != "" {
			break
		}
	}

	cmd.Process.Kill()
	cmd.Wait()
	return s
}

func getBestMove(s string) string {
	words := strings.Split(s, " ")
	if len(words) < 2 || words[0] != "bestmove" {
		return ""
	}
	return words[1]
}

