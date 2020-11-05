# Basics
This page is dedicated to teaching you the basics like receiving/sending messages, logging the bot in, etc...

## Connecting to discord
```golang
import (
  "os"
  "os/signal"
  "syscall"
  
  "github.com/bwmarrin/discordgo"
)

// this function gets called when the bot is logged in
func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Bot running. Press CTRL+C to exit.")
}

func main() {
  	// create new discordgo instance
  	discord, err := discordgo.New("Bot " + "your discord bot token")
  	if err != nil {
      	fmt.Println("Error: ", err)
      	os.Exit(1)
  	}
  
  	discord.AddHandler(ready)
  
  	// open the websocket connection (actually connect to discord)
  	err = discord.Open()
	if err != nil {
		fmt.Println("Error while connecting to discord: ", err)
   	 	os.Exit(2)
	}
  
  	// wait until the user presses CTRL+C, otherwise the connection would get terminated instantly after connecting
  	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
  
  	// close the connection
  	discord.Close()
}
```

## Receiving a message
The `messageCreate` function gets called when a user sends a message.
```golang
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  	// your code here...
}
```

You have to add a handler for this function.
```golang
discord.AddHandler(ready)
```

It's a good practice to make the bot doesn't respond it's self because this can end in an infinite loop.
```golang
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  	if m.Author.ID == s.State.User.ID {
		return
	}
  	// your code here...
}
```

## Awnsering a message
You can send a message to a channel with the `ChannelMessageSend(channel, message)` function.

e.g. in the `messageCreate` function.
```golang
s.ChannelMessageSend(m.ChannelID, "cool message")
```
