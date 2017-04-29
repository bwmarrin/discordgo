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

  g.guildID = "Test Guild"

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
