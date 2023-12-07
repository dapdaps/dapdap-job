package dal

const (
	FindActionCategorySql = `select id,name from category`
	FindNetworksSql       = `select id,chain_id from network`
)
