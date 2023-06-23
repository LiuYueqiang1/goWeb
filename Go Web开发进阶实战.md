# Go操作MySQL

```go
import (
   "database/sql"
   "fmt"
   _ "github.com/go-sql-driver/mysql" // init()
)
```

## 初始化数据库

```go
dsn := "root:961024@tcp(localhost:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
//打开数据库
sql.Open("mysql", dsn)
// 尝试与数据库建立连接
err = db.Ping()
	db.SetMaxOpenConns(100) //最大连接数
	db.SetMaxIdleConns(10)  //最大空闲连接数
```

## MySQL语句

```go
sqlStr := `select id, name, gender ,hobby from user_infos where id=?;`
sqlStr := "select id, name, gender ,hobby from user_infos where id>?"
sqlStr := "insert into user_infos(name,gender,hobby) values(?,?,?)"
sqlStr := "update user_infos set hobby = ?,name=? where id = ?"
sqlStr := "delete from user_infos where id = ?"
```

## 查询

### 单条

```go
err := db.QueryRow(sqlStr, 2).Scan(&u.id, &u.name, &u.hobby, &u.gender)
```

### 多条

```go
sqlStr := "select id, name, gender ,hobby from user_infos where id>?"
rows, err := db.Query(sqlStr, 1)
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.hobby, &u.gender)
		fmt.Printf("id:%d name:%s hobby:%s gender:%s \n", u.id, u.name, u.gender, u.hobby)
	}
```

## 插入、更新、删除

```go
ret, err := db.Exec(sqlStr, "年轻人", "男", "练死劲儿")
n, err = ret.RowsAffected() // 操作影响的行数
```

## 预处理

```
stmt, err := db.Prepare(sqlStr)
defer stmt.Close() //关闭命令
rows, err := stmt.Query(0)
defer rows.Close() //关闭数据
```

## 事务

```
tx, err := db.Begin() // 开启事务
tx.Rollback() // 操作失败则回滚事务
如果两次影响的行数都为1，则提交
tx.Commit() // 提交事务
```