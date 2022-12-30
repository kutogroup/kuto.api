package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/kutogroup/kuto.api/models"
	"github.com/kutogroup/kuto.api/utils"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
)

var timeZone = url.QueryEscape("Asia/Shanghai")

//KutoDB 数据库结构体
type KutoDB struct {
	dbmap *gorp.DbMap
}

//KutoTx 数据库事务体
type KutoTx struct {
	dbtx *gorp.Transaction
}

//NewDatabase 新建数据库对象
func NewDatabase(table, addr, user, pwd string, loggable bool) *KutoDB {
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=%s", user, pwd, addr, table, timeZone)
	fmt.Println("database conn string is", sqlStr)
	db, err := sql.Open("mysql", sqlStr)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	database := &KutoDB{
		dbmap: &gorp.DbMap{
			Db: db,
			Dialect: gorp.MySQLDialect{
				Engine:   "InnoDB",
				Encoding: "UTF8"},
		},
	}

	if loggable {
		database.dbmap.TraceOn("[gorp]", log.New(os.Stdout, "kk:", log.Lmicroseconds))
	}
	for _, v := range models.ModelTables {
		database.dbmap.AddTableWithName(v, utils.StructGetLineName(v)).SetKeys(true, "ID")
	}

	return database
}

func NewDatabaseCustom(table, addr, user, pwd string, loggable bool, structTables []interface{}) *KutoDB {
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=%s", user, pwd, addr, table, timeZone)
	fmt.Println("database conn string is", sqlStr)
	db, err := sql.Open("mysql", sqlStr)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	database := &KutoDB{
		dbmap: &gorp.DbMap{
			Db: db,
			Dialect: gorp.MySQLDialect{
				Engine:   "InnoDB",
				Encoding: "UTF8"},
		},
	}

	if loggable {
		database.dbmap.TraceOn("[gorp]", log.New(os.Stdout, "kk:", log.Lmicroseconds))
	}

	for _, v := range structTables {
		database.dbmap.AddTableWithName(v, utils.StructGetLineName(v)).SetKeys(true, "ID")
	}

	return database

}

//SelectByID 根据ID查询
func (db *KutoDB) SelectByID(holder interface{}, id int64) error {
	v, _ := utils.ReflectGetVT(holder)
	res := ""
	for n := 0; n < v.NumField(); n++ {
		res = res + v.Field(n).Tag.Get("db") + ","
	}

	res = strings.TrimRight(res, ",")
	return db.dbmap.SelectOne(holder,
		fmt.Sprintf("SELECT "+res+" FROM %s WHERE deleted=0 AND (%s)", utils.StructGetLineName(holder), "id=?"),
		id)
}

//Select 根据条件查询
func (db *KutoDB) Select(holder interface{}, where string, args ...interface{}) error {
	v, t := utils.ReflectGetVT(holder)

	if t.Kind() != reflect.Slice {
		return errors.New("First params must be slice")
	}

	s := v.Elem()
	res := ""
	for n := 0; n < s.NumField(); n++ {
		res = res + s.Field(n).Tag.Get("db") + ","
	}

	res = strings.TrimRight(res, ",")
	sql := "SELECT " + res + " FROM " + utils.ConvertCamel2Line(s.Name()) + " WHERE "
	if len(where) > 0 {
		if where[0] >= 'A' && where[0] <= 'Z' {
			//首字母大写表示没有where条件了
			sql = sql + "deleted=0 " + where
		} else {
			if !strings.Contains(where, "deleted") {
				//如果不包含deleted，则默认选择未删除的
				sql = sql + "deleted=0 AND " + where
			} else {
				//如果包含deleted，则直接追加
				sql = sql + where
			}
		}
	}

	_, err := db.dbmap.Select(holder, sql, args...)
	return err
}

//Insert 插入数据库
func (db *KutoDB) Insert(holder ...interface{}) error {
	return db.dbmap.Insert(holder...)
}

//Update 更新数据库
func (db *KutoDB) Update(holder interface{}) (int64, error) {
	return db.dbmap.Update(holder)
}

//Delete 删除数据(这里仅仅只修改数据库deleted字段)
func (db *KutoDB) Delete(holder interface{}) (int64, error) {
	utils.ReflectSetValue(holder, "Deleted", int64(1))
	return db.Update(holder)
}

//Exec 执行sql语句
func (db *KutoDB) Exec(sql string, args ...interface{}) error {
	_, err := db.dbmap.Exec(sql, args...)
	return err
}

//Begin 开启事务
func (db *KutoDB) Begin() (*KutoTx, error) {
	tx, err := db.dbmap.Begin()
	return &KutoTx{
		dbtx: tx,
	}, err
}

//Commit 提交事务
func (tx *KutoTx) Commit() error {
	return tx.dbtx.Commit()
}

//Rollback 回滚事务
func (tx *KutoTx) Rollback() error {
	return tx.dbtx.Rollback()
}

//SelectByID 根据ID查询
func (tx *KutoTx) SelectByID(holder interface{}, id int64) error {
	v, _ := utils.ReflectGetVT(holder)
	res := ""
	for n := 0; n < v.NumField(); n++ {
		res = res + v.Field(n).Tag.Get("db") + ","
	}

	res = strings.TrimRight(res, ",")

	return tx.dbtx.SelectOne(holder,
		fmt.Sprintf("SELECT "+res+" FROM %s WHERE deleted=0 AND (%s)", utils.StructGetLineName(holder), "id=?"),
		id)
}

//Select 根据条件查询
func (tx *KutoTx) Select(holder interface{}, where string, args ...interface{}) error {
	v, t := utils.ReflectGetVT(holder)

	if t.Kind() != reflect.Slice {
		return errors.New("First params must be slice")
	}

	s := v.Elem()
	res := ""
	for n := 0; n < s.NumField(); n++ {
		res = res + s.Field(n).Tag.Get("db") + ","
	}

	res = strings.TrimRight(res, ",")

	sql := "SELECT " + res + " FROM " + utils.ConvertCamel2Line(s.Name()) + " WHERE "
	if len(where) > 0 {
		if where[0] >= 'A' && where[0] <= 'Z' {
			//首字母大写表示没有where条件了
			sql = sql + "deleted=0 " + where
		} else {
			if !strings.Contains(where, "deleted") {
				//如果不包含deleted，则默认选择未删除的
				sql = sql + "deleted=0 AND " + where
			} else {
				//如果包含deleted，则直接追加
				sql = sql + where
			}
		}
	}

	_, err := tx.dbtx.Select(holder, sql, args...)
	return err
}

//Insert 插入数据库
func (tx *KutoTx) Insert(holder ...interface{}) error {
	return tx.dbtx.Insert(holder...)
}

//Update 更新数据库
func (tx *KutoTx) Update(holder interface{}) (int64, error) {
	return tx.dbtx.Update(holder)
}

//Delete 删除数据(这里仅仅只修改数据库deleted字段)
func (tx *KutoTx) Delete(holder interface{}) (int64, error) {
	utils.ReflectSetValue(holder, "Deleted", 1)
	return tx.Update(holder)
}

//Exec 执行sql语句
func (tx *KutoTx) Exec(sql string, args ...interface{}) error {
	_, err := tx.dbtx.Exec(sql, args...)
	return err
}
