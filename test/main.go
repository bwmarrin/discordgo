package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
	"path"
)

func main() {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	path.Join(dir, "tmp2.json")
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		panic(b)
	}

	g := &discordgo.GuildCreate{}
	if err := json.Unmarshal(b, &g); err != nil {
		panic(err)
	}
}
