package dal

const (
	FindMaxRecordIdFromActionQuestSql = `select max(record_id) from t_action_quest`
	FindAllQuestSql                   = `select id,quest_campaign_id,quest_category_id,start_time,end_time,total_action,status from quest order by id asc`
	FindAllQuestActionSql             = `select id,quest_campaign_id,quest_id,times,category_id,source,dapps,networks,to_networks from quest_action where category='dapp' order by id asc`
	FindUserQuestSql                  = `select id,quest_campaign_id,quest_id,account_id,action_completed,status from user_quest`
	FindUserQuestActionSql            = `select id,quest_campaign_id,quest_id,quest_action_id,account_id,times,status from user_quest_action`
)
