# go-trace-demo
golang的分布式服务链路追踪案例，最近调研了网上的关于http和grpc结合的例子比较凌乱，而且相关的demo也比较少，所以就自己出一个案例教程

主要服务

3个http服务和一个rpc服务：

svc_1  
svc_2  
svc_3  
grpc_server  


服务间的调用顺序为：  

svc_1 -> svc_2 -> svc_3 -> grpc_server  

使用方法：  

直接运行根目录下面的main.go文件就可以了，然后在浏览器访问输入http://127.0.0.1:8081/r1 就可以去jaeger控制台看效果图了(前提是已经安装好jaeger)。

下面是jaeger的控制台效果：  

![avatar](https://github.com/Huangsh17/go-trace-demo/blob/master/image/20210108105836.png)

![avatar](https://github.com/Huangsh17/go-trace-demo/blob/master/image/20210108105852.png)