package main

import (
	"math/rand"
	"time"

	"github.com/buglloc/rip/v2/commands"
)

func main() {
	rand.Seed(time.Now().Unix())
	commands.Execute()
}
