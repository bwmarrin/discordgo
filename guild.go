package discordgo

/*
TODO: How to name these? If we make a variable to store
channels from READY packet, etc.  We can't have a Channel
func?  And which is better.  Channels func gets live up
to date data on each call.. so, there is some benefit there.

Maybe it should have both, but make the Channels check and
pull new data based on a cache time?

func (s *Server) Channels() (c []Channel, err error) {
	c, err = Channels(s.Session, s.Id)
	return
}
*/
/*

func (s *Server) Members() (m []Users, err error) {
}
*/
