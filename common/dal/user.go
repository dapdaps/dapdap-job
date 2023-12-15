package dal

const (
	FindAccountIdSql          = `select id,address from user_info`
	FindAccountIdByAddressSql = FindAccountIdSql + ` where address in(`
	FindAccountForUpdateSql   = `select id from user_info where id=$1 for update`
	FindAccountIdByTgSql      = FindAccountIdSql + ` where tg_user_id=$1`
	FindAccountExtSql         = `select account_id,twitter_user_id,twitter_access_token_type,twitter_access_token_expires,twitter_access_token,twitter_refresh_token,telegram_user_id from user_info`
)
