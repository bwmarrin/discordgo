package discordgo

import (
  "testing"
  "os"
)

//////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////// VARS NEEDED FOR TESTING
var (
	dg *Session // Stores global discordgo session

	envToken    = os.Getenv("DG_TOKEN")    // Token to use when authenticating
	envEmail    = os.Getenv("DG_EMAIL")    // Email to use when authenticating
	envPassword = os.Getenv("DG_PASSWORD") // Password to use when authenticating
	envGuild    = os.Getenv("DG_GUILD")    // Guild ID to use for tests
	envChannel  = os.Getenv("DG_CHANNEL")  // Channel ID to use for tests
	//	envUser     = os.Getenv("DG_USER")     // User ID to use for tests
  envEmoji    = os.Getenv("DG_EMOJI")
	envAdmin = os.Getenv("DG_ADMIN") // User ID of admin user to use for tests
)

func init() {
	if envEmail == "" || envPassword == "" || envToken == "" {
		return
	}

	if d, err := New(envEmail, envPassword, envToken); err == nil {
		dg = d
	}
}
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

  g.ID = envGuild

  _, err = s.GuildAdd(g)
  if err != nil {
    t.Errorf("GuildAdd() returned error: %+v", err)
  }

  g2, err := s.Guild(envGuild)
  if err != nil {
    t.Errorf("Guild() returned error: %+v", err)
  }
  if g2 == nil {
    t.Errorf("Guild not found when added")
  }

  _, err = s.GuildRemove(g)
  if err != nil {
    t.Errorf("GuildRemove() returned error: %+v", err)
  }

  g3, err := s.Guild(envGuild)
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
  channel.ID = envChannel

  _, err = s.ChannelAdd(channel)
  if err != nil {
    t.Errorf("ChannelAdd() returned error: %+v", err)
  }

  channel2, err := s.Channel(envChannel)
  if err != nil {
    t.Errorf("Channel() returned error: %+v", err)
  }
  if channel2 == nil {
    t.Errorf("Channel not found when added")
  }

  _, err = s.ChannelRemove(channel)
  if err != nil {
    t.Errorf("ChannelRemove() returned error: %+v", err)
  }

  channel3, err := s.Channel(envChannel)
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

  guild.ID = envGuild
  emoji.ID = "Test Emoji"
  s.GuildAdd(guild)

  _, err = s.EmojiAdd(envGuild, emoji)
  if err != nil {
    t.Errorf("EmojiAdd() returned error: %+v", err)
  }

  emoji2, err := s.Emoji(envGuild, "Test Emoji")
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
  _, err = s.EmojisAdd(envGuild, emojis)
  if err != nil {
    t.Errorf("EmojisAdd() returned error: %+v", err)
  }

  for i := 0; i < 3; i++ {
    emoji2, err := s.Emoji(envGuild, "Test Emoji" + i)
    if err != nil {
      t.Errorf("Emoji() returned error: %+v", err)
    }
    if emoji2 == nil {
      t.Errorf("Emoji not found when added")
    }
  }
}

func TestMemberOperations(t *testing.T) {
    s, err := NewState()
    var m member
    var g Guild

    m.ID = "Test Member"
    guild.ID = envGuild
    s.GuildAdd(guild)

    _,err = s.MemberAdd(envGuild, m)
    if err != nil {
      t.Errorf("GuildAdd() returned error: %+v", err)
    }

    m2, err := s.Member(envGuild, "Test Member")
    if err != nil {
      t.Errorf("Member() returned error: %+v", err)
    }
    if m2==nil {
      t.Errorf("Guild not found when added")
    }

    _, err = m.MemberRemove(envGuild, m)
    if err != nil {
      t.Errorf("MemberRemove() returned error: %+v", err)
    }

    m3, err := m.Member(envGuild, "Test Member")
    if err != nil {
      t.ErrorF("Member() returned error: %+v", err)
    }

    if m3 == nil {
      t.Errorf("Member correctly not found after removal expeced error")
    }
  }

  func TestRoleOperations(t *testing.T) {
    s, err := NewState()
    var r Role
    var g Guild

    r.ID = "Test Role"
    guild.ID = envGuild
    s.GuildAdd(guild)

    _,err = s.RoleAdd(envGuild, r)
    if err != nil {
      t.Errorf("GuildAdd() returned error: %+v", err)
    }

    r2, err := s.Role(envGuild, "Test Role")
    if err != nil {
      t.Errorf("Role() returned error: %+v", err)
    }

    if r2==nil {
      t.Errorf("Guild not found when added")
    }

    _, err = m.RoleRemove(envGuild, r)
    if err != nil {
      t.Errorf("RoleRemove() returned error: %+v", err)
    }

    r3, err := M.Role(envGuild, "Test Role")
    if err != nil {
      t.ErrorF("Role() returned error: %+v", err)
    }

    if r3 == nil {
      t.Errorf("Role correctly not found after removal expeced error")
    }
  }
