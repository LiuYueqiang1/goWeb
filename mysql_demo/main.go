package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // init()
)

// 生命全局变量
var db *sql.DB

func initMySQL() (err error) {
	dsn := "root:961024@tcp(localhost:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// 做完错误之后检查，确保 db 不为 nil
	// 尝试与数据库建立连接
	err = db.Ping()
	if err != nil {
		fmt.Printf("connect to db failed ,err %v\n", err)
		return
	}
	// 数组需要业务具体情况来确定
	db.SetMaxOpenConns(100) //最大连接数
	db.SetMaxIdleConns(10)  //最大空闲连接数
	return
}

type user struct {
	id     int
	name   string
	gender string
	hobby  string
}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := `select id, name, gender ,hobby from user_infos where id=?;`
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	//row := db.QueryRow(sqlStr, 1) // 并没有对其做scan操作
	err := db.QueryRow(sqlStr, 2).Scan(&u.id, &u.name, &u.hobby, &u.gender)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s hobby:%s gender:%s \n", u.id, u.name, u.hobby, u.gender)
}

func main() {
	if err := initMySQL(); err != nil {
		fmt.Printf("connect to db failed, err:%v\n", err)
	}
	// Close() 用来释放掉数据库连接相关的资源
	defer db.Close() // 注意这行代码要写在上面err判断的下面
	fmt.Println("connect to db success")
	queryRowDemo()
}
