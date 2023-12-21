package model

type Invite struct {
	Id            int
	CreatorUserId int64
	UsedUserId    int64
	Reward        int64
}
