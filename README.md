context 上下文:  
https://blog.golang.org/context   

context不仅可以控制并发逻辑，而且本身也可以携带变量，类似于Map。并且提供Value方法用于获取指定Key的Value值.   
context包，可以轻松地将请求范围的值，取消信号以及截止时间跨越API边界传递到处理请求所涉及的所有参数   
Context之间是具有父子关系的，新的Context往往从已有的Context中创建， 因此，最终所有的context组成一颗树状的结构.   
context包中提供一个创建初始Context的方法:  func Background() ContextBackgraund就是所有context树的根。   
WithCancel和WithTimeout两个方法用于在已有的context上创建新的context，同时从新的context中可以获取到旧的context中保存的Key，Value  

context可以被多个并发的Goroutine使用，对context的访问是并发安全的   

c.Abort()只是设置了一些内部标志，标志着上下文以某种方式表明异常终止。  但是在后续中间件中可以根据Isabort() 判断，从而进行 C.Json(500,"") 返回一些error信息   
