package discordgo

import{
  "testing"
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

func TestMemberOperations(t *testing.T) {
    s, err := NewState()
    var m member
    var g Guild

    m.ID = "Test Member"
    guild.ID = "Test Guild"
    s.GuildAdd(guild)

    _,err := s.MemberAdd("Test Guild", m)
    if err != nil {
      t.Errorf("GuildAdd() returned error: %+v", err)
    }

    m2, err := s.Member("Test Guild", "Test Member")
    if err != nil {
      t.Errorf("Member() returned error: %+v", err)
    }
    if m2==nil {
      t.Errorf("Guild not found when added")
    }

    _, err := m.MemberRemove("Test Guild", m)
    if err != nil {
      t.Errorf("MemberRemove() returned error: %+v", err)
    }

    m3, err := m.Member("Test Guild", "Test Member")
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
    guild.ID = "Test Guild"
    s.GuildAdd(guild)

    _,err := s.RoleAdd("Test Guild", r)
    if err != nil {
      t.Errorf("GuildAdd() returned error: %+v", err)
    }

    r2, err := s.Role("Test Guild", "Test Role")
    if err != nil {
      t.Errorf("Role() returned error: %+v", err)
    }
    
    if r2==nil {
      t.Errorf("Guild not found when added")
    }

    _, err := m.RoleRemove("Test Guild", r)
    if err != nil {
      t.Errorf("RoleRemove() returned error: %+v", err)
    }

    r3, err := M.Role("Test Guild", "Test Role")
    if err != nil {
      t.ErrorF("Role() returned error: %+v", err)
    }

    if r3 == nil {
      t.Errorf("Role correctly not found after removal expeced error")
    }
  }
}
