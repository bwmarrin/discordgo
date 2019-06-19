package discordgo

import (
	"reflect"
	"testing"
)

func TestSession_State(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *Session
		want  func() *State
	}{
		{
			name: "nil",
			setup: func() *Session {
				return &Session{}
			},
			want: func() *State {
				return nil
			},
		},
		{
			name: "has_guild_id",
			setup: func() *Session {
				s := &Session{state: NewState()}
				_ = s.state.GuildAdd(&Guild{ID: "42"})
				_ = s.state.RoleAdd("guild", &Role{ID: "84"})
				return s
			},
			want: func() *State {
				s := NewState()
				_ = s.GuildAdd(&Guild{ID: "42"})
				_ = s.RoleAdd("guild", &Role{ID: "84"})
				return s
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.setup()
			want := tt.want()

			if inState := in.State(); !reflect.DeepEqual(inState, want) {
				t.Fatalf("%#v\n\nnot equal to:\n\n%#v", inState, want)
			}
		})
	}
}
