package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
)

//_ "github.com/go-sql-driver/mysql"

func temp() {
	db, err := sql.Open("mysql", "root@/127.0.0.1") //user:password@/dbname
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db.Close()
}

// SaveData 保存数据到数据库
func SaveData(objs ...interface{}) {
	for i := len(objs); i >= 0; i-- {
		s := reflect.ValueOf(objs[i]).Elem()
		for j := 0; j < s.NumField(); j++ {
			f := s.Field(j)
			reflect.TypeOf(f).Name
		}
	}
}
