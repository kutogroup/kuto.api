package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	"kuto/models"
	"kuto/utils"
	"log"
	"net/url"
	"reflect"
	"time"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
)

var timeZone = url.QueryEscape("Asia/Shanghai")
var defaultExcludedUpdateColumns = []string{
	"id", "create_at", "update_at",
}

//WahaDB 数据库结构体
type WahaDB struct {
	dbmap *gorp.DbMap
}

//WahaTx 数据库事务体
type WahaTx struct {
	dbtx *gorp.Transaction
}

//NewDatabase 新建数据库对象
func NewDatabase(table, addr, user, pwd string) *WahaDB {
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=%s", user, pwd, addr, table, timeZone)
	fmt.Println("database conn string is", sqlStr)
	db, err := sql.Open("mysql", sqlStr)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	database := &WahaDB{
		dbmap: &gorp.DbMap{
			Db: db,
			Dialect: gorp.MySQLDialect{
				Engine:   "InnoDB",
				Encoding: "UTF8"},
		},
	}

	for _, v := range models.ModelTables {
		database.dbmap.AddTableWithName(v, utils.StructGetLineName(v)).SetKeys(true, "ID")
	}

	return database
}

//SelectByID 根据ID查询
func (db *WahaDB) SelectByID(holder interface{}, id int64) error {
	return db.dbmap.SelectOne(holder,
		fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND (%s)", utils.StructGetLineName(holder), "id=?"),
		id)
}

//Select 根据条件查询
func (db *WahaDB) Select(holder interface{}, where string, args ...interface{}) error {
	v, t := utils.ReflectGetVT(holder)

	if t.Kind() != reflect.Slice {
		return errors.New("First params must be slice")
	}

	sql := "SELECT * FROM " + utils.ConvertCamel2Line(v.Elem().Name()) + " WHERE deleted=0"
	if len(where) > 0 {
		if where[0] >= 'A' && where[0] <= 'Z' {
			//首字母大写表示不是where条件了
			sql = sql + " " + where
		} else {
			sql = sql + " AND " + where
		}
	}

	_, err := db.dbmap.Select(holder, sql, args...)
	return err
}

//Insert 插入数据库
func (db *WahaDB) Insert(holder ...interface{}) error {
	for _, v := range holder {
		t := &models.Time{}
		t.Scan(time.Now())
		utils.ReflectSetValue(v, "CreateAt", *t)
		utils.ReflectSetValue(v, "UpdateAt", *t)
	}

	return db.dbmap.Insert(holder...)
}

//Update 更新数据库
func (db *WahaDB) Update(holder interface{}, filterColumns ...string) error {
	_, err := db.dbmap.UpdateColumns(func(m *gorp.ColumnMap) bool {
		for _, c := range defaultExcludedUpdateColumns {
			if c == m.ColumnName {
				return false
			}
		}

		for _, c := range filterColumns {
			if c == m.ColumnName {
				return true
			}
		}

		return false
	}, holder)
	return err
}

//Delete 删除数据(这里仅仅只修改数据库deleted字段)
func (db *WahaDB) Delete(holder interface{}) error {
	utils.ReflectSetValue(holder, "Deleted", int64(1))
	return db.Update(holder, "deleted")
}

//Exec 执行sql语句
func (db *WahaDB) Exec(sql string, args ...interface{}) error {
	_, err := db.dbmap.Exec(sql, args)
	return err
}

//Begin 开启事务
func (db *WahaDB) Begin() *WahaTx {
	tx, err := db.dbmap.Begin()
	fmt.Println(err)
	return &WahaTx{
		dbtx: tx,
	}
}

//Commit 提交事务
func (tx *WahaTx) Commit() error {
	return tx.dbtx.Commit()
}

//Rollback 回滚事务
func (tx *WahaTx) Rollback() error {
	return tx.dbtx.Rollback()
}

//SelectByID 根据ID查询
func (tx *WahaTx) SelectByID(holder interface{}, id int) error {
	return tx.dbtx.SelectOne(holder,
		fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND (%s)", utils.StructGetLineName(holder), "id=?"),
		id)
}

//Select 根据条件查询
func (tx *WahaTx) Select(holder interface{}, where string, args ...interface{}) error {
	v, t := utils.ReflectGetVT(holder)

	if t.Kind() != reflect.Slice {
		return errors.New("First params must be slice")
	}

	sql := "SELECT * FROM " + utils.ConvertCamel2Line(v.Elem().Name()) + " WHERE deleted=0"
	if len(where) > 0 {
		if where[0] >= 'A' && where[0] <= 'Z' {
			//首字母大写表示不是where条件了
			sql = sql + " " + where
		} else {
			sql = sql + " AND " + where
		}
	}

	_, err := tx.dbtx.Select(holder, sql, args...)
	return err
}

//Insert 插入数据库
func (tx *WahaTx) Insert(holder ...interface{}) error {
	for _, v := range holder {
		t := &models.Time{}
		t.Scan(time.Now())
		utils.ReflectSetValue(v, "CreateAt", *t)
		utils.ReflectSetValue(v, "UpdateAt", *t)
	}

	return tx.dbtx.Insert(holder...)
}

//Update 更新数据库
func (tx *WahaTx) Update(holder interface{}, filterColumns ...string) error {
	_, err := tx.dbtx.UpdateColumns(func(m *gorp.ColumnMap) bool {
		for _, c := range defaultExcludedUpdateColumns {
			if c == m.ColumnName {
				return false
			}
		}

		for _, c := range filterColumns {
			if c == m.ColumnName {
				return true
			}
		}

		return false
	}, holder)
	return err
}

//Delete 删除数据(这里仅仅只修改数据库deleted字段)
func (tx *WahaTx) Delete(holder interface{}) error {
	utils.ReflectSetValue(holder, "Deleted", 1)
	return tx.Update(holder, "deleted")
}

//Exec 执行sql语句
func (tx *WahaTx) Exec(sql string, args ...interface{}) error {
	_, err := tx.dbtx.Exec(sql, args)
	return err
}
