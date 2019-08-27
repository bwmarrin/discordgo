package discordgo

// MessageEmbedFooter is a part of a MessageEmbed struct.
type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedImage is a part of a MessageEmbed struct.
type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedThumbnail is a part of a MessageEmbed struct.
type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedVideo is a part of a MessageEmbed struct.
type MessageEmbedVideo struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedProvider is a part of a MessageEmbed struct.
type MessageEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

// MessageEmbedAuthor is a part of a MessageEmbed struct.
type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedField is a part of a MessageEmbed struct.
type MessageEmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// An MessageEmbed stores data for message embeds.
type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []*MessageEmbedField   `json:"fields,omitempty"`
}

// NewEmbed creates an empty MessageEmbed object that you can use to chain
func NewEmbed() *MessageEmbed {
	return &MessageEmbed{}
}

// SetDescription can be used to set the embed description in a chain
// desc :   the embed description
func (e *MessageEmbed) SetDescription(desc string) *MessageEmbed {
	e.Description = desc
	return e
}

// SetTitle can be used to set the embed title in a chain
// title :   the embed title
func (e *MessageEmbed) SetTitle(title string) *MessageEmbed {
	e.Title = title
	return e
}

// SetColor can be used to set the embed color in a chain
// color :   the embed color
func (e *MessageEmbed) SetColor(c int) *MessageEmbed {
	e.Color = c
	return e
}

// SetFooterText can be used to only set the text of the embed footer in a chain
// text :   the footer text
func (e *MessageEmbed) SetFooterText(text string) *MessageEmbed {
	e.Footer = &MessageEmbedFooter{
		Text: text,
	}
	return e
}

// SetFooter can be used to set the text and url of the embed footer in a chain
// text :   the footer text
// url :    the footer url
func (e *MessageEmbed) SetFooter(text, url string) *MessageEmbed {
	e.Footer = &MessageEmbedFooter{
		Text:    text,
		IconURL: url,
	}
	return e
}

// SetAuthorName can be used to only set the name of the author in a chain
// name :   the author name
func (e *MessageEmbed) SetAuthorName(name string) *MessageEmbed {
	e.Author = &MessageEmbedAuthor{
		Name: name,
	}
	return e
}

// SetAuthor can be used to set the name, url and icon url of the author in a chain
// name :     the author name
// url :      the author url
// iconUrl :  the url of the author icon
func (e *MessageEmbed) SetAuthor(name, url, iconUrl string) *MessageEmbed {
	e.Author = &MessageEmbedAuthor{
		Name:    name,
		URL:     url,
		IconURL: iconUrl,
	}
	return e
}

// AddField can be used to add an embed field in a chain
// name :   the field name
// value :  the field value
// inline : determines if the field should be placed inline or not
func (e *MessageEmbed) AddField(name, value string, inline bool) *MessageEmbed {
	e.Fields = append(e.Fields, &MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
	return e
}

// ClearFields removes all fields in a chain
func (e *MessageEmbed) ClearFields() *MessageEmbed {
	e.Fields = nil
	return e
}

// SetFieldAt replaces the field at position with new values in a chain
// name :     the field name
// value :    the field value
// inline :   determines if the field should be placed inline or not
// position : the position of the field to replace
func (e *MessageEmbed) SetFieldAt(name, value string, inline bool, position int) *MessageEmbed {
	if position > len(e.Fields) {
		position = len(e.Fields)
	}
	e.Fields[position] = &MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
	return e
}

// RemoveField removes the field at position in a chain
// position :  the position of the field to remove
func (e *MessageEmbed) RemoveField(position int) *MessageEmbed {
	e.Fields = append(e.Fields[:position], e.Fields[position+1:]...)
	return e
}

// SetImage sets the image url in a chain
// url :  url of the image
func (e *MessageEmbed) SetImage(url string) *MessageEmbed {
	e.Image = &MessageEmbedImage{
		URL: url,
	}
	return e
}

// SetThumbnail sets the image url for the thumbnail in a chain
// url :  url of the image
func (e *MessageEmbed) SetThumbnail(url string) *MessageEmbed {
	e.Thumbnail = &MessageEmbedThumbnail{
		URL: url,
	}
	return e
}
