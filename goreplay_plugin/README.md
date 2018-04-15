goreplay工具:     
https://github.com/buger/goreplay/wiki   

实现：   
gin + vuejs  
 
work：    
1.添加web界面     
2.后台封装goreplay process struct,提供start stop restart等等function     
3.goroutine维护process alivelist，如果终止则会及时拉起。time ticker扫描对应目录下的文件更新    
4.所有配置信息持久化到本地，不会造成重新启动丢失配置情况   
4.systemed --user start xxx.service   