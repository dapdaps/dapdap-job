package dal

const (
	FindQuestActionRecordIdSql        = `select max(record_id) from t_action_quest`
	FindQuestCampaignInfoSql          = `select total_users,total_reward,total_quest_execution from quest_campaign_info limit 1`
	FindQuestCampaignByStatusSql      = `select id from quest_campaign where status=$1 order by id asc`
	FindQuestCampaignByNotStatusSql   = `select id,start_time,end_time,status from quest_campaign where status != $1`
	FindQuestByNotStatusSql           = `select id,start_time,end_time,status from quest where status != $1`
	SaveQuestActionRecordIdSql        = `insert into t_action_quest(record_id) VALUES($1)`
	UpdateQuestActionRecordIdSql      = `update t_action_quest set record_id=$1`
	FindAllQuestByStatusIdSql         = `select id,quest_campaign_id,quest_category_id,total_action,status,reward from quest where quest_campaign_id=$1 and status=$2 order by id asc`
	FindAllQuestActionByCategoryIdSql = `select id,quest_campaign_id,quest_id,times,category_id,source,dapps,networks,to_networks from quest_action where quest_campaign_id=$1 and category=$2 order by id asc`
	FindUserQuestSql                  = `select id,quest_campaign_id,quest_id,account_id,action_completed,status from user_quest`
	FindUserQuestActionSql            = `select id,quest_campaign_id,quest_id,quest_action_id,account_id,times,status from user_quest_action`
	UpdateUserQuestSql                = `insert into user_quest(account_id,quest_id,quest_campaign_id,action_completed,status,updated_at) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT (account_id,quest_id) DO UPDATE SET action_completed=EXCLUDED.action_completed,status=EXCLUDED.status,updated_at=EXCLUDED.updated_at`
	UpdateUserQuestActionSql          = `insert into user_quest_action(account_id,quest_action_id,quest_id,quest_campaign_id,times,status,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (account_id,quest_action_id) DO UPDATE SET times=EXCLUDED.times,status=EXCLUDED.status,updated_at=EXCLUDED.updated_at`
	UpdateQuestCampaignStatusSql      = `update quest_campaign set status=$1 where id=$2`
	UpdateQuestStatusSql              = `update quest set status=$1 where id=$2`
	UpdateUserQuestStatusSql          = `update user_quest set status=$1 where quest_id=$2`
	UpdateQuestCampaignInfoSql        = `update quest_campaign_info set total_users=$1,total_quest_execution=$2,updated_at=$3`
	UpdateTotalRewardSql              = `update quest_campaign_info set total_reward=$1,updated_at=$2`
	FindCountQuestCampaignInfoSql     = `select count(*) from quest_campaign_info`
	InitQuestCampaignInfoSql          = `insert into quest_campaign_info(total_users,total_reward,total_quest_execution) VALUES(0,0,0)`
	//UpdateQuestCampaignSql          = `update quest_campaign set total_users=$1,total_reward=$2,total_quest_execution=$3,updated_at=$4 where id=$5`
)
