
## 2015下半年-2016上半年 ##
（该期间，银联基于docker 的 DBAAS云平台 项目: docker生态+ 监控系统）

我想研究docker,就必然会涉及到微服务，devops等关联性比较强的概念，可以从google的 k8s 文档直观感受到 .这三个在这两年能一起火起来，我想都是正回馈的。 

# 微服务&容器  #
  
- [**DevOps**](https://en.wikipedia.org/wiki/DevOps)
- [**12-Factor（SAAS 软件即服务**）](http://12factor.net/zh_cn/)
- [**introduction-to-microservices**](https://www.nginx.com/blog/introduction-to-microservices/)（网上有中文翻译）
- [**microservices**](http://martinfowler.com/articles/microservices.html)
- [**Microservice架构模式简介** ](http://www.cnblogs.com/loveis715/p/4644266.html)
- [**Microservice Architecture - A Quick Guide**](http://colobu.com/2015/04/10/microservice-architecture-a-quick-guide/)


# **算法** #

《大数据日知录__架构与算法》上 paxos,LSM树，bloom filter，ｈａｓｈ，数据分片，数据路由　等等（**强烈推荐这本书**）

**在Youtube上输入关键字 也有很多不错的教学视频。**

## **raft** ##
- [raft 演示图](http://thesecretlivesofdata.com/raft/)
- [《InSearch of an Understandable Consensus Algorithm》](https://ramcloud.stanford.edu/wiki/download/attachments/11370504/raft.pdf)
- [Raft一致性算法](http://blog.csdn.net/cszhouwei/article/details/38374603)
- [Raft一致性算法分析与总结](http://www.thinkingyu.com/articles/Raft/)

 go语言实现包: 
    
- https://github.com/hashicorp/raft （consul项目）         
-  godoc.org/github.com/coreos/etcd/raft

## **paxos** ##

（paxos 算是P2p的算法，复杂，尤其证明它的正确性，还是喜欢相对简单master-slave逻辑。）

- [图解分布式一致性协议Paxos](http://codemacro.com/2014/10/15/explain-poxos/)
- [Paxos](https://en.wikipedia.org/wiki/Paxos_(computer_science) )
- [Paxos算法](http://zh.wikipedia.org/zh-cn/Paxos算法)

## **Gossip** ##
- [gossip](https://github.com/yucs/yucs-awesome-resource/blob/master/algorithms/gossip.pptx)
- [gossip base algorithms](https://github.com/yucs/yucs-awesome-resource/blob/master/algorithms/gossip%20base%20algorithms.pdf)
- [Gossip protocols for large-scale distributed systems](https://github.com/yucs/yucs-awesome-resource/blob/master/algorithms/Gossip%20protocols%20for%20large-scale%20distributed%20systems.pdf)
- [algorithms for cloud computing](https://github.com/yucs/yucs-awesome-resource/blob/master/algorithms/algorithms%20for%20cloud%20computing.pdf)

**go语言相关包: https://github.com/hashicorp/memberlist（consul项目）**


## **hash** ##
- [一致性 hash 算法](http://blog.csdn.net/sparkliang/article/details/5279393)
- [dht-consistent-hash](https://github.com/yucs/yucs-awesome-resource/blob/master/algorithms/dht-consistent-hash.pdf)
- [Couchbase-Server-vBuckets](https://github.com/yucs/yucs-awesome-resource/blob/master/algorithms/Couchbase-Server-vBuckets(hash).pdf)
(二级hash映射,产品实例思路，ceph 的CRUSH算法感觉也类似。)

**go语言实现包: https://github.com/stathat/consistent**





# **docker生态** #

- [10张图带你深入理解Docker容器和镜像](http://dockone.io/article/783)











   

