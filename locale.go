package discordgo

type Locale string

func (l Locale) String() string {
	switch l {
	case LocaleGerman:
		return "German"
	case LocaleEnglishUK:
		return "English, UK"
	case LocaleEnglishUS:
		return "English, US"
	case LocaleSpanish:
		return "Spanish"
	case LocaleFrench:
		return "French"
	case LocaleCroatioan:
		return "Croatioan"
	case LocaleItalian:
		return "Italian"
	case LocaleLithuanian:
		return "Lithuanian"
	case LocaleHungarian:
		return "Hungarian"
	case LocaleDutch:
		return "Dutch"
	case LocaleNorwegian:
		return "Norwegian"
	case LocalePolish:
		return "Polish"
	case LocalePortuguese:
		return "Portuguese, Brazilian"
	case LocaleRomanian:
		return "Romanian, Romania"
	case LocaleFinnish:
		return "Finnish"
	case LocaleSwedish:
		return "Swedish"
	case LocaleVietnamese:
		return "Vietnamese"
	case LocaleTurkish:
		return "Turkish"
	case LocaleCzech:
		return "Czech"
	case LocaleGreek:
		return "Greek"
	case LocaleBulgarian:
		return "Bulgarian"
	case LocaleRussian:
		return "Russian"
	case LocaleUkrainian:
		return "Ukrainian"
	case LocaleHindi:
		return "Hindi"
	case LocaleThai:
		return "Thai"
	case LocaleChinese:
		return "Chinese, China"
	case LocaleJapanese:
		return "Japanese"
	case LocaleTaiwan:
		return "Chinese, Taiwan"
	case LocaleKorean:
		return "Korean"
	}

	return "Unknown"
}

// All known locales
const (
	LocaleGerman                      Locale = "de"
	LocaleEnglishUK                   Locale = "en-GB"
	LocaleEnglishUS                   Locale = "en-US"
	LocaleSpanish                     Locale = "es-ES"
	LocaleFrench                      Locale = "fr"
	LocaleCroatioan                   Locale = "hr"
	LocaleItalian                     Locale = "it"
	LocaleLithuanian                  Locale = "lt"
	LocaleHungarian                   Locale = "hu"
	LocaleDutch                       Locale = "nl"
	LocaleNorwegian                   Locale = "no"
	LocalePolish                      Locale = "pl"
	LocalePortuguese, LocaleBrazilian Locale = "pt-BR", "pt-BR"
	LocaleRomanian, LocaleRomania     Locale = "ro", "ro"
	LocaleFinnish                     Locale = "fi"
	LocaleSwedish                     Locale = "sv-SE"
	LocaleVietnamese                  Locale = "vi"
	LocaleTurkish                     Locale = "tr"
	LocaleCzech                       Locale = "cs"
	LocaleGreek                       Locale = "el"
	LocaleBulgarian                   Locale = "bg"
	LocaleRussian                     Locale = "ru"
	LocaleUkrainian                   Locale = "uk"
	LocaleHindi                       Locale = "hi"
	LocaleThai                        Locale = "th"
	LocaleChinese, LocaleChina        Locale = "zh-CN", "zh-CN"
	LocaleJapanese                    Locale = "ja"
	LocaleTaiwan                      Locale = "zh-TW"
	LocaleKorean                      Locale = "ko"
	LocaleUnknown                     Locale = ""
)
