package dal

const (
	FindAccountIdSql          = `select id,address from user_info`
	FindAccountIdByAddressSql = FindAccountIdSql + ` where address in(`
)
