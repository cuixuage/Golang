// func function_name( [parameter list] ) [return_types] {
//    函数体
// }


// package main
// import "fmt"
//
// func main() {
//    /* 定义局部变量 */
//    var a int = 100
//    var b int= 200
//
//    fmt.Printf("交换前，a 的值 : %d\n", a )
//    fmt.Printf("交换前，b 的值 : %d\n", b )
//
//    swap(&a, &b)
//
//    fmt.Println("交换后，a 的值 :", a ,"\n")
//    fmt.Println("交换后，b 的值 :", b ,"\n")
// }
//
// //地址传递
// func swap(x *int, y *int) {
//    var temp int
//    temp = *x
//    *x = *y
//    *y = temp
// }


package main
import (
   "fmt"
   "math"
)

func main(){
   /* 声明函数变量 */
   getSquareRoot := func(x float64) float64 {
      return math.Sqrt(x)
   }

   /* 使用函数 */
   fmt.Println(getSquareRoot(9))          //函数变量的使用

}
