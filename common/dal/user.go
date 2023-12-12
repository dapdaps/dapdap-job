package dal

const (
	FindAccountIdSql          = `select id,address from user_info`
	FindAccountIdByAddressSql = FindAccountIdSql + ` where address in(`
	FindAccountForUpdateSql   = `select id from user_info where id=$1 for update`
	FindAccountIdByTgSql      = FindAccountIdSql + ` where tg_user_id=$1`
)
