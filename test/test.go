package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	// 可根据需要自行添加字段
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

var CustomSecret = []byte("夏天夏天悄悄过去")

// var mc = new(CustomClaims)
// mc :=new(jwt.CustomClaims)
func main() {
	//r := gin.Default()

	authHeader := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxODM3Mjk2MDI2MzkwNTI4LCJ1c2VybmFtZSI6IuS4g-exsyIsImV4cCI6MTY4OTkwOTY3NywiaXNzIjoiYmx1ZWJlbGwuY29tL2JsdWViZWxsIn0.owvUVA3UAC9qXjWEyTr6_Z4wt5WRVZsEpcVF7_bToPA"
	//parts := strings.SplitN(authHeader, " ", 2)
	////	fmt.Println(parts[1])
	//_, err := ParseToken(parts[0])
	//if err != nil {
	//	fmt.Println("无效的Token")
	//	return
	//}
	ret, err := ParseToken(authHeader)
	if err != nil {
		fmt.Println("无效的Token")
		return
	}
	fmt.Println(ret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	var mc = new(CustomClaims)
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
