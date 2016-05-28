# go-app-scatter
	go http app echarts scatter 
	golang 并发测试操作并使用echarts输出表表统计运行时间
	
##目录介绍
	运行app生成并发报表 分析基本每次运行执行时间
	out 结果统计输出目录 可在config/download.config 修改保存目录
	config 配置目录

	golang代码修改，可添加自己的测试
	修改config文件夹里面的config/output.html里面需要替换的内容 用%s代替


##修改扩展
	
##测试结果
## 预先创建好大小为1000的 redis静态连接池，并发10000 的LPUSH操作，可在 0.7 秒内完成；
## 3  Intel(R) Xeon(R) CPU           X5650  @ 2.67GHz
## CPU、内存、网络吞吐量还未做详细统计
## 代码版本号为 7ed0e05c56779c0901200e98b671d74e2ff50dc6

##测试scatter图链接

* [10000 并发图](http://138.128.192.237:30300/10000-20160528_070201.html)


