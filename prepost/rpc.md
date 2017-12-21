# rpc


https://my.oschina.net/u/1378920/blog/904669#comment-list

https://zhuanlan.zhihu.com/p/29028054

https://www.zhihu.com/question/25536695


https://mp.weixin.qq.com/s?__biz=MzIwNDU2MTI4NQ==&mid=2247483772&idx=1&sn=ee3d6e3937dffb2d45d555ae753482ad&chksm=973f0f96a04886804f90e36e1e79408fe4bb666c0a6b8d1ecc17c40b74f761dcf3658aa64d74#rd

分布式环境。

写起来就跟调用本地函数一样。

程序员无需关注与远程的交互细节。




- **RPC**:rpc作为系统组件通讯协议， 避免自己实现底层通讯协议，高性能，跨语言等优点 成为主流选择，在面向微服务架构流行下，google开源grpc成为趋势（swamkit,etcd3.0，k8s等最近项目基本都采用grpc），可以对外调用的系统组件往往都会提供rpc接口，一种高性能调用接口。
 



###RPC: thrift/grpc
   -  [Thrift架构介绍](http://www.91it.org/articles/thrift-framework-intro.html)
   -  [gRPC vs Thrift](http://blog.csdn.net/dazheng/article/details/48830511) & [Why gRPC?](http://www.grpc.io/posts/principles)。
   -  个人理解：重要一点都是**跨语言**的； grpc基于http2与 protocol buffers,文档比thrift完善，开发比较友好，更多强调 **面向微服务架构**，主流趋势。thrift基于raw socket性能更高，目前语言的支持更多,相对复杂特性多而文档不是很完善，代码生成也比较多，panic问题。