package model

import (
	"math/big"
)

type Quest struct {
	Id              int
	QuestCampaignId int
	QuestCategoryId int
	StartTime       big.Int
	EndTime         big.Int
	TotalAction     int
	Status          string
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
	QuestCampaignId int
	QuestId         int
	AccountId       int
	ActionCompleted int
	Status          string
}

type UserQuestAction struct {
	Id              int
	QuestCampaignId int
	QuestId         int
	QuestActionId   int
	AccountId       int
	Times           int
	Status          string
}
