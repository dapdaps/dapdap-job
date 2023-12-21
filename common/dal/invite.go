package dal

const (
	FindInvitesSql           = `select id, creator_user_id, used_user_id from invite_code_pool`
	FindTotalInviteRewardSql = `select sum(reward) from invite_code_pool where creator_user_id=$1`
	UpdateInviteSql          = `update invite_code_pool set status=$1,reward=$2,updated_at=$3 where id=$4`
)
