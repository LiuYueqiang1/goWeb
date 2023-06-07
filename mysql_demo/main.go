package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // init()
)

// 声明全局变量
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

// 查询多条数据
func queryMultiRowDemo() {
	sqlStr := "select id, name, gender ,hobby from user_infos where id>?"
	rows, err := db.Query(sqlStr, 1)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	//循环读取结果中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.hobby, &u.gender)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s hobby:%s gender:%s \n", u.id, u.name, u.gender, u.hobby)
	}
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user_infos(name,gender,hobby) values(?,?,?)"

	ret, err := db.Exec(sqlStr, "年轻人", "男", "练死劲儿")
	if err != nil {
		fmt.Printf("insert failed,err:%v\n", err)
		return
	}
	var theID int64
	theID, err = ret.LastInsertId() // 新插入数据的ID 自动增量
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success,the id is %d.\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user_infos set hobby = ?,name=? where id = ?"
	ret, err := db.Exec(sqlStr, "颈椎练坏了", "第二个年轻人", "5")
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	var n int64
	n, err = ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user_infos where id = ?"
	ret, err := db.Exec(sqlStr, 5)
	if err != nil {
		fmt.Printf("delete failed,err:%v\n", err)
		return
	}
	var n int64
	n, err = ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

// 预处理 1、优化服务器性能，降低服务器成本 2、解决了sql注入问题
// 预处理查询
// 1、先传命令 2、再传数据
func prepareQueryDemo() {
	sqlStr := "select id ,name,gender,hobby from user_infos where id > ?"
	// 预处理命令
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close() //关闭命令
	// 预处理数据
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	defer rows.Close() //关闭数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.gender, &u.hobby)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s hobby:%s gender:%s \n", u.id, u.name, u.hobby, u.gender)
	}
}

// 预处理插入
func prepareInsertDemo() {
	sqlStr := "insert into user_infos(name,gender,hobby) values(?,?,?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec("掌门人", "男", "五连鞭")
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	_, err = stmt.Exec("真爱粉", "男", "油饼")
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	fmt.Println("insert success.")
}

// sql注入示例
func sqlInjectDemo(name string) {
	sqlStr := fmt.Sprintf("select id, name, gender ,hobby from user_infos where name='%s'", name)
	fmt.Printf("SQL:%s\n", sqlStr)
	var u user
	err := db.QueryRow(sqlStr).Scan(&u.id, &u.name, &u.gender, &u.hobby)
	if err != nil {
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	fmt.Printf("user:%#v\n", u)
}

// 事务操作示例
func transactionDemo() {
	tx, err := db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 操作失败则回滚事务
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "Update user_infos set gender=? where id=?"
	ret1, err := tx.Exec(sqlStr1, "女", 2)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	//看影响的行数
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "Update user_infos set name=? where id=?"
	ret2, err := tx.Exec(sqlStr2, "马保国", 3)
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

	// 当affRow1 == 1 && affRow2 == 1
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
func main() {
	if err := initMySQL(); err != nil {
		fmt.Printf("connect to db failed, err:%v\n", err)

	}
	// Close() 用来释放掉数据库连接相关的资源
	defer db.Close() // 注意这行代码要写在上面err判断的下面
	fmt.Println("connect to db success")
	//queryRowDemo()

	//queryMultiRowDemo()
	//insertRowDemo()
	//prepareInsertDemo()
	//prepareQueryDemo()

	//sql注入示例
	//sqlInjectDemo("xxx ' or 2=2#")
	//事务
	transactionDemo()
	// 查询
	prepareQueryDemo()
}
