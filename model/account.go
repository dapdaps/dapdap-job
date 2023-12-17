package model

type AccountExt struct {
	AccountId                 int
	TwitterUserId             string
	TwitterAccessTokenType    string
	TwitterAccessTokenExpires string
	TwitterAccessToken        string
	TwitterRefreshToken       string
	TelegramUserId            string
	DiscordUserId             string
	UpdateTime                int64
	TwitterQuestCompleted     bool
}
