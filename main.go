package main

import (
	"log"

	"ai-shell-windows/config"
	"ai-shell-windows/shell"
)

func main() {
	cfg := config.Default()
	if err := shell.StartShell(cfg); err != nil {
		log.Fatalf("shell exited with error: %v", err)
	}
}
