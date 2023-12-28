package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //不直接使用，执行init()
)

// 定义一个全局对象db
var db *sql.DB //内部维护一个连接池，并发安全！ 如果需要其它模块访问，需要修改为大写！
func initDB() (err error) {
	dsn := "root:123456@tcp(localhost:3306)/mysql_demo"
	//去初始化全局的db变量而不是新声明一个
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		//panic(err)
		return err
	}
	//fmt.Println("success??")

	err = db.Ping() //尝试与数据库建立连接（校验dsn的正确性）
	if err != nil {
		//fmt.Println("connect to db failed")
		return err
	}

	//fmt.Println("success!!")
	return nil
}

func main() {
	if err := initDB(); err != nil {
		fmt.Printf("connect to db failed, err%v\n", err)
	}
	defer db.Close() //细节问题，defer不能在err之前 -> 需要先保证db不为nil ！

}
