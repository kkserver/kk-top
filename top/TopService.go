package top

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"strings"
)

type TopService struct {
	app.Service

	Create      *TopCreateTask
	Set         *TopSetTask
	Get         *TopTask
	Remove      *TopRemoveTask
	Query       *TopQueryTask
	Exchange    *TopExchangeTask
	BatchCreate *TopBatchCreateTask
}

func (S *TopService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *TopService) HandleTopCreateTask(a ITopApp, task *TopCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Top{}

	v.Alias = task.Alias
	v.Oid = NewOid()
	v.Tags = task.Tags
	v.ItemId = task.ItemId
	v.ItemType = task.ItemType
	v.CityId = task.CityId
	v.CityPath = task.CityPath
	v.ClassifyId = task.ClassifyId
	v.ClassifyPath = task.ClassifyPath
	v.GroupId = task.GroupId

	if task.Options != nil {
		b, err := json.Encode(task.Options)
		if err != nil {
			task.Result.Errno = ERROR_TOP
			task.Result.Errmsg = err.Error()
			return nil
		}
		v.Options = string(b)
	}

	_, err = kk.DBInsert(db, a.GetTopTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Top = &v

	return nil
}

func (S *TopService) HandleTopBatchCreateTask(a ITopApp, task *TopBatchCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	tx, err := db.Begin()

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	err = func() error {

		itemIds := []int64{}

		if task.ItemIds != "" {
			for _, itemId := range strings.Split(task.ItemIds, ",") {
				id := dynamic.IntValue(itemId, 0)
				if id != 0 {
					itemIds = append(itemIds, id)
				}
			}
		}

		v := Top{}

		v.Alias = task.Alias
		v.Tags = task.Tags
		v.ItemType = task.ItemType
		v.CityId = task.CityId
		v.CityPath = task.CityPath
		v.ClassifyId = task.ClassifyId
		v.ClassifyPath = task.ClassifyPath
		v.GroupId = task.GroupId
		v.Oid = NewOid()

		if task.Options != nil {
			b, err := json.Encode(task.Options)
			if err != nil {
				task.Result.Errno = ERROR_TOP
				task.Result.Errmsg = err.Error()
				return nil
			}
			v.Options = string(b)
		}

		for _, itemId := range itemIds {
			v.Oid = v.Oid + 1
			v.ItemId = itemId

			_, err = kk.DBInsert(tx, a.GetTopTable(), a.GetPrefix(), &v)

			if err != nil {
				return err
			}
		}

		return nil

	}()

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {
		tx.Rollback()
		e, ok := err.(*app.Error)
		if ok {
			task.Result.Errno = e.Errno
			task.Result.Errmsg = e.Errmsg
			return nil
		}
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = e.Error()
		return nil
	}

	return nil
}

func (S *TopService) HandleTopSetTask(a ITopApp, task *TopSetTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_TOP_NOT_FOUND_ID
		task.Result.Errmsg = "Not found top id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Top{}

	rows, err := kk.DBQuery(db, a.GetTopTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_TOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		keys := map[string]bool{}

		if task.ItemId != nil {
			v.ItemId = dynamic.IntValue(task.ItemId, v.ItemId)
			keys["itemid"] = true
		}

		if task.ItemType != nil {
			v.ItemType = dynamic.StringValue(task.ItemType, v.ItemType)
			keys["itemtype"] = true
		}

		if task.Tags != nil {
			v.Tags = dynamic.StringValue(task.Tags, v.Tags)
			keys["tags"] = true
		}

		if task.CityId != nil {
			v.CityId = dynamic.IntValue(task.CityId, v.CityId)
			keys["cityid"] = true
		}

		if task.CityPath != nil {
			v.CityPath = dynamic.StringValue(task.CityPath, v.CityPath)
			keys["citypath"] = true
		}

		if task.ClassifyId != nil {
			v.ClassifyId = dynamic.IntValue(task.ClassifyId, v.ClassifyId)
			keys["classifyid"] = true
		}

		if task.ClassifyPath != nil {
			v.ClassifyPath = dynamic.StringValue(task.ClassifyPath, v.ClassifyPath)
			keys["classifypath"] = true
		}

		if task.Options != nil {
			var options interface{} = nil
			if v.Options != "" {
				_ = json.Decode([]byte(v.Options), &options)
			}
			if options == nil {
				options = map[interface{}]interface{}{}
			}
			dynamic.Each(task.Options, func(key interface{}, value interface{}) bool {
				dynamic.Set(options, dynamic.StringValue(key, ""), value)
				return true
			})
			b, _ := json.Encode(options)
			v.Options = string(b)
			keys["options"] = true
		}

		_, err = kk.DBUpdateWithKeys(db, a.GetTopTable(), a.GetPrefix(), &v, keys)

		if err != nil {
			task.Result.Errno = ERROR_TOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Top = &v

	} else {
		task.Result.Errno = ERROR_TOP_NOT_FOUND
		task.Result.Errmsg = "Not found top"
		return nil
	}

	return nil
}

func (S *TopService) HandleTopRemoveTask(a ITopApp, task *TopRemoveTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Top{}

	rows, err := kk.DBQuery(db, a.GetTopTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_TOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		_, err = kk.DBDelete(db, a.GetTopTable(), a.GetPrefix(), " WHERE id=?", v.Id)

		if err != nil {
			task.Result.Errno = ERROR_TOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Top = &v

	} else {
		task.Result.Errno = ERROR_TOP_NOT_FOUND
		task.Result.Errmsg = "Not found top"
		return nil
	}

	return nil
}

func (S *TopService) HandleTopTask(a ITopApp, task *TopTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Top{}

	rows, err := kk.DBQuery(db, a.GetTopTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_TOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Top = &v

	} else {
		task.Result.Errno = ERROR_TOP_NOT_FOUND
		task.Result.Errmsg = "Not found top"
		return nil
	}

	return nil
}

func (S *TopService) HandleTopQueryTask(a ITopApp, task *TopQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	var tops = []Top{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE 1")

	if task.Alias != "" {
		sql.WriteString(" AND alias=?")
		args = append(args, task.Alias)
	}

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	}

	if task.ItemId != 0 {
		sql.WriteString(" AND itemid=?")
		args = append(args, task.ItemId)
	}

	if task.ItemType != "" {
		sql.WriteString(" AND itemtype=?")
		args = append(args, task.ItemType)
	}

	if task.CityId != 0 {
		sql.WriteString(" AND (cityid=? OR cityid=0)")
		args = append(args, task.CityId)
	}

	if task.CityPrefix != "" {
		sql.WriteString(" AND (citypath LIKE ? OR cityid=0)")
		args = append(args, task.CityPrefix+"%")
	}

	if task.ClassifyId != 0 {
		sql.WriteString(" AND classifyid=?")
		args = append(args, task.ClassifyId)
	}

	if task.ClassifyPrefix != "" {
		sql.WriteString(" AND classifypath LIKE ?")
		args = append(args, task.ClassifyPrefix+"%")
	}

	if task.GroupId != 0 {
		sql.WriteString(" AND (groupid=? OR groupid=0)")
		args = append(args, task.GroupId)
	}

	if task.Keyword != "" {
		q := "%" + task.Keyword + "%"
		sql.WriteString(" AND (tags LIKE ?)")
		args = append(args, q)
	}

	sql.WriteString(" ORDER BY oid DESC,id DESC")

	var pageIndex = task.PageIndex
	var pageSize = task.PageSize

	if pageIndex < 1 {
		pageIndex = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	if task.Counter {

		var counter = TopQueryCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize
		counter.RowCount, err = kk.DBQueryCount(db, a.GetTopTable(), a.GetPrefix(), sql.String(), args...)

		if counter.RowCount%pageSize == 0 {
			counter.PageCount = counter.RowCount / pageSize
		} else {
			counter.PageCount = counter.RowCount/pageSize + 1
		}

		task.Result.Counter = &counter
	}

	sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	var v = Top{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetTopTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_TOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		tops = append(tops, v)
	}

	task.Result.Tops = tops

	return nil
}

func (S *TopService) HandleTopExchangeTask(a ITopApp, task *TopExchangeTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.FromId == 0 {
		task.Result.Errno = ERROR_TOP_NOT_FOUND_ID
		task.Result.Errmsg = "Not found fromId"
		return nil
	}

	v := Top{}

	rows, err := kk.DBQuery(db, a.GetTopTable(), a.GetPrefix(), " WHERE id=?", task.FromId)

	if err != nil {
		task.Result.Errno = ERROR_TOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_TOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		if task.ToId == 0 {
			v.Oid = NewOid()
			_, err = kk.DBUpdateWithKeys(db, a.GetTopTable(), a.GetPrefix(), &v, map[string]bool{"oid": true})
			if err != nil {
				task.Result.Errno = ERROR_TOP
				task.Result.Errmsg = err.Error()
				return nil
			}
		} else {

			rows, err = kk.DBQuery(db, a.GetTopTable(), a.GetPrefix(), " WHERE id=?", task.ToId)

			if err != nil {
				task.Result.Errno = ERROR_TOP
				task.Result.Errmsg = err.Error()
				return nil
			}

			defer rows.Close()

			from := v

			if rows.Next() {

				err = scanner.Scan(rows)

				if err != nil {
					task.Result.Errno = ERROR_TOP
					task.Result.Errmsg = err.Error()
					return nil
				}

				oid := from.Oid
				from.Oid = v.Oid
				v.Oid = oid

				_, err = kk.DBUpdateWithKeys(db, a.GetTopTable(), a.GetPrefix(), &from, map[string]bool{"oid": true})
				if err != nil {
					task.Result.Errno = ERROR_TOP
					task.Result.Errmsg = err.Error()
					return nil
				}

				_, err = kk.DBUpdateWithKeys(db, a.GetTopTable(), a.GetPrefix(), &v, map[string]bool{"oid": true})
				if err != nil {
					task.Result.Errno = ERROR_TOP
					task.Result.Errmsg = err.Error()
					return nil
				}

			} else {
				task.Result.Errno = ERROR_TOP_NOT_FOUND
				task.Result.Errmsg = "Not found to top"
				return nil
			}

		}

	} else {
		task.Result.Errno = ERROR_TOP_NOT_FOUND
		task.Result.Errmsg = "Not found from top"
		return nil
	}

	return nil
}
