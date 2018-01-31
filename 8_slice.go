// s := arr[startIndex:endIndex]
// len切片<=cap切片<=len数组，切片由三部分组成：指向底层数组的指针、len、cap。

package main
import "fmt"

func main() {
   numbers := []int{0,1,2,3,4,5,6,7,8}
   printSlice(numbers)
   fmt.Println("numbers ==",numbers)
   numbers1 := numbers[1:4]
   printSlice(numbers1)
   fmt.Println("numbers[:3] ==", numbers[:3])
   numbers2 :=numbers[4:]
   printSlice(numbers2)
   numbers3 := make([]int,0,5)
   printSlice(numbers3)
   numbers4 := numbers[:2]
   printSlice(numbers4)
}

func printSlice(x []int){
   fmt.Printf("len=%d cap=%d slice=%v\n",len(x),cap(x),x)
}
