package dal

const (
	FindAccountIdSql          = `select id,address from user_info`
	FindAccountIdByAddressSql = FindAccountIdSql + ` where address in(`
	FindAccountForUpdateSql   = `select id from user_info where id=$1 for update`
	FindAccountExtSql         = `select account_id,twitter_user_id,twitter_access_token_type,twitter_access_token_expires,twitter_access_token,telegram_user_id,discord_user_id,updated_at from user_info_ext`
	FindAccountIdByExtSql     = `select account_id from user_info_ext`
	FindAccountIdByTgSql      = FindAccountIdByExtSql + ` where tg_user_id=$1`
	FindAccountIdByDiscordSql = FindAccountIdByExtSql + ` where discord_user_id=$1`
)
