**NOSQL**
couchbase
hbase


**协调系统**

- **《大数据日知录__架构与算法》第五章**
- **作用**：在大规模系统分布式系统下，简单高效的实现 分布式锁服务，领导者选举，配置管理（配置信息自动获取，被改变时，notify其影响的进程），组成员管理，任务分配，节点状态信息，基于协调的Lease机制 等等。（5.2.4&&5.2.5 章节）

- **系统**：chubby作为google的系统，其论文chubby 催生开源zookeeper,两者都基于paxos协议；consul/etcd 基于raft协议的协调系统。
- **技术**：
 - 共同：
     - 类似文件系统的目录层次管理+notify/watch机制；
     - 数据存内存+日记恢复方式 
     - 。。。
 - **chubby**（**主从+paxos+Lease机制**）：
     - 通过paxos 选举主控服务器，所有读写请求有该服务响应；数据一致性通过paxos协议来保证；采用Lease机制：主控服务器租约，客户端缓存的一致性也通过Lease来保证，大大减少服务端的压力,提高并发度；
	 
 - **zookeeper(paxos + 每台都能响应请求+也有主从)**: 
 		- 任意服务器都能响应请求，写只能住响应，存在读到过期数据风险(API接口提供sync);一致性也是通过paxos来保证.
  -  zookeeper强调高吞吐；Chubby强调系统可靠性和高可用性及语义已于理解；
  - **consul/etcd**:一致性基于raft协议/go实现；协议consul服务发现还基于gossip ：
   - [etcd：从应用场景到实现原理的全方位解读](http://www.infoq.com/cn/articles/etcd-interpretation-application-scenario-implement-principle)
 - [Zookeeper与paxos算法](http://blog.jobbole.com/45721/)
- **对比：**[服务发现:Zookeeper vs etcd vs Consul](http://dockone.io/article/667):
      - Zookeeper的主要优势是其成熟、健壮以及丰富的特性，然而，它也有自己的缺点，其中采用Java开发以及复杂性是罪魁祸首.
      - 与Zookeeper和etcd不一样，Consul内嵌实现了服务发现系统，所以这样就不需要构建自己的系统或使用第三方系统务.
      - etcd 轻量级，简洁，高效
- **官网**：
  - etcd : [https://github.com/coreos/etcd](https://github.com/coreos/etcd)
  - consul:[https://www.consul.io/](https://www.consul.io/)
  - zookeeper :[http://zookeeper.apache.org/](http://zookeeper.apache.org/)
  - chubby论文：[Chubby：面向松散耦合的分布式系统的锁服务](http://duanple.blog.163.com/blog/static/70971767201142412058672/)

**分布式通信**

- **《大数据日知录__架构与算法》第六章**
- 大数据系统抽象归纳出3种常见的通信机制：序列化与远程过程调用，消息队列，多播通信。分别典型代表：thrift/grpc, kafka , gossip协议。
- **RPC: thrift/grpc**  
   -  [Thrift架构介绍](http://www.91it.org/articles/thrift-framework-intro.html)
   -  [gRPC vs Thrift](http://blog.csdn.net/dazheng/article/details/48830511) & [Why gRPC?](http://www.grpc.io/posts/principles)。
   -  个人理解：重要一点都是**跨语言**的； grpc基于http2,文档比thrift完善，开发比较友好，更多强调 **面向微服务架构**，主流趋势。thrift基于raw socket性能更高，目前语言的支持更多,相对复杂特性多而文档不是很完善，代码生成也比较多，panic问题。
- **消息队列：kafka** 
    -  **主要作用: 系统级解耦 + 异步通信 + 峰值处理能力 +[送达保证：至少送达一次|丢失|保证只送达一次...]**
    -  **通用技术：生产者消费者模式 + pub-sub订阅模式+ pull|push [+持久化]**
    - **RabbitMQ/kafka **: 在大数据处理在收集各类资源，以kafka主流；RabbitMQ重量级系统，遵循AMQP协议，具有较强功能和相对广泛的使用场景（偏功能而非收集），但是扩展性和性能相对低。
    - [Kafka深度解析](http://www.jasongj.com/2015/01/02/Kafka深度解析)（消息队列作用+常用Message Queue对比+Kafka解析 都总结的很好 ）

- **官网**：
  - grpc :[http://www.grpc.io/](http://www.grpc.io/)
  - thrift:[https://thrift.apache.org/](https://thrift.apache.org/)
  - kafka :[http://kafka.apache.org/](http://kafka.apache.org/)
