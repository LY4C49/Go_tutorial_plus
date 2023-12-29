package main

import (
	"database/sql"
	"fmt"

	//第三方库只要实现database/sql中driver.go 规定的方法即可！
	_ "github.com/go-sql-driver/mysql" //不直接使用，执行init()。调用database/sql 的regeister方法，注册到全局的map中
)

// 定义一个全局对象db
var db *sql.DB //内部维护一个连接池，并发安全！ 如果需要其它模块访问，需要修改为大写！
func initDB() (err error) {
	dsn := "root:123456@tcp(localhost:3306)/sql_test"
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

type user struct {
	id   int
	age  int
	name string
}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id = ?"
	var u user
	//非常重要：确保QueryRow之后调用Scan方法，否则数据库连接不会被释放! 否则资源耗尽以后，就一直卡在这！！！
	//QueryRow返回 Row 对象
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err: %v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
}

// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接。防止for循环中间出错，从而没有释放连接！
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

func transactionDemo() {
	tx, err := db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "Update user set age=30 where id=?"
	ret1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "Update user set age=40 where id=?"
	ret2, err := tx.Exec(sqlStr2, 3)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql2 failed, err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}

	fmt.Println("exec trans success!")
}

func transactionDemo1() {
	tx, err := db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "Update user set age=25 where id=?"
	ret1, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if affRow1 == 1 {
		//tx.Commit()
	} else {
		tx.Rollback()
	}
	fmt.Println(affRow1)
}

func main() {
	if err := initDB(); err != nil {
		fmt.Printf("connect to db failed, err%v\n", err)
	}
	defer db.Close() //细节问题，defer不能在err之前 -> 需要先保证db不为nil ！
	queryRowDemo()
	queryMultiRowDemo()

	//transactionDemo()
	transactionDemo1()
}
