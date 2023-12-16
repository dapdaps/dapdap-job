package model

type AccountExt struct {
	AccountId                 int
	TwitterUserId             string
	TwitterAccessTokenType    string
	TwitterAccessTokenExpires string
	TwitterAccessToken        string
	TwitterRefreshToken       string
	TelegramUserId            int64
	DiscordUserId             string
}
