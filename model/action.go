package model

type Action struct {
	Id              uint64
	AccountId       string
	ActionTitle     string
	ActionType      string
	ActionTokens    string
	ActionAmount    string
	Template        string
	ActionNetworkId string
	DappId          int
	DappCategoryId  int
	NetworkId       int
	ToNetworkId     int
	Source          string
}

type ActionDapp struct {
	Id           uint64
	RecordId     uint64
	Count        uint64
	Participants uint64
	ActionType   string
	DappId       int
}

type ActionChain struct {
	Id          uint64
	RecordId    uint64
	Count       uint64
	NetworkId   int
	DappId      int
	ActionTitle string
}
