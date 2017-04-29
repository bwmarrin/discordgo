package discordgo

import{
  "testing"
}

//////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////// VARS NEEDED FOR TESTING
var (
  tstate = NewState()

)

//////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////// START OF TESTS

func TestNewState(t *testing.T)  {

  _, err := NewState()
  if err != nil {
    t.Errorf("NewState() returned error: %+v", err)
  }
}

func TestGuildOperations(t *testing.T) {
  s, err := NewState()
  if err != nil {
    t.Errorf("NewState() returned error: %+v", err)
  }

  var g Guild

  g.ID = "Test Guild"

  _, err := s.GuildAdd(g)
  if err != nil {
    t.Errorf("GuildAdd() returned error: %+v", err)
  }

  g2, err := s.Guild("Test Guild")
  if err != nil {
    t.Errorf("Guild() returned error: %+v", err)
  }
  if g2 == nil {
    t.Errorf("Guild not found when added")
  }

  _, err := s.GuildRemove(g)
  if err != nil {
    t.Errorf("GuildRemove() returned error: %+v", err)
  }

  g3, err := s.Guild("Test Guild")
  if err != nil {
    t.Errorf("Guild() returned error: %+v", err)
  }
  if g3 == nil {
    t.Errorf("Guild correctly not found after removal expected error")
  }
}

func TestChannelOperations(t *testing.T)  {
  s, err := NewState()
  var channel Channel
  channel.ID = "Test Channel"

  _, err := s.ChannelAdd(channel)
  if err != nil {
    t.Errorf("ChannelAdd() returned error: %+v", err)
  }

  channel2, err := s.Channel("Test Channel")
  if err != nil {
    t.Errorf("Channel() returned error: %+v", err)
  }
  if channel2 == nil {
    t.Errorf("Channel not found when added")
  }

  _, err := s.ChannelRemove(channel)
  if err != nil {
    t.Errorf("ChannelRemove() returned error: %+v", err)
  }

  channel3, err := s.Channel("Test Channel")
  if err != nil {
    t.Errorf("Channel() returned error: %+v", err)
  }
  if channel3 == nil {
    t.Errorf("Channel correctly not found after removal expected error")
  }
}

func TestEmojiOperations(t *testing.T)  {
  s, err := NewState()
  var emoji Emoji
  var guild Guild

  guild.ID = "Test Guild"
  emoji.ID = "Test Emoji"
  s.GuildAdd(guild)

  _, err := s.EmojiAdd("Test Guild", emoji)
  if err != nil {
    t.Errorf("EmojiAdd() returned error: %+v", err)
  }

  emoji2, err := s.Emoji("Test Guild", "Test Emoji")
  if err != nil {
    t.Errorf("Emoji() returned error: %+v", err)
  }
  if emoji2 == nil {
    t.Errorf("Emoji not found when added")
  }

  var emojis []*Emoji
  for i := 0; i < 3; i++ {
    emoji.ID = "Test Emoji" + i
    emojis[i] = emoji
  }
  _, err := s.EmojisAdd("Test Guild", emojis)
  if err != nil {
    t.Errorf("EmojisAdd() returned error: %+v", err)
  }

  for i := 0; i < 3; i++ {
    emoji2, err := s.Emoji("Test Guild", "Test Emoji" + i)
    if err != nil {
      t.Errorf("Emoji() returned error: %+v", err)
    }
    if emoji2 == nil {
      t.Errorf("Emoji not found when added")
    }
  }
}
