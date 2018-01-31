package main

import (
	"fmt"
	"os"
	"regexp"
)

// func Match(pattern string, b []byte) (matched bool, error error)
// func MatchReader(pattern string, r io.RuneReader) (matched bool, error error)
// func MatchString(pattern string, s string) (matched bool, error error)

func IsIP(ip string) (b bool) {
	m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip)
	if !m {
		return false
	}
	return true
}

func test1() {
	if len(os.Args) == 1 { //args[1]是第一个参数
		fmt.Println("Usage:")
		os.Exit(1)
	}
	//判断输入的字符是否为数字   测试结果不太正确？？？
	m, _ := regexp.MatchString("^[0-9]+$", os.Args[1])

	if !m {
		fmt.Println("数字")
	} else {
		fmt.Println("非数字")
	}

}

// //去除连续的换行符
//     re, _ = regexp.Compile("\\s{2,}")
//     src = re.ReplaceAllString(src, "\n")
//
//     使用复杂的正则首先是Compile，它会解析正则表达式是否合法，如果正确，那么就会返回一个Regexp，
//     然后就可以利用返回的Regexp在任意的字符串上面执行需要的操作

func test2() {
	a := "I am learning Go language"

	re, _ := regexp.Compile("[a-z]{2,4}")

	//查找符合正则的第一个
	one := re.Find([]byte(a))
	fmt.Println("Find:", one)

	//查找符合正则的所有slice,n小于0表示返回全部符合的字符串，不然就是返回指定的长度
	// all := re.FindAll([]byte(a), -1)
	all := re.FindAllStringSubmatch(a, -1)
	fmt.Println("FindAll", all)

	//查找符合条件的index位置,开始位置和结束位置
	index := re.FindIndex([]byte(a))
	fmt.Println("FindIndex", index)

	//查找符合条件的所有的index位置，n同上
	allindex := re.FindAllIndex([]byte(a), -1)
	fmt.Println("FindAllIndex", allindex)

}

// func main() {
// 	// test1();
// 	test2()
// }

// https://astaxie.gitbooks.io/build-web-application-with-golang/content/zh/07.3.html
