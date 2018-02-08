package main

import (
    "fmt"
)

type Phone interface {
    call()
}

type NokiaPhone struct {
}

//成员函数的定义  func（实例 类） 内部还可以使用类内的变量 那么实例相当于this指针?
func (nokiaPhone NokiaPhone) call() {
    fmt.Println("I am Nokia, I can call you!")
}

type IPhone struct {
}

func (iPhone IPhone) call() {
    fmt.Println("I am iPhone, I can call you!")
}

func main() {
    var phone Phone

    phone = new(NokiaPhone)
    phone.call()

    phone = new(IPhone)
    phone.call()

}

// /* 定义结构体 */
// type Circle struct {
//   radius float64
// }
//
// func main() {
//   var c1 Circle
//   c1.radius = 10.00
//   fmt.Println("Area of Circle(c1) = ", c1.getArea())
// }
//
// //该 method 属于 Circle 类型对象中的方法
// func (c Circle) getArea() float64 {
//   //c.radius 即为 Circle 类型对象中的属性
//   return 3.14 * c.radius * c.radius
// }
//
