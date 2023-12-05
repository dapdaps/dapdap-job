package dal

const (
	FindQuestActionRecordIdSql              = `select max(record_id) from t_action_quest`
	FindQuestCampaignByStatusSql            = `select id from quest_campaign where status=$1 order by id asc`
	SaveQuestActionRecordIdSql              = `insert into t_action_quest(record_id) VALUES($1)`
	UpdateQuestActionRecordIdSql            = `update t_action_quest set record_id=$1`
	FindAllQuestByStatusIdSql               = `select id,quest_campaign_id,quest_category_id,total_action,status,reward from quest where quest_campaign_id=$1 and status=$2 order by id asc`
	FindAllQuestActionByCategoryIdSql       = `select id,quest_campaign_id,quest_id,times,category_id,source,dapps,networks,to_networks from quest_action where quest_campaign_id=$1 and category=$2 order by id asc`
	FindUserQuestSql                        = `select id,quest_campaign_id,quest_id,account_id,action_completed,status from user_quest`
	FindUserQuestActionSql                  = `select id,quest_campaign_id,quest_id,quest_action_id,account_id,times,status from user_quest_action`
	UpdateUserQuestSql                      = `insert into user_quest(account_id,quest_id,quest_campaign_id,action_completed,status,updated_at) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT (account_id,quest_id) DO UPDATE SET action_completed=EXCLUDED.action_completed,status=EXCLUDED.status,updated_at=EXCLUDED.updated_at`
	UpdateUserQuestActionSql                = `insert into user_quest_action(account_id,quest_action_id,quest_id,quest_campaign_id,times,status,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (account_id,quest_action_id) DO UPDATE SET times=EXCLUDED.times,status=EXCLUDED.status,updated_at=EXCLUDED.updated_at`
	FindUserRewardByIdSql                   = `select reward from user_reward where account_id=$1`
	UpdateUserRewardByIdSql                 = `insert into user_reward(account_id,reward,updated_at) VALUES($1,$2,$3) ON CONFLICT (account_id) DO UPDATE SET reward=EXCLUDED.reward,updated_at=EXCLUDED.updated_at`
	FindQuestCampaignRewardByAccountIdSql   = `select reward from quest_campaign_reward where account_id=$1 and quest_campaign_id=$2`
	UpdateQuestCampaignRewardByAccountIdSql = `insert into quest_campaign_reward(account_id,quest_campaign_id,reward,updated_at) VALUES($1,$2,$3,$4) ON CONFLICT (account_id,quest_campaign_id) DO UPDATE SET reward=EXCLUDED.reward,updated_at=EXCLUDED.updated_at`
)