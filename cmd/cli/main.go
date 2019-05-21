package main

import "github.com/VideoCoin/cli/internal/cmd"

var (
	Version = "0.0.1"
	Build   = "3"
)

func main() {
	cmd.Execute(Version, Build)
}
