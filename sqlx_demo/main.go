package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func initDB() (err error) {
	dsn := "root:961024@tcp(localhost:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	return
}

type user struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func main() {
	if err := initDB(); err != nil {
		fmt.Printf("init DB failed, err:%v\n", err)
		return
	}
	fmt.Println("init DB success...")
	//queryRowDemo()
	//queryMultiRowDemo()
	//deleteRowDemo()
	//queryMultiRowDemo()
	//insertUserDemo()

	// sqlx.In批量插入
	//u1 := user{Name: "xx", Age: 18}
	//u2 := user{Name: "xxx", Age: 28}
	//u3 := user{Name: "xxxx", Age: 38}
	//users := []interface{}{u1, u2, u3}
	//BatchInsertUsers2(users)

	// NamedExec 批量插入
	//users3 := []*user{&u1, &u2, &u3}
	//err := BatchInsertUsers3(users3)
	//if err != nil {
	//	fmt.Printf("BatchInsertUsers3 failed, err:%v\n", err)
	//}

	transactionDemo2()
}

// 查询单行
func queryRowDemo() {
	sqlStr := `select id, name, age from users where id=?;`
	var u user
	err := db.Get(&u, sqlStr, 14)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Printf("id :%d,name:%s age:%d\n", u.ID, u.Name, u.Age)
}

// 查询多行
// 查询多条数据示例
// 不需要扫描，建一个切片
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from users where id > ?"
	var users []user
	err := db.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", users)
}

// exec 增删改
// 插入数据
func insertRowDemo() {
	sqlStr := "insert into users(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "沙河小王子", 19)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update users set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 17)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from users where id = ?"
	ret, err := db.Exec(sqlStr, 17)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

// DB.NamedExec方法用来绑定SQL语句与结构体或map中的同名字段。
func insertUserDemo() (err error) {
	sqlStr := "INSERT INTO users (name,age) VALUES (:name,:age)"
	_, err = db.NamedExec(sqlStr,
		map[string]interface{}{
			"name": "七米",
			"age":  28,
		})
	return
}

// 使用sqlx.In实现批量插入
// 前提是需要我们的结构体实现driver.Valuer接口：
func (u user) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func BatchInsertUsers2(users []interface{}) error {
	query, args, _ := sqlx.In(
		"INSERT INTO users (name, age) VALUES (?), (?), (?)",
		users..., // 如果arg实现了 driver.Valuer, sqlx.In 会通过调用 Value()来展开它
	)
	fmt.Println(query) // 查看生成的querystring
	fmt.Println(args)  // 查看生成的args
	_, err := db.Exec(query, args...)
	return err
}

// 使用NamedExec实现批量插入的代码如下
func BatchInsertUsers3(users []*user) error {
	_, err := db.NamedExec("INSERT INTO users (name, age) VALUES (:name, :age)", users)
	return err
}

// 事务操作
func transactionDemo2() (err error) {
	tx, err := db.Beginx() //开启事务
	if err != nil {
		fmt.Printf("begin trans failed,err:%v\n", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
			fmt.Println("commit")
		}
	}()

	sqlStr1 := "Update users set age=69 where id = ?"

	rs, err := tx.Exec(sqlStr1, 15)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	sqlStr2 := "Update users set age=31 where id = ?"

	rs, err = tx.Exec(sqlStr2, 19)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr2 failed")
	}
	return err
}
