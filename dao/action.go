package dao

import (
	"dapdap-job/common/dal"
	"dapdap-job/model"
	"database/sql"
	"errors"
	"fmt"
)

func (d *Dao) ActionRows(rows *sql.Rows) (action *model.Action, err error) {
	var (
		accountId       sql.NullString
		amount          sql.NullString
		actionTitle     sql.NullString
		actionType      sql.NullString
		actionTokens    sql.NullString
		template        sql.NullString
		actionNetworkId sql.NullString
		dappId          sql.NullInt64
		networkId       sql.NullInt64
		dappCategoryId  sql.NullInt64
		toNetworkId     sql.NullInt64
		source          sql.NullString
	)
	action = &model.Action{}
	if err = rows.Scan(&action.Id, &accountId, &actionTitle, &actionType, &actionTokens, &amount, &template, &actionNetworkId, &dappId, &networkId, &dappCategoryId, &toNetworkId, &source); err != nil {
		return
	}
	if accountId.Valid {
		action.AccountId = accountId.String
	}
	if amount.Valid {
		action.ActionAmount = amount.String
	}
	if actionTitle.Valid {
		action.ActionTitle = actionTitle.String
	}
	if actionType.Valid {
		action.ActionType = actionType.String
	}
	if actionTokens.Valid {
		action.ActionTokens = actionTokens.String
	}
	if template.Valid {
		action.Template = template.String
	}
	if actionNetworkId.Valid {
		action.ActionNetworkId = actionNetworkId.String
	}
	if dappId.Valid {
		action.DappId = int(dappId.Int64)
	}
	if networkId.Valid {
		action.NetworkId = int(networkId.Int64)
	}
	if dappCategoryId.Valid {
		action.DappCategoryId = int(dappCategoryId.Int64)
	}
	if toNetworkId.Valid {
		action.ToNetworkId = int(toNetworkId.Int64)
	}
	if source.Valid {
		action.Source = source.String
	}
	return
}

func (d *Dao) FindActions(id uint64, limit uint64) (data []*model.Action, err error) {
	var (
		rows *sql.Rows
	)
	rows, err = d.db.Query(dal.FindActionByLimitSql, id, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var action *model.Action
		if action, err = d.ActionRows(rows); err != nil {
			return
		}
		data = append(data, action)
	}
	return
}

func (d *Dao) FindAllActions(id uint64) (data []*model.Action, err error) {
	var (
		maxId uint64
		rows  *sql.Rows
		limit = uint64(5000)
	)
	maxId, err = d.FindMaxId()
	if err != nil {
		return
	}
	for {
		rows, err = d.db.Query(dal.FindActionByBetweenSql, id, id+limit)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return
		}
		err = nil
		for rows.Next() {
			var action *model.Action
			if action, err = d.ActionRows(rows); err != nil {
				return
			}
			data = append(data, action)
		}
		_ = rows.Close()
		if len(data) == 0 {
			return
		}
		if data[len(data)-1].Id >= maxId {
			return
		}
		id = id + limit + 1
	}
}

func (d *Dao) FindMaxId() (id uint64, err error) {
	var (
		maxRecordIdSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindMaxRecordIdSql).Scan(&maxRecordIdSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if maxRecordIdSql.Valid {
		id = uint64(maxRecordIdSql.Int64)
	}
	return
}

func (d *Dao) FindMaxRecordId() (maxRecordIdActionDapp uint64, maxRecordIdActionChain uint64, err error) {
	var (
		maxRecordIdActionDappSql  sql.NullInt64
		maxRecordIdActionChainSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindMaxRecordIdFromActionDappSql).Scan(&maxRecordIdActionDappSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		} else {
			return
		}
	}
	if maxRecordIdActionDappSql.Valid {
		maxRecordIdActionDapp = uint64(maxRecordIdActionDappSql.Int64)
	}

	err = d.db.QueryRow(dal.FindMaxRecordIdFromActionChainSql).Scan(&maxRecordIdActionChainSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		} else {
			return
		}
	}
	if maxRecordIdActionChainSql.Valid {
		maxRecordIdActionChain = uint64(maxRecordIdActionChainSql.Int64)
	}
	return
}

func (d *Dao) FindMaxActionChainId() (id uint64, err error) {
	var (
		maxIdSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindMaxIdFromActionChainSql).Scan(&maxIdSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if maxIdSql.Valid {
		id = uint64(maxIdSql.Int64)
	}
	return
}

func (d *Dao) FindActionsDapp() (data map[string]*model.ActionDapp, err error) {
	var (
		rows *sql.Rows
	)
	data = map[string]*model.ActionDapp{}
	rows, err = d.db.Query(dal.FindActionsDappSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var (
			action = &model.ActionDapp{}
		)
		if err = rows.Scan(&action.Id, &action.RecordId, &action.Count, &action.Participants, &action.Template); err != nil {
			return
		}
		data[action.Template] = action
	}
	return
}

func (d *Dao) FindActionsChain() (data map[string]*model.ActionChain, err error) {
	var (
		maxId uint64
		rows  *sql.Rows
		id    = uint64(1)
		limit = uint64(5000)
	)
	maxId, err = d.FindMaxActionChainId()
	if err != nil {
		return
	}
	data = map[string]*model.ActionChain{}
	for {
		var findMaxId uint64
		rows, err = d.db.Query(dal.FindActionsChainSql, id, id+limit)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = nil
			}
			return
		}
		for rows.Next() {
			var (
				action = &model.ActionChain{}
			)
			if err = rows.Scan(&action.Id, &action.RecordId, &action.Count, &action.ActionTitle, &action.Template, &action.ActionNetworkId); err != nil {
				return
			}
			data[fmt.Sprintf("%s_%s_%s", action.ActionNetworkId, action.Template, action.ActionTitle)] = action
			findMaxId = action.Id
		}
		_ = rows.Close()
		if findMaxId >= maxId {
			return
		}
		id = id + limit + 1
	}
}

func (d *Dao) UpdateActionsDapp(data map[string]*model.ActionDapp) (err error) {
	err = d.WithTrx(func(db *sql.Tx) (err error) {
		for _, actionDapp := range data {
			_, err = db.Exec(dal.UpdateActionDappSql, actionDapp.RecordId, actionDapp.Count, actionDapp.Participants, actionDapp.Template)
			if err != nil {
				return
			}
		}
		return
	})
	return
}

func (d *Dao) UpdateActionsChain(data map[string]*model.ActionChain) (err error) {
	err = d.WithTrx(func(db *sql.Tx) (err error) {
		for _, actionChain := range data {
			_, err = db.Exec(dal.UpdateActionChainSql, actionChain.RecordId, actionChain.Count, actionChain.ActionTitle, actionChain.Template, actionChain.ActionNetworkId)
			if err != nil {
				return
			}
		}
		return
	})
	return
}
