package mysql

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func Test(t *testing.T) {
	log.Println("helloworld")
	db, err := sql.Open("mysql", "root@/192.168.0.189")
	if err != nil {
		log.Println(err.Error())
	}
	db.Close()
}
