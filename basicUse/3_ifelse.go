package main
import "fmt"
func main(){

  // var count,c int   //定义变量不使用也会报错
  var count int
  var flag bool
  count=1

  //while(count<100) {    //go没有while
  for count < 100 {
    count++
    flag = true;
  //注意tmp变量  :=
    for tmp:=2;tmp<count;tmp++ {
      if count%tmp==0{
          flag = false
        }
    }
    //每一个 if else 都需要加入括号 同时 else位置不能在新一行
    if flag == true {
      fmt.Println(count,"素数")
    }else{
      continue
    }
  }

}


// 寻找到所有的素数
// for init; condition; post { }          和C 的for 一样
// for condition { }                         和while 一样
// for { }                                      和C 的for(;;) 一样（死循环）
