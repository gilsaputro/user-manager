package main

import (
	"gilsaputro/user-manager/cmd/user-manager/server"
	"os"
)

func main() {
	os.Exit(server.Run())
}
