package dal

const (
	FindQuestActionRecordIdSql              = `select max(record_id) from t_action_quest`
	FindQuestCampaignByStatusSql            = `select id,total_users,total_reward,total_quest_execution from quest_campaign where status=$1 order by id asc`
	FindQuestCampaignByNotStatusSql         = `select id,start_time,end_time,status from quest_campaign where status != $1`
	FindQuestByNotStatusSql                 = `select id,start_time,end_time,status from quest where status != $1`
	SaveQuestActionRecordIdSql              = `insert into t_action_quest(record_id) VALUES($1)`
	UpdateQuestActionRecordIdSql            = `update t_action_quest set record_id=$1`
	FindAllQuestByStatusIdSql               = `select id,quest_campaign_id,quest_category_id,total_action,status,reward from quest where quest_campaign_id=$1 and status=$2 order by id asc`
	FindAllQuestActionByCategoryIdSql       = `select id,quest_campaign_id,quest_id,times,category_id,source,dapps,networks,to_networks from quest_action where quest_campaign_id=$1 and category=$2 order by id asc`
	FindUserQuestSql                        = `select id,quest_campaign_id,quest_id,account_id,action_completed,status from user_quest`
	FindUserQuestActionSql                  = `select id,quest_campaign_id,quest_id,quest_action_id,account_id,times,status from user_quest_action`
	UpdateUserQuestSql                      = `insert into user_quest(account_id,quest_id,quest_campaign_id,action_completed,status,updated_at) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT (account_id,quest_id) DO UPDATE SET action_completed=EXCLUDED.action_completed,status=EXCLUDED.status,updated_at=EXCLUDED.updated_at`
	UpdateUserQuestActionSql                = `insert into user_quest_action(account_id,quest_action_id,quest_id,quest_campaign_id,times,status,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (account_id,quest_action_id) DO UPDATE SET times=EXCLUDED.times,status=EXCLUDED.status,updated_at=EXCLUDED.updated_at`
	FindQuestCampaignRewardByAccountIdSql   = `select reward from quest_campaign_reward where account_id=$1 and quest_campaign_id=$2`
	FindQuestCampaignRewardByCampaignIdSql  = `select id,quest_campaign_id,account_id,reward,rank from quest_campaign_reward where quest_campaign_id=$1`
	UpdateQuestCampaignRewardByAccountIdSql = `insert into quest_campaign_reward(account_id,quest_campaign_id,reward,updated_at) VALUES($1,$2,$3,$4) ON CONFLICT (account_id,quest_campaign_id) DO UPDATE SET reward=EXCLUDED.reward,updated_at=EXCLUDED.updated_at`
	UpdateQuestCampaignSql                  = `update quest_campaign set total_users=$1,total_reward=$2,total_quest_execution=$3,updated_at=$4 where id=$5`
	UpdateQuestCampaignStatusSql            = `update quest_campaign set status=$1 where id=$2`
	UpdateQuestStatusSql                    = `update quest set status=$1 where id=$2`
	UpdateUserQuestStatusSql                = `update user_quest set status=$1 where quest_id=$2`
)
