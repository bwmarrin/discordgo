package discordgo

import (
	"fmt"
	"net/url"
	"time"
)

const (
	ScopeActivitiesRead                         = "activities.read"
	ScopeActivitieswrite                        = "activities.write"
	ScopeApplicationsBuildsRead                 = "applications.builds.read"
	ScopeApplicationsBuildsUpload               = "applications.builds.upload"
	ScopeApplicationsCommands                   = "applications.commands"
	ScopeApplicationsCommandsUpdate             = "applications.commands.update"
	ScopeApplicationsCommandsPermissionsUpdate  = "applications.commands.permissions.update"
	ScopeApplicationsEntitlements               = "applications.entitlements"
	ScopeApplicationsStoreUpdate                = "applications.store.update"
	ScopeBot                                    = "bot"
	Scopeconnections                            = "connections"
	ScopeDM_ChannelsRead                        = "dm_channels.read"
	ScopeEmail                                  = "email"
	ScopeGDM_Join                               = "gdm.join"
	ScopeGuilds                                 = "guilds"
	ScopeGuildsJoin                             = "guilds.join"
	ScopeGuildsMembersRead                      = "guilds.members.read"
	ScopeIdentify                               = "identify"
	ScopeMessagesRead                           = "messages.read"
	ScopeRelationshipsRead                      = "relationships.read"
	ScopeRoleConnectionsWrite                   = "role_connections.write"
	ScopeRPC                                    = "rpc"
	ScopeRPC_ActivitiesWrite                    = "rpc.activities.write"
	ScopeRPC_NotificationsRead                  = "rpc.notifications.read"
	ScopeRPC_VoiceRead                          = "rpc.voice.read"
	ScopeRPC_VoiceWrite                         = "rpc.voice.write"
	ScopeVoice                                  = "voice"
	ScopeWebhookIncoming                        = "webhook.incoming"
)

// Token hint type for use in oauth2 token revocation requests
type TokenHint string
const (
	// No token hint given for an oauth2 token revocation request
	TokenHintNone     TokenHint = ""
	// Token given for an oauth2 token revocation request is ( likely ) an access token
	TokenHintAccess   TokenHint = "access_token"
	// Token given for an oauth2 token revocation requsst is ( likely ) a refresh token
	TokenHintRefresh  TokenHint = "refresh_token"
)

// An oauth2 access token response from the Discord Api
type AccessToken struct {
	// Oauth2 access token for accessing the Discord Api
	AccessToken   string    `json:"access_token"`
	// Oauth2 token type either: Bearer or Bot
	TokenType     string    `json:"token_type"`
	// Expire time for the access in seconds from the time the response was sent
	ExpiresIn     int64     `json:"expires_in"`

	expiresAt     int64
	
	// Refresh token for requesting a new oauth2 access token ( empty for client credentials )
	RefreshToken  string    `json:"refresh_token,omitempty"`
	// Scope for the token
	Scope         string    `json:"scope"`
	// Guild object for advanced bot authorization scoped tokens ( nil for others )
	Guild         *Guild    `json:"guild,omitempty"`
	// Webhook object for webhook tokens ( nil for others )
	Webhook       *Webhook  `json:"webhook,omitempty"`
}

// Creates a new Discord Api session using the token from an access token response
func ( at *AccessToken ) NewSession ( ) ( st *Session, err error ) {
	if at.TokenType != "Bearer" { return nil, ErrNotBearerToken }
	return New( at.TokenType + " " + at.AccessToken )
}

// Expires returns a time.Time object representing the true expiration time of an access token and whether it is known ( unknown unless the AccessToken object was created by a response from the Discord Api )
func ( at *AccessToken ) Expires ( ) ( t time.Time, has bool ) {
	if at.expiresAt == 0 { return time.Unix( 0, 0 ), false }
	return time.Unix( at.expiresAt, 0 ), true
}

// Oauth2 authorization info
type AuthorizationInfo struct {
	// Application object
	Application   *Application  `json:"application"`
	// Scopes authorized
	Scopes        []string      `json:"scopes"`
	// Expire date/time in the ISO 8601 time format for the token
	Expires       string        `json:"expires"`
	// User object
	User          *User         `json:"user"`
}

func ( s *Session ) accessToken ( data url.Values, options ...RequestOption ) ( st *AccessToken, err error ) {
	if err = s.checkBasicSession(); err != nil { return }
	body, err := s.RequestWithBucketID( "POST", EndpointOAuth2Token, data, EndpointOAuth2Token, options... )
	if err == nil {
		err = unmarshal( body, &st )
	}
	st.expiresAt = time.Now().UTC().Unix() + st.ExpiresIn
	return
}

// AccessToken performs an oauth2 token request for the code and returns the access token response ( requires a session created with a 'Basic' token using the applications client id and client secret )
// code       :  code obtained from Discord oauth2 redirect
// redirectUri:  redirect uri used in the Discord oauth2 redirect
func ( s *Session ) AccessToken ( code string, redirectUri string, options ...RequestOption ) ( st *AccessToken, err error ) {
	data := url.Values{
		"grant_type": []string{ "authorization_code" },
		"code": []string{ code },
		"redirect_uri": []string{ redirectUri },
	}
	return s.accessToken( data, options... )
}

// AccessTokenRefresh performs an oauth2 token request using the refresh token of a previous access token response ( requires a session created with a 'Basic' token using the applications client id and client secret; and not one created from the 'Bearer'/'Bot' token of the user )
// refreshToken:  refresh token obtained from a prior oauth2 token response
func ( s *Session ) AccessTokenRefresh ( refreshToken string, options ...RequestOption ) ( st *AccessToken, err error ) {
	data := url.Values{
		"grant_type": []string{ "refresh_token" },
		"refresh_token": []string{ refreshToken },
	}
	return s.accessToken( data, options... )
}

// ClientCredentials performs a client credentials request using scopes retreiving developers oauth2 token for testing purposes. ( requires a session created with a 'Basic' token using the applications client id and client secret )
// scope:  the scope to request the token be issued for
func ( s *Session ) ClientCredentials ( scope string, options ...RequestOption ) ( st *AccessToken, err error ) {
	data := url.Values{
		"grant_type": []string{ "client_credentials" },
		"scope": []string{ scope },
	}
	return s.accessToken( data, options... )
}

// AccessTokenRevoke preforms a oauth2 token revocation request for token ( requires a session created with a 'Basic' token using the applications client id and client secret; and not one created from the 'Bearer'/'Bot' token )
// token    :  the access token or that access token's refresh token being revoked
// tokenType:  the hint for which type of token has been provided to search for
func ( s *Session ) AccessTokenRevoke ( token string, tokenType TokenHint, options ...RequestOption ) ( st []byte, err error ) {
	if err = s.checkBasicSession( ); err != nil { return }
	data := url.Values{ "token": []string{ token } }
	if tokenType != TokenHintNone {
		data.Set( "token_type_hint", string( tokenType ) )
	}
	st, err = s.RequestWithBucketID( "POST", EndpointOAuth2TokenRevoke, data, EndpointOAuth2TokenRevoke, options... )
	return
}

// AuthorizationInfo retreives the authorization information of the access token ( requires a session created with the 'Bearer' token received from a previous oauth2 access token response; and not a 'Basic' token created with the applications client id and client secret )
func ( s *Session ) AuthorizationInfo ( options ...RequestOption ) ( st *AuthorizationInfo, err error ) {
	if err = s.checkBearerSession( ); err != nil { return }
	fmt.Println(s.Identify.Token)
	body, err := s.RequestWithBucketID( "GET", EndpointOAuth2AuthorizationInfo, nil, EndpointOAuth2AuthorizationInfo, options... )
	if err == nil { err = unmarshal( body, &st ) }
	return
}
