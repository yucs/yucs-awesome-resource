#理解

- **协调系统**:像etcd/consul等第三方协调程序，对分布式系统设计来说基本是标配组件， 能简化设计，降低逻辑复杂度有很大的作用（像kubernetes,swamkit,mesos，kafka,自己设计的系统,就如ceph分布式文件系统的monitor某种意义也是个协调系统）：可以简单高效的实现 分布式锁服务，领导者选举，配置管理（配置信息自动获取，被改变时，notify其影响的进程），组成员管理，任务分配，节点状态信息，基于协调的Lease机制。
- **消息队列**: 系统级解耦 + 异步通信 + 峰值处理能力 +[送达保证：至少送达一次|丢失|保证只送达一次...]，也是分布式系统常见组件，大数据处理kafka基本主流。
- **RPC**:rpc作为系统内部通讯协议， 避免自己实现底层通讯协议，高性能，跨语言等优点 成为主流选择，在面向微服务架构流行下，google开源grpc成为趋势（swamkit,etcd3.0，k8s等最近项目基本都采用grpc）
- **Nosql**：在海量数据需要存储（当初关系型数据库设计本身不考虑集群:ACID事务 很难在集群上实现） 和 更友好的数据交互提高开发效率（“无模式”数据模型） 趋势下， Nosql渐渐在大数据方向下崛起，更多强调可用性和可扩展性。根据数据模型 分为：键值，文档，列族（该三类 都是面向聚合），图四大类型。数据类分布式系统，技术面的话基本都要考虑 数据路由+数据分片+多副本策略（主流主从）+一致性模型（读我所写，session一致性）+ membership(Gossip协议)+查询机制+索引（skipList|b+树） +[其他：CAS机制(Campare and set)]，数据保存在内存也是个方向。


#具体技术
---
##协调系统

- **《大数据日知录__架构与算法》第五章 分布式协调系统**
- **作用**：在大规模系统分布式系统下，通过协调系统可以简单高效的实现 分布式锁服务，领导者选举，配置管理（配置信息自动获取，被改变时，notify其影响的进程），组成员管理，任务分配，节点状态信息，基于协调的Lease机制 等等。（5.2.4&&5.2.5 章节）

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

---

##分布式通信

- **《大数据日知录__架构与算法》第六章 分布式通信**
- 大数据系统抽象归纳出3种常见的通信机制：**序列化与远程过程调用，消息队列，多播通信**。分别典型代表：thrift/grpc, kafka , gossip协议。

###RPC: thrift/grpc
   -  [Thrift架构介绍](http://www.91it.org/articles/thrift-framework-intro.html)
   -  [gRPC vs Thrift](http://blog.csdn.net/dazheng/article/details/48830511) & [Why gRPC?](http://www.grpc.io/posts/principles)。
   -  个人理解：重要一点都是**跨语言**的； grpc基于http2与 protocol buffers,文档比thrift完善，开发比较友好，更多强调 **面向微服务架构**，主流趋势。thrift基于raw socket性能更高，目前语言的支持更多,相对复杂特性多而文档不是很完善，代码生成也比较多，panic问题。

###kafka

-  **主要作用: 系统级解耦 + 异步通信 + 峰值处理能力 +扩展性+[送达保证：至少送达一次|丢失|保证只送达一次...]+[顺序保证]**
-  **通用技术：生产者消费者模式 + pub-sub订阅模式+ pull|push [+持久化]**
- **RabbitMQ/kafka**: 在大数据处理在收集各类资源，以kafka主流；RabbitMQ重量级系统，遵循AMQP协议，具有较强功能和相对广泛的使用场景（企业级偏功能而非收集），但是扩展性和性能相对低。
- **参考资源**:**[Kafka深度解析](http://www.jasongj.com/2015/01/02/Kafka深度解析)**（消息队列作用+常用Message Queue对比+Kafka解析 都总结的很好 ）&& [Kafka技术内幕](http://zqhxuyuan.github.io/2017/01/01/Kafka-Code-Index/)&&**[apache kafka技术分享目录索引（好）](http://blog.csdn.net/lizhitao/article/details/39499283)**(ps:网上相关文章还是比较多且质量都比较高的)
- **kafka主要技术**：
	 - 生产者 根据topic **push**数据到 Broker,消费者 从Broker **pull** 数据； 至少送达一次；
	 - 支持topic进行数据分片，并且数据是有序，不可更改的**追加**到消息队列上，并存储到具体文件上进行持久化；消费者端内存维护Offset索引，可以通过修改索引来读过期数据；（**高效处理大批数据的重要原因就是将读写操作尽可能转化为顺序读写**）,利用linux内核的文件系统自身的缓存机制，sendfile system call减少数据的内核用户的拷贝。
	 - topic上消息会**广播**到所有的comsumer groups,每个comsumer group只会把消息**单传**到一个comsumer.
	 - 通过Zokeeper来存放配置信息，offset索引，侦测Broker 的加入和删除 来扩展；
	 - 副本机制也是leader-follower模式，读写都有leader来响应，通过ISR机制来保证多副本的一致性,flolloer 间断从leader拉数据；(无处不显示Kafka高吞吐量设计思想)
	 -[apache Kafka概要介绍](http://blog.csdn.net/lizhitao/article/details/23743821)(**好**)

	 - [design](http://kafka.apache.org/documentation.html#design)(官网 讲的比较好了)&&[为什么Kafka那么快](http://mp.weixin.qq.com/s?__biz=MzIxMjAzMDA1MQ==&mid=2648945468&idx=1&sn=b622788361b384e152080b60e5ea69a7#rd)
   		
   		
   		 
 
- **官网**：
  - grpc :[http://www.grpc.io/](http://www.grpc.io/)&[protocol buffers](https://developers.google.com/protocol-buffers/docs/overview)
  - thrift:[https://thrift.apache.org/](https://thrift.apache.org/)
  - kafka :[http://kafka.apache.org/](http://kafka.apache.org/)

---


###NOSQL

- **《大数据日知录__架构与算法》 第九，十章 内存KV数据库 列式数据库**

- **《nosql精粹》**：大规模数据推进Nosql和Newsql的技术发展，这本书作为NoSQL领域入门科普书籍，通俗易懂，讲解nosql类型分类,基本技术，和应用场景；

- 都要要回答的问题是：
  - 集群中，数据如何路由分片路由（要考虑 节点加入删除的数据迁移，负载均衡的问题）
  - 数据的多副本策略 与 一致性模型分类/一致性的保证。
  - 数据模型 与 查询索引逻辑。
   
.
 

- **couchbase**：基于内存的文档型数据库
  - 数据路由分布(二级哈希):[Technical-Whitepaper-Couchbase-Server-vBuckets](http://www.couchbase.com/sites/default/files/uploads/all/whitepapers/Technical-Whitepaper-Couchbase-Server-vBuckets.pdf)
  - **主要技术**：由多个服务组件组成的，**基于内存**；**分布式索引**（通过skipLIist数据结构存在内存）；数据路由分片基于vbucket(二级hash) ；N1QL(类SQL语句查询)；
  -  官方文档：无论操作还是基本技术原理 ，培训视频 等各种资源文档比较都是完善！
- **hbase**：列族数据库
  - 。。。
- **官网**：
  - [couchbase ](http://www.couchbase.com/) 
  - [hbase](http://www.couchbase.com/)