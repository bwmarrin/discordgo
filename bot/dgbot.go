package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	Discord "github.com/bwmarrin/discordgo"
)

// Global Variables
var (
	Session  Discord.Session
	Username string
	Password string
)

func main() {
	fmt.Printf("\nDiscordgo Bot Starting.\n\n")

	// Register all the Event Handlers
	RegisterHandlers()

	// read in the config file.
	ParseFile()

	// seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// main program loop to keep dgbot running
	// will add stuff here to track goroutines
	// and monitor for CTRL-C or other things.
	for {
		time.Sleep(1000 * time.Millisecond)
	}

	fmt.Println("\nDiscordgo Bot shutting down.\n")
}

// ParseFile will read a file .dgbotrc and run all included
// commands
func ParseFile() {
	file, err := os.Open(".dgbotrc")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		fmt.Println(admin(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
