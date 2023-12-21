package dal

const (
	FindUserRewardMaxIdSql     = `select max(id) from user_reward`
	FindUserRewardByBetweenSql = `select id,account_id,claimed_reward from user_reward where id between $1 and $2 order by id asc`
	UpdateUserRewardRankSql    = `insert into user_reward_rank(account_id,reward,rank,created_at) VALUES($1,$2,$3,$4) ON CONFLICT (account_id) DO UPDATE SET reward=EXCLUDED.reward,rank=EXCLUDED.rank,updated_at=EXCLUDED.updated_at`
)
