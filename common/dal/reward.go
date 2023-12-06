package dal

const (
	FindUserRewardMaxIdSql           = `select max(id) from user_reward`
	FindUserRewardByIdSql            = `select reward,invite_reward from user_reward where account_id=$1`
	FindUserRewardByBetweenSql       = `select id,account_id,reward from user_reward where id between $1 and $2 order by id asc`
	UpdateUserRewardByIdSql          = `insert into user_reward(account_id,reward,updated_at) VALUES($1,$2,$3) ON CONFLICT (account_id) DO UPDATE SET reward=EXCLUDED.reward,updated_at=EXCLUDED.updated_at`
	UpdateUserInviteRewardByIdSql    = `insert into user_reward(account_id,reward,invite_reward,updated_at) VALUES($1,$2,$3,$4) ON CONFLICT (account_id) DO UPDATE SET reward=EXCLUDED.reward,invite_reward=EXCLUDED.invite_reward,updated_at=EXCLUDED.updated_at`
	FindUserRewardByCategorySql      = `select id,account_id,reward from user_reward where category=$1`
	UpdateUserQuestCampaignRewardSql = `insert into user_reward(account_id,category,reward,updated_at) VALUES($1,$2,$3,$4) ON CONFLICT (account_id,category) DO UPDATE SET reward=EXCLUDED.reward,updated_at=EXCLUDED.updated_at`
	UpdateUserRewardRankSql          = `insert into user_reward_rank(account_id,reward,rank,created_at) VALUES($1,$2,$3,$4) ON CONFLICT (account_id) DO UPDATE SET reward=EXCLUDED.reward,rank=EXCLUDED.rank,updated_at=EXCLUDED.updated_at`
)
