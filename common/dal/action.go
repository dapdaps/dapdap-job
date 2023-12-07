package dal

const (
	findActionSql                     = `select id,account_id,action_title,action_type,action_tokens,action_amount,template,action_network_id,dapp_id,chain_id,to_chain_id,source from t_action_record`
	FindActionByLimitSql              = findActionSql + ` where id >= $1 order by id asc limit $2`
	FindActionByBetweenSql            = findActionSql + ` where id between $1 and $2 order by id asc`
	FindMaxRecordIdSql                = `select max(id) from t_action_record`
	FindMaxRecordIdFromActionDappSql  = `select max(record_id) from t_action_dapp`
	FindMaxRecordIdFromActionChainSql = `select max(record_id) from t_action_chain`
	FindMaxIdFromActionChainSql       = `select max(id) from t_action_chain`
	FindActionsDappSql                = ` select id,record_id,count,participants,dapp_id from t_action_dapp`
	UpdateActionDappSql               = ` insert into t_action_dapp(record_id,count,participants,dapp_id) VALUES($1,$2,$3,$4) ON CONFLICT (dapp_id) DO UPDATE SET record_id=EXCLUDED.record_id, count=EXCLUDED.count, participants=EXCLUDED.participants`
	FindActionsChainSql               = ` select id,record_id,count,action_title,dapp_id,network_id from t_action_chain where id between $1 and $2 order by id asc`
	UpdateActionChainSql              = ` insert into t_action_chain(record_id,count,action_title,dapp_id,network_id) VALUES($1,$2,$3,$4,$5) ON CONFLICT (dapp_id,network_id,action_title) DO UPDATE SET record_id=EXCLUDED.record_id,count=EXCLUDED.count`
)
