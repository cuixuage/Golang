package main

import "unsafe"

//常量作为枚举类型
const (
    a = "abc"
    b = len(a)      //b=3
    c = unsafe.Sizeof(a)   //？？？为什么是16
    /*
    字符串类型在 go 里是个结构, 包含指向底层数组的指针和长度,这两部分每部分都是 8 个字节，所以字符串类型大小为 16 个字节
    */
)

func main(){
    println(a, b, c)
}
