package discordgo

import "testing"

func TestUser_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		u    *User
		want string
	}{
		{
			name: "User with a discriminator",
			u: &User{
				Username:      "bob",
				Discriminator: "8192",
			},
			want: "bob#8192",
		},
		{
			name: "User without a discriminator",
			u: &User{
				Username: "aldiwildan",
			},
			want: "aldiwildan",
		},
		{
			name: "User with discriminator set to 0",
			u: &User{
				Username:      "aldiwildan",
				Discriminator: "0",
			},
			want: "aldiwildan",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.u.String(); got != tc.want {
				t.Errorf("User.String() = %v, want %v", got, tc.want)
			}
		})
	}
}
