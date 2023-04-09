package dto

type Message struct {
	// MessageID is a unique message identifier inside this chat
	MessageID int `json:"message_id"`
	// From is a sender, empty for messages sent to channels;
	//
	// optional
	From *User `json:"from,omitempty"`
	// SenderChat is the sender of the message, sent on behalf of a chat. The
	// channel itself for channel messages. The supergroup itself for messages
	// from anonymous group administrators. The linked channel for messages
	// automatically forwarded to the discussion group
	//
	// optional
	SenderChat *Chat `json:"sender_chat,omitempty"`
	// Date of the message was sent in Unix time
	Date int `json:"date"`
	// Chat is the conversation the message belongs to
	Chat *Chat `json:"chat"`
	// ForwardFrom for forwarded messages, sender of the original message;
	//
	// optional
	ForwardFrom *User `json:"forward_from,omitempty"`
	// ForwardFromChat for messages forwarded from channels,
	// information about the original channel;
	//
	// optional
	ForwardFromChat *Chat `json:"forward_from_chat,omitempty"`
	// ForwardFromMessageID for messages forwarded from channels,
	// identifier of the original message in the channel;
	//
	// optional
	ForwardFromMessageID int `json:"forward_from_message_id,omitempty"`
	// ForwardSignature for messages forwarded from channels, signature of the
	// post author if present
	//
	// optional
	ForwardSignature string `json:"forward_signature,omitempty"`
	// ForwardSenderName is the sender's name for messages forwarded from users
	// who disallow adding a link to their account in forwarded messages
	//
	// optional
	ForwardSenderName string `json:"forward_sender_name,omitempty"`
	// ForwardDate for forwarded messages, date the original message was sent in Unix time;
	//
	// optional
	ForwardDate int `json:"forward_date,omitempty"`
	// IsAutomaticForward is true if the message is a channel post that was
	// automatically forwarded to the connected discussion group.
	//
	// optional
	IsAutomaticForward bool `json:"is_automatic_forward,omitempty"`
	// ReplyToMessage for replies, the original message.
	// Note that the Message object in this field will not contain further ReplyToMessage fields
	// even if it itself is a reply;
	//
	// optional
	ReplyToMessage *Message `json:"reply_to_message,omitempty"`
	// ViaBot through which the message was sent;
	//
	// optional
	ViaBot *User `json:"via_bot,omitempty"`
	// EditDate of the message was last edited in Unix time;
	//
	// optional
	EditDate int `json:"edit_date,omitempty"`
	// HasProtectedContent is true if the message can't be forwarded.
	//
	// optional
	HasProtectedContent bool `json:"has_protected_content,omitempty"`
	// MediaGroupID is the unique identifier of a media message group this message belongs to;
	//
	// optional
	MediaGroupID string `json:"media_group_id,omitempty"`
	// AuthorSignature is the signature of the post author for messages in channels;
	//
	// optional
	AuthorSignature string `json:"author_signature,omitempty"`
	// Text is for text messages, the actual UTF-8 text of the message, 0-4096 characters;
	//
	// optional
	Text string `json:"text,omitempty"`
	// Entities are for text messages, special entities like usernames,
	// URLs, bot commands, etc. that appear in the text;
	//
	// optional
	Entities []MessageEntity `json:"entities,omitempty"`
	// Animation message is an animation, information about the animation.
	// For backward compatibility, when this field is set, the document field will also be set;
	//
	// optional
	Animation *Animation `json:"animation,omitempty"`
	// PremiumAnimation message is an animation, information about the animation.
	// For backward compatibility, when this field is set, the document field will also be set;
	//
	// optional
	PremiumAnimation *Animation `json:"premium_animation,omitempty"`
	// Audio message is an audio file, information about the file;
	//
	// optional
	Audio *Audio `json:"audio,omitempty"`
	// Document message is a general file, information about the file;
	//
	// optional
	Document *Document `json:"document,omitempty"`
	// Photo message is a photo, available sizes of the photo;
	//
	// optional
	Photo []PhotoSize `json:"photo,omitempty"`
	// Sticker message is a sticker, information about the sticker;
	//
	// optional
	Sticker *Sticker `json:"sticker,omitempty"`
	// Video message is a video, information about the video;
	//
	// optional
	Video *Video `json:"video,omitempty"`
	// VideoNote message is a video note, information about the video message;
	//
	// optional
	VideoNote *VideoNote `json:"video_note,omitempty"`
	// Voice message is a voice message, information about the file;
	//
	// optional
	Voice *Voice `json:"voice,omitempty"`
	// Caption for the animation, audio, document, photo, video or voice, 0-1024 characters;
	//
	// optional
	Caption string `json:"caption,omitempty"`
	// CaptionEntities;
	//
	// optional
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	// Contact message is a shared contact, information about the contact;
	//
	// optional
	Contact *Contact `json:"contact,omitempty"`
	// Dice is a dice with random value;
	//
	// optional
	Dice *Dice `json:"dice,omitempty"`
	// Game message is a game, information about the game;
	//
	// optional
	//TODO:Game *Game `json:"game,omitempty"`
	// Poll is a native poll, information about the poll;
	//
	// optional
	Poll *Poll `json:"poll,omitempty"`
	// Venue message is a venue, information about the venue.
	// For backward compatibility, when this field is set, the location field
	// will also be set;
	//
	// optional
	Venue *Venue `json:"venue,omitempty"`
	// Location message is a shared location, information about the location;
	//
	// optional
	Location *Location `json:"location,omitempty"`
	// NewChatMembers that were added to the group or supergroup
	// and information about them (the bot itself may be one of these members);
	//
	// optional
	NewChatMembers []User `json:"new_chat_members,omitempty"`
	// LeftChatMember is a member was removed from the group,
	// information about them (this member may be the bot itself);
	//
	// optional
	LeftChatMember *User `json:"left_chat_member,omitempty"`
	// NewChatTitle is a chat title was changed to this value;
	//
	// optional
	NewChatTitle string `json:"new_chat_title,omitempty"`
	// NewChatPhoto is a chat photo was change to this value;
	//
	// optional
	NewChatPhoto []PhotoSize `json:"new_chat_photo,omitempty"`
	// DeleteChatPhoto is a service message: the chat photo was deleted;
	//
	// optional
	DeleteChatPhoto bool `json:"delete_chat_photo,omitempty"`
	// GroupChatCreated is a service message: the group has been created;
	//
	// optional
	GroupChatCreated bool `json:"group_chat_created,omitempty"`
	// SuperGroupChatCreated is a service message: the supergroup has been created.
	// This field can't be received in a message coming through updates,
	// because bot can't be a member of a supergroup when it is created.
	// It can only be found in ReplyToMessage if someone replies to a very first message
	// in a directly created supergroup;
	//
	// optional
	SuperGroupChatCreated bool `json:"supergroup_chat_created,omitempty"`
	// ChannelChatCreated is a service message: the channel has been created.
	// This field can't be received in a message coming through updates,
	// because bot can't be a member of a channel when it is created.
	// It can only be found in ReplyToMessage
	// if someone replies to a very first message in a channel;
	//
	// optional
	ChannelChatCreated bool `json:"channel_chat_created,omitempty"`
	// MessageAutoDeleteTimerChanged is a service message: auto-delete timer
	// settings changed in the chat.
	//
	// optional
	//TODO:MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	// MigrateToChatID is the group has been migrated to a supergroup with the specified identifier.
	// This number may be greater than 32 bits and some programming languages
	// may have difficulty/silent defects in interpreting it.
	// But it is smaller than 52 bits, so a signed 64-bit integer
	// or double-precision float type are safe for storing this identifier;
	//
	// optional
	MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`
	// MigrateFromChatID is the supergroup has been migrated from a group with the specified identifier.
	// This number may be greater than 32 bits and some programming languages
	// may have difficulty/silent defects in interpreting it.
	// But it is smaller than 52 bits, so a signed 64-bit integer
	// or double-precision float type are safe for storing this identifier;
	//
	// optional
	MigrateFromChatID int64 `json:"migrate_from_chat_id,omitempty"`
	// PinnedMessage is a specified message was pinned.
	// Note that the Message object in this field will not contain further ReplyToMessage
	// fields even if it is itself a reply;
	//
	// optional
	PinnedMessage *Message `json:"pinned_message,omitempty"`
	// Invoice message is an invoice for a payment;
	//
	// optional
	//TODO:Invoice *Invoice `json:"invoice,omitempty"`
	// SuccessfulPayment message is a service message about a successful payment,
	// information about the payment;
	//
	// optional
	//TODO:SuccessfulPayment *SuccessfulPayment `json:"successful_payment,omitempty"`
	// ConnectedWebsite is the domain name of the website on which the user has
	// logged in;
	//
	// optional
	ConnectedWebsite string `json:"connected_website,omitempty"`
	// PassportData is a Telegram Passport data;
	//
	// optional
	//TODO:PassportData *PassportData `json:"passport_data,omitempty"`

	// ProximityAlertTriggered is a service message. A user in the chat
	// triggered another user's proximity alert while sharing Live Location
	//
	// optional
	ProximityAlertTriggered *ProximityAlertTriggered `json:"proximity_alert_triggered,omitempty"`
	// VideoChatScheduled is a service message: video chat scheduled.
	//
	// optional
	//TODO:VideoChatScheduled *VideoChatScheduled `json:"video_chat_scheduled,omitempty"`
	// VideoChatStarted is a service message: video chat started.
	//
	// optional
	//TODO:VideoChatStarted *VideoChatStarted `json:"video_chat_started,omitempty"`
	// VideoChatEnded is a service message: video chat ended.
	//
	// optional
	//TODO:VideoChatEnded *VideoChatEnded `json:"video_chat_ended,omitempty"`
	// VideoChatParticipantsInvited is a service message: new participants
	// invited to a video chat.
	//
	// optional
	//TODO:VideoChatParticipantsInvited *VideoChatParticipantsInvited `json:"video_chat_participants_invited,omitempty"`
	// WebAppData is a service message: data sent by a Web App.
	//
	// optional
	//WebAppData *WebAppData `json:"web_app_data,omitempty"`
	// ReplyMarkup is the Inline keyboard attached to the message.
	// login_url buttons are represented as ordinary url buttons.
	//
	// optional
	//ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type MessageEntity struct {
	// Type of the entity.
	// Can be:
	//  “mention” (@username),
	//  “hashtag” (#hashtag),
	//  “cashtag” ($USD),
	//  “bot_command” (/start@jobs_bot),
	//  “url” (https://telegram.org),
	//  “email” (do-not-reply@telegram.org),
	//  “phone_number” (+1-212-555-0123),
	//  “bold” (bold text),
	//  “italic” (italic text),
	//  “underline” (underlined text),
	//  “strikethrough” (strikethrough text),
	//  "spoiler" (spoiler message),
	//  “code” (monowidth string),
	//  “pre” (monowidth block),
	//  “text_link” (for clickable text URLs),
	//  “text_mention” (for users without usernames)
	Type string `json:"type"`
	// Offset in UTF-16 code units to the start of the entity
	Offset int `json:"offset"`
	// Length
	Length int `json:"length"`
	// URL for “text_link” only, url that will be opened after user taps on the text
	//
	// optional
	URL string `json:"url,omitempty"`
	// User for “text_mention” only, the mentioned user
	//
	// optional
	User *User `json:"user,omitempty"`
	// Language for “pre” only, the programming language of the entity text
	//
	// optional
	Language string `json:"language,omitempty"`
}

type ProximityAlertTriggered struct {
	// Traveler is the user that triggered the alert
	Traveler User `json:"traveler"`
	// Watcher is the user that set the alert
	Watcher User `json:"watcher"`
	// Distance is the distance between the users
	Distance int `json:"distance"`
}
