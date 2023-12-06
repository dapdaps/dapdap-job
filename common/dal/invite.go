package dal

const (
	FindInvitesSql  = `select i.id, i.creator_user_id, i.used_user_id, u.invite_reward from invite_code_pool as i left join user_reward as u on i.creator_user_id=u.account_id`
	UpdateInviteSql = `update invite_code_pool set status=$1 where id=$2`
)
