package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// structure for holding session and settings for http server and oauth2 requests
type handler struct {
	// Discord session for preforming oauth2 requests
	session *discordgo.Session
	// Discord application oauth2 client id
	id string
	/* Applications redirect uri registered with discord on the application/oauth2 
			dev portal page
	*/
	uri *url.URL
	// Whether log response to stdout
	log bool
	// What scope ( permissions ) to request from discord in oauth2 requests
	scope string
	// Cache of states used in oauth2 redirect
	states map[string]int64
}

// ServeHTTP handler for go's http server
func( oauth2 handler ) ServeHTTP( response http.ResponseWriter, request *http.Request ) {
	if request.URL.Path != oauth2.uri.Path {
		response.Header().Set( "Content-Type", "text/plain" )
		response.WriteHeader( http.StatusNotFound )
		response.Write( []byte( 
			"Only the path in -redirect <uri> ( " +
			oauth2.uri.Path +
			" ) is valid for this test server" ) )
		return
	}
	
	values := request.URL.Query( )

	switch {
	case values.Has( "state" ) || values.Has( "code" ):
		oauth2.doAccessToken( response, request )
	
	case values.Has( "token" ):
		oauth2.doAuthorizationInfo( response, request )
	
	case values.Has( "cred" ):
		oauth2.doClientCredentials( response, request )

	default:
		// create a random 32 byte state and base64 encode it
		var buffer = make( []byte, 32 )
		rand.Read( buffer )
		state := base64.URLEncoding.EncodeToString( buffer )

		oauth2.states[ state ] = time.Now( ).UTC( ).Unix( )

		location := 
			"https://discord.com/oauth2/authorize?" +
			"client_id=" + oauth2.id +
			"&response_type=code" +
			"&redirect_uri=" + oauth2.uri.String( ) +
			"&scope=" + oauth2.scope +
			"&state=" + state

		response.Header( ).Set( "Location", location )
		response.WriteHeader( http.StatusTemporaryRedirect )
	}
}

func ( oauth2 handler ) doAccessToken ( response http.ResponseWriter, request *http.Request ) {
	code := request.URL.Query( ).Get( "code" )
	state := request.URL.Query( ).Get( "state" )

	created, has := oauth2.states[ state ]
	if !has {
		response.Header( ).Set( "Content-Type", "text/plain" )
		response.WriteHeader( http.StatusForbidden )
		response.Write( []byte( "OAuth2 state not found" ) )
		return
	}

	delete ( oauth2.states, state )
	if created < time.Now( ).UTC( ).Unix( ) - 450 {
		response.Header( ).Set( "Content-Type", "text/plain" )
		response.WriteHeader( http.StatusForbidden )
		response.Write( []byte( "OAuth2 state is over 7.5 minutes old" ) )
		return
	}

	var data []byte
	accessToken, e := oauth2.session.AccessToken( code, oauth2.uri.String( ) )
	if e == nil {
		data, e = json.MarshalIndent( accessToken, "  ", "  " )
	}

	if e != nil {
		response.Header( ).Set( "Content-Type", "text/plain" )
		response.WriteHeader( http.StatusInternalServerError )
		response.Write( []byte( "Error while getting/processing access token request\n" + e.Error( ) ) )
		return
	}

	if oauth2.log {
		log.Printf( "Access Token Response:\n%s\n", data )
	}

	body := bytes.NewBuffer( make( []byte, 0, 256 ) )
	body.WriteString( 
`<!DOCTYPE html>
	<html>
	<head><title>Access Token Response</title></head>
	<body>
		<h3>Access Token</h3>
		<pre>`)
	body.Write( data )
	body.WriteString(
`		</pre>
		<p>
			For more information see the 
			<a href=https://discord.com/developers/docs/topics/oauth2#authorization-code-grant>
				Discord Developer's Documentation: Authorization Code Grant</a>
		</p>
	</body>
</html>`)

	response.Header( ).Set( "Content-Type", "text/html" )
	response.WriteHeader( http.StatusOK )
	response.Write( body.Bytes( ) )

}

func ( oauth2 handler ) doAuthorizationInfo ( response http.ResponseWriter, request *http.Request ) {
	token := request.URL.Query( ).Get( "token" )
	if len( token ) < 1 {
		response.Header( ).Set( "Content-Type", "text/plain" )
		response.WriteHeader( http.StatusBadRequest )
		response.Write( []byte( "Token required in the format token='token'\nDo not add the 'Bearer ' to the start" ) )
		return
	}
	
	discord, e := discordgo.New( discordgo.BearerToken( token ) )

	if e != nil {
		response.Header( ).Set( "Content-Type", "text/plain" )
		response.WriteHeader( http.StatusInternalServerError )
		response.Write( []byte( e.Error( ) ) )
		return
	}

	var data []byte
	authinfo, e := discord.AuthorizationInfo( )
	if e == nil {
		data, e = json.MarshalIndent( authinfo, "  ", "  " )
	}
	
	if e != nil {
		response.Header( ).Set( "Content-Type", "text/plain" )
		response.WriteHeader( http.StatusInternalServerError )
		response.Write( 
			[]byte( 
				"Unable to get or process token authorization info due to error:\n" + e.Error( ) ) )
		return
	}

	if oauth2.log {
		log.Printf( "Authorization Info:\n%s\n", data )
	}

	body := bytes.NewBuffer( make( []byte, 0, 256 ) )
	body.WriteString( 
`<!DOCTYPE html>
	<html>
	<head><title>Token Authorization Info</title></head>
	<body>
		<h3>Authorization Info</h3>
		<pre>`)
	body.Write( data )
	body.WriteString(
`		</pre>
		<p>
			For more information see the 
			<a href=https://discord.com/developers/docs/topics/oauth2#get-current-authorization-information>
				Discord Developer's Documentation: Current Authorization Information</a>
		</p>
	</body>
</html>`)

	response.Header( ).Set( "Content-Type", "text/html" )
	response.WriteHeader( http.StatusOK )
	response.Write( body.Bytes( ) )
}

func ( oauth2 handler ) doClientCredentials ( response http.ResponseWriter, _ *http.Request ) {
	credentials, e := oauth2.session.ClientCredentials( oauth2.scope )
	if e != nil {
		response.Header( ).Set( "Content-Type", "text/plain" )
		response.WriteHeader( http.StatusBadGateway )
		response.Write( 
			[]byte( 
				"Unable to get client credentials due to error:\n" + e.Error( ) ) )
		return
	}

	var data []byte

	if oauth2.log {
		if data, e = json.MarshalIndent( credentials, "  ", "  " ); e != nil {
			log.Println( e )
		} else {
			log.Printf( "Client Credentials:\n%s\n", data )
		}
	}

	credentials.AccessToken = "ABCEFG123REDACTEDTOKEN"
	credentials.RefreshToken = "ABCDEFG123REDACTEDTOKEN"
	
	if data, e = json.MarshalIndent( credentials, "", "  " ); e != nil {
		response.Header( ).Set( "Content-Type", "text/plain" )
		response.WriteHeader( http.StatusInternalServerError )
		response.Write( []byte( "Error while serializing client credentials\n" + e.Error( ) ) )
		return
	}
	
	body := bytes.NewBuffer( make( []byte, 0, 256 ) )
	body.WriteString( 
`<!DOCTYPE html>
	<html>
	<head><title>Client Credentials Grant</title></head>
	<body>
		<h3>Client Credentials Grant</h3>
		<pre>`)
	body.Write( data )
	body.WriteString(
`		</pre>
		<p>
			Client Credentials Grants are a Discord oauth2 bearer token for the application 
			developer's discord account and are ment for testing purposes.
		</p>
		<p>
			For more information see the 
			<a href=https://discord.com/developers/docs/topics/oauth2#client-credentials-grant>
				Discord Developer's Documentation: Client Credential Grants</a>
		</p>
	</body>
</html>`)

	response.Header( ).Set( "Content-Type", "text/html" )
	response.WriteHeader( http.StatusOK )
	response.Write( body.Bytes( ) )
}

func main ( ) {
	var secret, redirect, cert, key  string
  var port int
	var server handler

	flag.StringVar( 
		&server.id, 
		"id", 
		"", 
		"Client id taken from the Settings>OAuth2 section of your application on the Discord Developer's Portal" )

	flag.StringVar( 
		&secret, 
		"secret", 
		"", 
		"Client secret taken from the Settings>Oauth2 section of your application on the Discord Developer's Portal" )

	flag.StringVar( 
		&redirect, 
		"redirect", 
		"", 
		"Redirect URI set in the Settings>OAuth2 section of your application on the Discord Developer's Portal" )

	flag.StringVar( 
		&server.scope, 
		"scope", 
		"identify", 
		"Scope to use for all OAuth2 requests made with the demo server" )

	flag.IntVar( 
		&port, 
		"port", 
		0, 
		"Port number for the demo server to bind to, must be a number between 1 and 65535" )

	flag.StringVar(
		&cert,
		"cert",
		"",
		"Full path to the ssl certificate file if you wish to run in https mode" )

	flag.StringVar(
		&key,
		"key",
		"",
		"Full path to the ssl key file if you wish to run in https mode" )

	flag.BoolVar( 
		&server.log, 
		"log",
		false,
		"Enable logging responses to stdout" )

	flag.Parse( )

	if server.id == "" { 
		log.Fatalln( 
			"-client <id> is required and must be a discord application client id" ) 
	}

	if secret == "" {
		log.Fatalln( 
			"-secret <secret> is required and must be a discord client secret" )
	}

	if redirect == "" {
		log.Fatalln( 
			"-redirect <uri> is required and must be a redirect uri registered to your application with discord" )
	}

	var e error
	server.uri, e = url.Parse( redirect )
	if e != nil {
		log.Fatalf( "-redirect <uri> must be a valid url\n%s\n", e.Error( ) )
	}

	if !( server.uri.Scheme == "http" || server.uri.Scheme == "https" ) || 
	len( server.uri.Hostname( ) ) < 5 {
		log.Fatalln( "-redirect <uri> must be the full url for discord to send user back to and must be the same as registered with discord" )
	}

	if len( server.scope ) < 1 {
		log.Fatalln( "-scope is required and must consist of valid discord oauth2 scopes" )
	}
	for _, scope := range strings.Fields( server.scope ) {
		switch scope {
		case 
			"activites.read", 
			"activities.write",
			"applications.builds.read",
			"applications.builds.upload",
			"applications.commands",
			"applications.commands.update",
			"applications.commands.permissions.update",
			"applications.entitlements",
			"applications.store.update",
			"bot",
			"connections",
			"dm_channels.read",
			"email",
			"gdm.join",
			"guilds",
			"guilds.join",
			"guilds.members.read",
			"identify",
			"messages.read",
			"relationships.read",
			"role_connections.write",
			"rpc",
			"rpc.activities.write",
			"rpc.notifications.read",
			"rpc.voice.read",
			"rpc.voice.write",
			"voice",
			"webhook.incoming":
				continue
		}
		log.Fatalf( 
			"Invalid scope %q see https://discord.com/developers/docs/topics/oauth2#shared-resources-oauth2-scopes for permitted scopes\n", e.Error( ) )
	}

	var bind string
	if port == 0 {
		bind = ":"
	} else if	 port > 0 && port < 65536 {
		bind = ":" + strconv.Itoa( port )
	} else {
		log.Fatal( "-port must be a number between 1 and 65535" )
	}

	ssl := false
	if certSet, keySet := len( cert ) > 0, len( key ) > 0; certSet && keySet {
		ssl = true
	} else if certSet || keySet {
		log.Fatalln( "Both -cert <filepath> and -key <filepath> must be set for https mode" )
	}

	server.session, e = discordgo.New( discordgo.BasicToken( server.id, secret ) )
	if e != nil {
		log.Fatalln( e )
	}

  server.states = make( map[string]int64 )

	if ssl {
		log.Fatalln( http.ListenAndServeTLS( bind, cert, key, server ) )
	} else {
		log.Fatalln( http.ListenAndServe( bind, server ) )
	}
}
