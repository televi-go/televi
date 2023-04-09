package dto

type Animation struct {
	// FileID is the identifier for this file, which can be used to download or reuse
	// the file
	FileID string `json:"file_id"`
	// FileUniqueID is the unique identifier for this file, which is supposed to
	// be the same over time and for different bots. Can't be used to download
	// or reuse the file.
	FileUniqueID string `json:"file_unique_id"`
	// Width video width as defined by sender
	Width int `json:"width"`
	// Height video height as defined by sender
	Height int `json:"height"`
	// Duration of the video in seconds as defined by sender
	Duration int `json:"duration"`
	// Thumbnail animation thumbnail as defined by sender
	//
	// optional
	Thumbnail *PhotoSize `json:"thumb,omitempty"`
	// FileName original animation filename as defined by sender
	//
	// optional
	FileName string `json:"file_name,omitempty"`
	// MimeType of the file as defined by sender
	//
	// optional
	MimeType string `json:"mime_type,omitempty"`
	// FileSize file size
	//
	// optional
	FileSize int `json:"file_size,omitempty"`
}

type PhotoSize struct {
	// FileID identifier for this file, which can be used to download or reuse
	// the file
	FileID string `json:"file_id"`
	// FileUniqueID is the unique identifier for this file, which is supposed to
	// be the same over time and for different bots. Can't be used to download
	// or reuse the file.
	FileUniqueID string `json:"file_unique_id"`
	// Width photo width
	Width int `json:"width"`
	// Height photo height
	Height int `json:"height"`
	// FileSize file size
	//
	// optional
	FileSize int `json:"file_size,omitempty"`
}

type Audio struct {
	// FileID is an identifier for this file, which can be used to download or
	// reuse the file
	FileID string `json:"file_id"`
	// FileUniqueID is the unique identifier for this file, which is supposed to
	// be the same over time and for different bots. Can't be used to download
	// or reuse the file.
	FileUniqueID string `json:"file_unique_id"`
	// Duration of the audio in seconds as defined by sender
	Duration int `json:"duration"`
	// Performer of the audio as defined by sender or by audio tags
	//
	// optional
	Performer string `json:"performer,omitempty"`
	// Title of the audio as defined by sender or by audio tags
	//
	// optional
	Title string `json:"title,omitempty"`
	// FileName is the original filename as defined by sender
	//
	// optional
	FileName string `json:"file_name,omitempty"`
	// MimeType of the file as defined by sender
	//
	// optional
	MimeType string `json:"mime_type,omitempty"`
	// FileSize file size
	//
	// optional
	FileSize int `json:"file_size,omitempty"`
	// Thumbnail is the album cover to which the music file belongs
	//
	// optional
	Thumbnail *PhotoSize `json:"thumb,omitempty"`
}

type Document struct {
	// FileID is an identifier for this file, which can be used to download or
	// reuse the file
	FileID string `json:"file_id"`
	// FileUniqueID is the unique identifier for this file, which is supposed to
	// be the same over time and for different bots. Can't be used to download
	// or reuse the file.
	FileUniqueID string `json:"file_unique_id"`
	// Thumbnail document thumbnail as defined by sender
	//
	// optional
	Thumbnail *PhotoSize `json:"thumb,omitempty"`
	// FileName original filename as defined by sender
	//
	// optional
	FileName string `json:"file_name,omitempty"`
	// MimeType  of the file as defined by sender
	//
	// optional
	MimeType string `json:"mime_type,omitempty"`
	// FileSize file size
	//
	// optional
	FileSize int `json:"file_size,omitempty"`
}

type Video struct {
	// FileID identifier for this file, which can be used to download or reuse
	// the file
	FileID string `json:"file_id"`
	// FileUniqueID is the unique identifier for this file, which is supposed to
	// be the same over time and for different bots. Can't be used to download
	// or reuse the file.
	FileUniqueID string `json:"file_unique_id"`
	// Width video width as defined by sender
	Width int `json:"width"`
	// Height video height as defined by sender
	Height int `json:"height"`
	// Duration of the video in seconds as defined by sender
	Duration int `json:"duration"`
	// Thumbnail video thumbnail
	//
	// optional
	Thumbnail *PhotoSize `json:"thumb,omitempty"`
	// FileName is the original filename as defined by sender
	//
	// optional
	FileName string `json:"file_name,omitempty"`
	// MimeType of a file as defined by sender
	//
	// optional
	MimeType string `json:"mime_type,omitempty"`
	// FileSize file size
	//
	// optional
	FileSize int `json:"file_size,omitempty"`
}

// VideoNote object represents a video message.
type VideoNote struct {
	// FileID identifier for this file, which can be used to download or reuse the file
	FileID string `json:"file_id"`
	// FileUniqueID is the unique identifier for this file, which is supposed to
	// be the same over time and for different bots. Can't be used to download
	// or reuse the file.
	FileUniqueID string `json:"file_unique_id"`
	// Length video width and height (diameter of the video message) as defined by sender
	Length int `json:"length"`
	// Duration of the video in seconds as defined by sender
	Duration int `json:"duration"`
	// Thumbnail video thumbnail
	//
	// optional
	Thumbnail *PhotoSize `json:"thumb,omitempty"`
	// FileSize file size
	//
	// optional
	FileSize int `json:"file_size,omitempty"`
}

// Voice represents a voice note.
type Voice struct {
	// FileID identifier for this file, which can be used to download or reuse the file
	FileID string `json:"file_id"`
	// FileUniqueID is the unique identifier for this file, which is supposed to
	// be the same over time and for different bots. Can't be used to download
	// or reuse the file.
	FileUniqueID string `json:"file_unique_id"`
	// Duration of the audio in seconds as defined by sender
	Duration int `json:"duration"`
	// MimeType of the file as defined by sender
	//
	// optional
	MimeType string `json:"mime_type,omitempty"`
	// FileSize file size
	//
	// optional
	FileSize int `json:"file_size,omitempty"`
}

// Contact represents a phone contact.
//
// Note that LastName and UserID may be empty.
type Contact struct {
	// PhoneNumber contact's phone number
	PhoneNumber string `json:"phone_number"`
	// FirstName contact's first name
	FirstName string `json:"first_name"`
	// LastName contact's last name
	//
	// optional
	LastName string `json:"last_name,omitempty"`
	// UserID contact's user identifier in Telegram
	//
	// optional
	UserID int64 `json:"user_id,omitempty"`
	// VCard is additional data about the contact in the form of a vCard.
	//
	// optional
	VCard string `json:"vcard,omitempty"`
}

// Dice represents an animated emoji that displays a random value.
type Dice struct {
	// Emoji on which the dice throw animation is based
	Emoji string `json:"emoji"`
	// Value of the dice
	Value int `json:"value"`
}

type Sticker struct {
	// FileID is an identifier for this file, which can be used to download or
	// reuse the file
	FileID string `json:"file_id"`
	// FileUniqueID is a unique identifier for this file,
	// which is supposed to be the same over time and for different bots.
	// Can't be used to download or reuse the file.
	FileUniqueID string `json:"file_unique_id"`
	// Width sticker width
	Width int `json:"width"`
	// Height sticker height
	Height int `json:"height"`
	// IsAnimated true, if the sticker is animated
	//
	// optional
	IsAnimated bool `json:"is_animated,omitempty"`
	// IsVideo true, if the sticker is a video sticker
	//
	// optional
	IsVideo bool `json:"is_video,omitempty"`
	// Thumbnail sticker thumbnail in the .WEBP or .JPG format
	//
	// optional
	Thumbnail *PhotoSize `json:"thumb,omitempty"`
	// Emoji associated with the sticker
	//
	// optional
	Emoji string `json:"emoji,omitempty"`
	// SetName of the sticker set to which the sticker belongs
	//
	// optional
	SetName string `json:"set_name,omitempty"`
	// PremiumAnimation for premium regular stickers, premium animation for the sticker
	//
	// optional
	//PremiumAnimation *File `json:"premium_animation,omitempty"`
	// MaskPosition is for mask stickers, the position where the mask should be
	// placed
	//
	// optional
	//MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	// CustomEmojiID for custom emoji stickers, unique identifier of the custom emoji
	//
	// optional
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
	// FileSize
	//
	// optional
	FileSize int `json:"file_size,omitempty"`
}

type Poll struct {
	// ID is the unique poll identifier
	ID string `json:"id"`
	// Question is the poll question, 1-255 characters
	Question string `json:"question"`
	// Options is the list of poll options
	Options []PollOption `json:"options"`
	// TotalVoterCount is the total numbers of users who voted in the poll
	TotalVoterCount int `json:"total_voter_count"`
	// IsClosed is if the poll is closed
	IsClosed bool `json:"is_closed"`
	// IsAnonymous is if the poll is anonymous
	IsAnonymous bool `json:"is_anonymous"`
	// Type is the poll type, currently can be "regular" or "quiz"
	Type string `json:"type"`
	// AllowsMultipleAnswers is true, if the poll allows multiple answers
	AllowsMultipleAnswers bool `json:"allows_multiple_answers"`
	// CorrectOptionID is the 0-based identifier of the correct answer option.
	// Available only for polls in quiz mode, which are closed, or was sent (not
	// forwarded) by the bot or to the private chat with the bot.
	//
	// optional
	CorrectOptionID int `json:"correct_option_id,omitempty"`
	// Explanation is text that is shown when a user chooses an incorrect answer
	// or taps on the lamp icon in a quiz-style poll, 0-200 characters
	//
	// optional
	Explanation string `json:"explanation,omitempty"`
	// ExplanationEntities are special entities like usernames, URLs, bot
	// commands, etc. that appear in the explanation
	//
	// optional
	ExplanationEntities []MessageEntity `json:"explanation_entities,omitempty"`
	// OpenPeriod is the amount of time in seconds the poll will be active
	// after creation
	//
	// optional
	OpenPeriod int `json:"open_period,omitempty"`
	// CloseDate is the point in time (unix timestamp) when the poll will be
	// automatically closed
	//
	// optional
	CloseDate int `json:"close_date,omitempty"`
}

type PollOption struct {
	// Text is the option text, 1-100 characters
	Text string `json:"text"`
	// VoterCount is the number of users that voted for this option
	VoterCount int `json:"voter_count"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll.
type PollAnswer struct {
	// PollID is the unique poll identifier
	PollID string `json:"poll_id"`
	// User who changed the answer to the poll
	User User `json:"user"`
	// OptionIDs is the 0-based identifiers of poll options chosen by the user.
	// May be empty if user retracted vote.
	OptionIDs []int `json:"option_ids"`
}

type Venue struct {
	// Location is the venue location
	Location Location `json:"location"`
	// Title is the name of the venue
	Title string `json:"title"`
	// Address of the venue
	Address string `json:"address"`
	// FoursquareID is the foursquare identifier of the venue
	//
	// optional
	FoursquareID string `json:"foursquare_id,omitempty"`
	// FoursquareType is the foursquare type of the venue
	//
	// optional
	FoursquareType string `json:"foursquare_type,omitempty"`
	// GooglePlaceID is the Google Places identifier of the venue
	//
	// optional
	GooglePlaceID string `json:"google_place_id,omitempty"`
	// GooglePlaceType is the Google Places type of the venue
	//
	// optional
	GooglePlaceType string `json:"google_place_type,omitempty"`
}
