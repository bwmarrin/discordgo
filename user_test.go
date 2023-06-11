package discordgo

import "testing"

func TestUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		u    *User
		want string
	}{
		{
			name: "Username with a discriminator",
			u: &User{
				Username:      "bob",
				Discriminator: "8192",
			},
			want: "bob#8192",
		},
		{
			name: "Given a user with a username and no discriminator, When String() is called, Then it should return the username",
			u: &User{
				Username: "aldiwildan",
			},
			want: "aldiwildan",
		},
		{
			name: "Given a user with a username and a 0 discriminator, When String() is called, Then it should return the username",
			u: &User{
				Username:      "aldiwildan",
				Discriminator: "0",
			},
			want: "aldiwildan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.want {
				t.Errorf("User.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
