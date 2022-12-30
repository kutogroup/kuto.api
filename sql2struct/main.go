package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kutogroup/kuto.api/utils"

	_ "github.com/go-sql-driver/mysql"
)

var IGNORE_COLUMN = []string{
	"create_at", "update_at", "deleted",
}

//用户名，密码，服务器地址，表名，目录地址
var user, pwd, host, table, dest string

//数据库实例
var db *sql.DB

//表的列结构体
type tableColumnsAttr struct {
	Field      sql.NullString
	Type       sql.NullString
	Collation  sql.NullString
	Null       sql.NullString
	Key        sql.NullString
	Default    sql.NullString
	Extra      sql.NullString
	Privileges sql.NullString
	Comment    sql.NullString
}

func main() {
	flag.StringVar(&user, "u", "root", "mysql login user")
	flag.StringVar(&pwd, "p", "root", "mysql login pwd")
	flag.StringVar(&host, "h", "localhost", "mysql host")
	flag.StringVar(&table, "t", "kuto", "mysql table")
	flag.StringVar(&dest, "d", ".", "gen folder")
	flag.PrintDefaults()
	flag.Parse()

	fmt.Println("user=", user, "pwd=", pwd, "host=", host, "table=", table, "dest=", dest)
	os.RemoveAll(dest + "/models")
	err := os.MkdirAll(dest+"/models", 0744)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pwd, host, table))

	if err != nil {
		panic(err)
	}

	tables, err := getTable(db)
	for _, table := range tables {
		columns := getColumns(db, table)

		structName := utils.ConvertLine2Camel(table)
		structString := "package models\n\n"
		structString = structString + "import (\n"
		structString = structString + getSpace(4) + "m \"github.com/kutogroup/kuto.api/models\"\n"
		structString = structString + ")\n\n"

		//生成数据库常量
		var maxConstLength = -1

		for _, column := range columns {
			l := len(utils.ConvertLine2Camel(column.Field.String))
			if l > maxConstLength {
				maxConstLength = l
			}
		}

		structString = structString + "const (\n"
		structString = structString + getSpace(4) + "Table" +
			structName +
			getSpace(maxConstLength+1) +
			" = \"" +
			table + "\"\n"
		for _, column := range columns {
			if isIgnoreColumn(column) {
				continue
			}

			fn := utils.ConvertLine2Camel(column.Field.String)
			structString = structString + getSpace(4) + "Column" +
				structName + fn +
				getSpace(maxConstLength-len(fn)) +
				" = \"" +
				column.Field.String + "\"\n"
		}
		structString = structString + ")\n\n"

		//生成结构体
		structString = structString + "type " + structName + " struct {\n"

		var maxFieldLineLength, maxFieldLength, maxTypeLength int
		for _, column := range columns {
			if isIgnoreColumn(column) {
				continue
			}

			field := utils.ConvertLine2Camel(column.Field.String)
			if len(field) > maxFieldLength {
				maxFieldLength = len(field)
			}

			if len(column.Field.String) > maxFieldLineLength {
				maxFieldLineLength = len(column.Field.String)
			}

			if len(getType(column)) > maxTypeLength {
				maxTypeLength = len(getType(column))
			}
		}

		for _, column := range columns {
			if isIgnoreColumn(column) {
				continue
			}

			var field string

			if column.Key.String == "PRI" {
				field = strings.ToUpper(column.Field.String)
			} else {
				field = utils.ConvertLine2Camel(column.Field.String)
			}

			structString = structString + getSpace(4) +
				field + getSpace(maxFieldLength+1-len(field)) +
				getType(column) + getSpace(maxTypeLength+1-len(getType(column))) +
				"`" +
				"db:\"" + column.Field.String + "\"" + getSpace(maxFieldLineLength+1-len(column.Field.String)) +
				"json:\"" + column.Field.String + "\"" + getSpace(maxFieldLineLength+1-len(column.Field.String)) +
				"comment:\"" + column.Comment.String + "\"" +
				"`" +
				"\n"
		}
		structString = structString + "}\n"

		structString = structString + "\n"
		structString = structString + "func init() {\n"
		structString = structString + getSpace(4) + "m.ModelTables = append(m.ModelTables, " + structName + "{})\n"
		structString = structString + "}"

		utils.FileWriteString(dest+"/models/"+table+".go", structString)
	}
}

func isIgnoreColumn(col tableColumnsAttr) bool {
	for n := 0; n < len(IGNORE_COLUMN); n++ {
		if IGNORE_COLUMN[n] == col.Field.String {
			return true
		}
	}

	return false
}

func getTable(db *sql.DB) ([]string, error) {
	rows, err := db.Query("show tables;")
	if err != nil {
		return nil, err
	}

	var tables []string
	defer rows.Close()
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tables = append(tables, tableName)
	}

	return tables, nil
}

func getColumns(db *sql.DB, tableName string) (columns []tableColumnsAttr) {
	rows, err := db.Query("show full columns from " + tableName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tabColAttr tableColumnsAttr
		err := rows.Scan(
			&tabColAttr.Field,
			&tabColAttr.Type,
			&tabColAttr.Collation,
			&tabColAttr.Null,
			&tabColAttr.Key,
			&tabColAttr.Default,
			&tabColAttr.Extra,
			&tabColAttr.Privileges,
			&tabColAttr.Comment,
		)
		if err != nil {
			panic(err)
		}
		columns = append(columns, tabColAttr)
	}

	return
}

func getType(field tableColumnsAttr) string {
	kv := strings.SplitN(field.Type.String, "(", 2)
	switch kv[0] {
	case "int", "tinyint":
		return "int64"
	case "datetime", "timestamp":
		return "m.Time"
	case "float", "decimal", "double":
		return "float64"
	default:
		return "string"
	}
}

func getSpace(n int) string {
	s := ""

	for i := 0; i < n; i++ {
		s = s + " "
	}

	return s
}
