package dto

type Update struct {
	// UpdateID is the update's unique identifier.
	// Update identifiers start from a certain positive number and increase
	// sequentially.
	// This ID becomes especially handy if you're using Webhooks,
	// since it allows you to ignore repeated updates or to restore
	// the correct update sequence, should they get out of order.
	// If there are no new updates for at least a week, then identifier
	// of the next update will be chosen randomly instead of sequentially.
	UpdateID int `json:"update_id"`
	// Message new incoming message of any kind — text, photo, sticker, etc.
	//
	// optional
	Message *Message `json:"message,omitempty"`
	// EditedMessage new version of a message that is known to the bot and was
	// edited
	//
	// optional
	EditedMessage *Message `json:"edited_message,omitempty"`
	// ChannelPost new version of a message that is known to the bot and was
	// edited
	//
	// optional
	ChannelPost *Message `json:"channel_post,omitempty"`
	// EditedChannelPost new incoming channel post of any kind — text, photo,
	// sticker, etc.
	//
	// optional
	EditedChannelPost *Message `json:"edited_channel_post,omitempty"`
	// InlineQuery new incoming inline query
	//
	// optional
	InlineQuery *InlineQuery `json:"inline_query,omitempty"`
	// ChosenInlineResult is the result of an inline query
	// that was chosen by a user and sent to their chat partner.
	// Please see our documentation on the feedback collecting
	// for details on how to enable these updates for your bot.
	//
	// optional
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	// CallbackQuery new incoming callback query
	//
	// optional
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
	// ShippingQuery new incoming shipping query. Only for invoices with
	// flexible price
	//
	// optional
	ShippingQuery *ShippingQuery `json:"shipping_query,omitempty"`
	// PreCheckoutQuery new incoming pre-checkout query. Contains full
	// information about checkout
	//
	// optional
	PreCheckoutQuery *PreCheckoutQuery `json:"pre_checkout_query,omitempty"`
	// Pool new poll state. Bots receive only updates about stopped polls and
	// polls, which are sent by the bot
	//
	// optional
	Poll *Poll `json:"poll,omitempty"`
	// PollAnswer user changed their answer in a non-anonymous poll. Bots
	// receive new votes only in polls that were sent by the bot itself.
	//
	// optional
	PollAnswer *PollAnswer `json:"poll_answer,omitempty"`
	// MyChatMember is the bot chat member status was updated in a chat. For
	// private chats, this update is received only when the bot is blocked or
	// unblocked by the user.
	//
	// optional
	MyChatMember *ChatMemberUpdated `json:"my_chat_member,omitempty"`
	// ChatMember is a chat member's status was updated in a chat. The bot must
	// be an administrator in the chat and must explicitly specify "chat_member"
	// in the list of allowed_updates to receive these updates.
	//
	// optional
	ChatMember *ChatMemberUpdated `json:"chat_member,omitempty"`
	// ChatJoinRequest is a request to join the chat has been sent. The bot must
	// have the can_invite_users administrator right in the chat to receive
	// these updates.
	//
	// optional
	ChatJoinRequest *ChatJoinRequest `json:"chat_join_request,omitempty"`
}

// SentFrom returns the user who sent an update. Can be nil, if Telegram did not provide information
// about the user in the update object.
func (u *Update) SentFrom() *User {
	switch {
	case u.Message != nil:
		return u.Message.From
	case u.EditedMessage != nil:
		return u.EditedMessage.From
	case u.InlineQuery != nil:
		return u.InlineQuery.From
	case u.ChosenInlineResult != nil:
		return u.ChosenInlineResult.From
	case u.CallbackQuery != nil:
		return u.CallbackQuery.From
	case u.ShippingQuery != nil:
		return u.ShippingQuery.From
	case u.PreCheckoutQuery != nil:
		return u.PreCheckoutQuery.From
	default:
		return nil
	}
}

type InlineQuery struct {
	// ID unique identifier for this query
	ID string `json:"id"`
	// From sender
	From *User `json:"from"`
	// Query text of the query (up to 256 characters).
	Query string `json:"query"`
	// Offset of the results to be returned, can be controlled by the bot.
	Offset string `json:"offset"`
	// Type of the chat, from which the inline query was sent. Can be either
	// “sender” for a private chat with the inline query sender, “private”,
	// “group”, “supergroup”, or “channel”. The chat type should be always known
	// for requests sent from official clients and most third-party clients,
	// unless the request was sent from a secret chat
	//
	// optional
	ChatType string `json:"chat_type,omitempty"`
	// Location sender location, only for bots that request user location.
	//
	// optional
	Location *Location `json:"location,omitempty"`
}

type ChosenInlineResult struct {
	// ResultID the unique identifier for the result that was chosen
	ResultID string `json:"result_id"`
	// From the user that chose the result
	From *User `json:"from"`
	// Location sender location, only for bots that require user location
	//
	// optional
	Location *Location `json:"location,omitempty"`
	// InlineMessageID identifier of the sent inline message.
	// Available only if there is an inline keyboard attached to the message.
	// Will be also received in callback queries and can be used to edit the message.
	//
	// optional
	InlineMessageID string `json:"inline_message_id,omitempty"`
	// Query the query that was used to obtain the result
	Query string `json:"query"`
}

type CallbackQuery struct {
	// ID unique identifier for this query
	ID string `json:"id"`
	// From sender
	From *User `json:"from"`
	// Message with the callback button that originated the query.
	// Note that message content and message date will not be available if the
	// message is too old.
	//
	// optional
	Message *Message `json:"message,omitempty"`
	// InlineMessageID identifier of the message sent via the bot in inline
	// mode, that originated the query.
	//
	// optional
	InlineMessageID string `json:"inline_message_id,omitempty"`
	// ChatInstance global identifier, uniquely corresponding to the chat to
	// which the message with the callback button was sent. Useful for high
	// scores in games.
	ChatInstance string `json:"chat_instance"`
	// Data associated with the callback button. Be aware that
	// a bad client can send arbitrary data in this field.
	//
	// optional
	Data string `json:"data,omitempty"`
	// GameShortName short name of a Game to be returned, serves as the unique identifier for the game.
	//
	// optional
	GameShortName string `json:"game_short_name,omitempty"`
}

type PreCheckoutQuery struct {
	// ID unique query identifier
	ID string `json:"id"`
	// From user who sent the query
	From *User `json:"from"`
	// Currency three-letter ISO 4217 currency code
	//	// (see https://core.telegram.org/bots/payments#supported-currencies)
	Currency string `json:"currency"`
	// TotalAmount total price in the smallest units of the currency (integer, not float/double).
	//	// For example, for a price of US$ 1.45 pass amount = 145.
	//	// See the exp parameter in currencies.json,
	//	// (https://core.telegram.org/bots/payments/currencies.json)
	//	// it shows the number of digits past the decimal point
	//	// for each currency (2 for the majority of currencies).
	TotalAmount int `json:"total_amount"`
	// InvoicePayload bot specified invoice payload
	InvoicePayload string `json:"invoice_payload"`
	// ShippingOptionID identifier of the shipping option chosen by the user
	//
	// optional
	ShippingOptionID string `json:"shipping_option_id,omitempty"`
	// OrderInfo order info provided by the user
	//
	// optional
	OrderInfo *OrderInfo `json:"order_info,omitempty"`
}

type OrderInfo struct {
	// Name username
	//
	// optional
	Name string `json:"name,omitempty"`
	// PhoneNumber user's phone number
	//
	// optional
	PhoneNumber string `json:"phone_number,omitempty"`
	// Email user email
	//
	// optional
	Email string `json:"email,omitempty"`
	// ShippingAddress user shipping address
	//
	// optional
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}
