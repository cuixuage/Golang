package main

//似乎默认初始化为0
var x, y int


var (  // 这种因式分解关键字的写法一般用于声明全局变量
    a int
    b bool
)

var c, d int = 1, 2
var e, f = 123, "hello"

//这种不带声明格式的只能在函数体中出现
//g, h := 123, "hello"

func main(){
    //函数体内部使用   :=
    g, h := 123, "hello"

    println(x, y, a, b, c, d, e, f, g, h)
    println(&x," ",&y)
}
