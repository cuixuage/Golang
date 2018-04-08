// func Contains(s, substr string) bool
// func Join(a []string, sep string) string  字符串链接，把slice a通过sep链接起来  strings.Join(s, ", ")
// func Replace(s, old, new string, n int) string 在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换
// func Split(s, sep string) []string  把s字符串按照sep分割，返回slice

package main

import (
	"fmt"
	"strconv"
)

func test_strconv() {
	str := make([]byte, 0, 100) //s :=make([]int,len,cap) 初始化
	str = strconv.AppendInt(str, 4567, 10)
	str = strconv.AppendBool(str, false)
	str = strconv.AppendQuote(str, "abcdefg")
	str = strconv.AppendQuoteRune(str, '单')
	fmt.Println(string(str))
}

// byte数据类型
// 类似 uint8
