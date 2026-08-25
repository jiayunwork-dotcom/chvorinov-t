package main

import (
	"fmt"
	"net/http"
	"os"

	"chvorinov-t/internal/cli"
	"chvorinov-t/internal/server"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && (args[0] == "-http" || args[0] == "--http") {
		addr := ":8080"
		if len(args) > 1 {
			addr = args[1]
		}
		fmt.Printf("chvorinov-t HTTP service on %s (POST /api/freeze)\n", addr)
		if err := server.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
			fmt.Fprintln(os.Stderr, "chvorinov-t:", err)
			os.Exit(1)
		}
		return
	}
	os.Exit(cli.Run(args, os.Stdout, os.Stderr))
}
