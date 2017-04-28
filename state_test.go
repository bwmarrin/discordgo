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

func TestGuildAdd(t *testing.T) {
  s, err := NewState()
  if err != nil {
    t.Errorf("NewState() returned error: %+v", err)
  }

  var g Guild

  _, err := s.GuildAdd(g)
  if err != nil {
    t.Errorf("GuildAdd() returned error: %+v", err)
  }
}
