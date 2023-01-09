package discordgo

import (
	"encoding/base64"
	"encoding/json"
	"regexp"
	"strconv"
)

// Get Discords build number
func (s *Session) GetBuildNumber() int {
	defaultBuildNumber := 165485

	loginPage, err := s.Request("GET", "https://discord.com/channels/@me", nil)
	if err != nil {
		return defaultBuildNumber
	}

	javascriptFiles := regexp.MustCompile(`assets/+([a-z0-9]+)\.js`).FindAllStringSubmatch(string(loginPage), -1)
	if len(javascriptFiles) < 2 {
		return defaultBuildNumber
	}

	jsFileWithBuildNumber := javascriptFiles[len(javascriptFiles)-1][0]

	javascript, err := s.Request("GET", "https://discord.com/"+jsFileWithBuildNumber, nil)
	if err != nil {
		return defaultBuildNumber
	}

	buildNumber := regexp.MustCompile(`buildNumber:"(.*?)"`).FindStringSubmatch(string(javascript))
	if len(buildNumber) < 2 {
		return defaultBuildNumber
	}

	buildNumberInt, err := strconv.Atoi(buildNumber[1])
	if err != nil {
		return defaultBuildNumber
	}

	return buildNumberInt
}

// Get super properties
func (s *Session) GetSuperProperties() string {
	superProperties := map[string]any{
		"os":                       s.Identify.Properties.OS,
		"browser":                  s.Identify.Properties.Browser,
		"device":                   s.Identify.Properties.Device,
		"system_locale":            s.Identify.Properties.SystemLocale,
		"browser_user_agent":       s.UserAgent,
		"browser_version":          s.Identify.Properties.BrowserVersion,
		"os_version":               s.Identify.Properties.OSVersion,
		"referrer":                 s.Identify.Properties.Referrer,
		"referring_domain":         s.Identify.Properties.ReferringDomain,
		"referrer_current":         s.Identify.Properties.ReferrerCurrent,
		"referring_domain_current": s.Identify.Properties.ReferringDomainCurrent,
		"release_channel":          s.Identify.Properties.ReleaseChannel,
		"client_build_number":      s.Identify.Properties.ClientBuildNumber,
		"client_event_source":      nil,
	}

	superPropertiesJson, err := json.Marshal(superProperties)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(superPropertiesJson)
}
