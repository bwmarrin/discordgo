package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var oauthConfig = oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://discord.com/oauth2/authorize",
		TokenURL: "https://discord.com/api/oauth2/token",
	},
	Scopes: []string{"identify", "role_connections.write"},
}

var (
	appID        = flag.String("app", "", "Application ID")
	token        = flag.String("token", "", "Application token")
	clientSecret = flag.String("secret", "", "OAuth2 secret")
	redirectURL  = flag.String("redirect", "", "OAuth2 Redirect URL")
)

func init() {
	flag.Parse()
	godotenv.Load()
	oauthConfig.ClientID = *appID
	oauthConfig.ClientSecret = *clientSecret
	oauthConfig.RedirectURL, _ = url.JoinPath(*redirectURL, "/linked-roles-callback")
}

func main() {
	s, _ := discordgo.New("Bot " + *token)

	_, err := s.ApplicationRoleConnectionMetadataUpdate(*appID, []*discordgo.ApplicationRoleConnectionMetadata{
		{
			Type:                     discordgo.ApplicationRoleConnectionMetadataIntegerGreaterThanOrEqual,
			Key:                      "loc",
			Name:                     "Lines of Code",
			NameLocalizations:        map[discordgo.Locale]string{},
			Description:              "Total lines of code written",
			DescriptionLocalizations: map[discordgo.Locale]string{},
		},
		{
			Type:                     discordgo.ApplicationRoleConnectionMetadataBooleanEqual,
			Key:                      "gopher",
			Name:                     "Gopher",
			NameLocalizations:        map[discordgo.Locale]string{},
			Description:              "Writes in Go",
			DescriptionLocalizations: map[discordgo.Locale]string{},
		},
		{
			Type:                     discordgo.ApplicationRoleConnectionMetadataDatetimeGreaterThanOrEqual,
			Key:                      "first_line",
			Name:                     "First line written",
			NameLocalizations:        map[discordgo.Locale]string{},
			Description:              "Days since the first line of code",
			DescriptionLocalizations: map[discordgo.Locale]string{},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Updated application metadata")
	http.HandleFunc("/linked-roles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		// Redirect the user to Discord OAuth2 page.
		http.Redirect(w, r, oauthConfig.AuthCodeURL("random-state"), http.StatusMovedPermanently)
	})
	http.HandleFunc("/linked-roles-callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		// A safeguard against CSRF attacks.
		// Usually tied to requesting user or random.
		// NOTE: Hardcoded for the sake of the example.
		if q["state"][0] != "random-state" {
			return
		}

		// Fetch the tokens with code we've received.
		tokens, err := oauthConfig.Exchange(r.Context(), q["code"][0])
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Construct a temporary session with user's OAuth2 access_token.
		ts, _ := discordgo.New("Bearer " + tokens.AccessToken)

		// Retrive the user data.
		u, err := ts.User("@me")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Fetch external metadata...
		// NOTE: Hardcoded for the sake of the example.
		metadata := map[string]string{
			"gopher":     "1", // 1 for true, 0 for false
			"loc":        "10000",
			"first_line": "1970-01-01", // YYYY-MM-DD
		}

		// And submit it back to discord.
		_, err = ts.UserApplicationRoleConnectionUpdate(*appID, &discordgo.ApplicationRoleConnection{
			PlatformName:     "Discord Gophers",
			PlatformUsername: u.Username,
			Metadata:         metadata,
		})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// Retrieve it to check if everything is ok.
		info, err := ts.UserApplicationRoleConnection(*appID)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		jsonMetadata, _ := json.Marshal(info.Metadata)
		// And show it to the user.
		w.Write([]byte(fmt.Sprintf("Your updated metadata is: %s", jsonMetadata)))
	})
	http.ListenAndServe(":8010", nil)
}
