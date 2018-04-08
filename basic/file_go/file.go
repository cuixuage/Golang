package main

import (
	"fmt"
	"os"
)

// func Mkdir(name string, perm FileMode) error
// 创建名称为name的目录，权限设置是perm，例如0777
// func MkdirAll(path string, perm FileMode) error
// 根据path创建多级子目录，例如astaxie/test1/test2。
func test_mkdir() {
	os.Mkdir("astaxie", 0777)
	os.MkdirAll("astaxie/test1/test2", 0777)

	// err := os.Remove("astaxie")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// os.RemoveAll("astaxie")       //多级目录使用removeall
}

func test_create() {
	userFile := "astaxie.txt"
	fout, err := os.Create(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	defer fout.Close() //defer关键字   是释放资源
	// defer后面的表达式会被放入一个列表中，在当前方法返回的时候，列表中的表达式就会被执行
	// 栈(stack)的结构，是一个后进先出的栈

	for i := 0; i < 10; i++ {
		fout.WriteString("Just a test!\r\n")
		fout.Write([]byte("Just a test!\r\n"))
		//byte仍然是按照string保存到文件中的？？？？
		//注意 WriteString 和 write 的区别

	}
}

// Go语言里面删除文件和删除文件夹是同一个函数
// func Remove(name string) Error

func main() {
	// test_mkdir()
	test_create()
}
