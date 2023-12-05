package model

type QuestCampaign struct {
	Id                  int
	TotalUsers          int
	TotalReward         int
	TotalQuestExecution int
}

type Quest struct {
	Id              int
	QuestCampaignId int
	QuestCategoryId int
	TotalAction     int
	Status          string
	Reward          int
}

type QuestAction struct {
	Id              int
	QuestCampaignId int
	QuestId         int
	Times           int
	CategoryId      int
	Source          string
	Dapps           string
	Networks        string
	ToNetworks      string
	DappsMap        map[int]int
	NetworksMap     map[int]int
	ToNetworksMap   map[int]int
}

type UserQuest struct {
	Id              int
	QuestId         int
	QuestCampaignId int
	AccountId       int
	ActionCompleted int
	Status          string
}

type UserQuestAction struct {
	Id              int
	QuestActionId   int
	QuestId         int
	QuestCampaignId int
	AccountId       int
	Times           int
	Status          string
}

type QuestCampaignReward struct {
	QuestCampaignId int
	AccountId       int
	Reward          int
}
