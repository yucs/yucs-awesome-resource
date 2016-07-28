#分布式文件系统#
- 都要要回答的问题是：  
  - 添加/删除节点的数据负载均衡的问题
  - 文件数据以 何种方式如何分布c存储在集群节点中,是否有元数据节点
  - 数据的多副本策略 与 一致性
  - IO读写流程
## ceph ##
## glusterfs ##


### HDFS ###
**主要参考**

- **《大数据日知录》第八章 分布式文件系统**
- **google论文：GFS**
- **HADOOP权威指南 第三章 分布式Hadoop文件系统**

**google论文 GFS要点**

- **技术需求** 
     - 组件失败成为一种常态而不是异常
     - 文件是巨大的。大部分的文件更新模式是通过在尾部追加数据而不是覆盖现有数据，文件内部的随机写操作几乎是不存在的。 一旦写完，文件就是只读的，而且通常是顺序读。
- **设计假设**
	 - 系统是由廉价的经常失败的商品化组件构建而来。
	 - 系统存储了适度个数的大文件
	 - 工作负载主要由两种类型的读组成：大的顺序流式读取和小的随机读取
	 - **追加写 ，少随机写；多顺序读，少随机读**：工作负载有很多大的对文件数据的 append 操作，系统必须为多个客户端对相同文件并行 append 操作的定义良好。
	 - 持续的高带宽比低延时更重要。

- **设计原理**
  - **数据路由分片**：分片：文件被划分成固定大小的 chunk，默认地我们存储三个备份，master记录位置；路由：client向master获取位置。 
  -  **简化设计**：采用元数据节点,Master 维护所有的文件系统元数据。包括名字空间，访问控制信息，文件与 chunk的映射信息， chunk 的当前位置。它也控制系统范围内的一些活动，比如 chunk租约管理， 无效 chunk 的垃圾回收， chunkserver 间的 chunk 迁移。 Master 与chunkserver 通过心跳信息进行周期性的通信，以发送指令和收集chunkserver 的状态。
  -  **大的 chunk size**： 降低了 client 与 master 的交互需
求，减少应用产生的负载是非常明显的，很有可能在一个给定的 chunk 上执行更多的操作，允许我们将元数据存放在内存中。
  -  Master 存储了三个主要类型的元数据：文件和 chunk 名字空间，文件到 chunk的映射信息，每个 chunk 的备份的位置。
  - **租约机制**： 我们使用租约机制来保持多个副本间变更顺序的一致性。 Master 授权给其中的一个副本一个该 chunk 的租约，我们把它叫做主副本(primary)。
  - 元数据存储在内存里， master 的操作是很快的； Master 并没有永久保存 chunk 的位置信息，而是在 master启动或者某个 chunkserver 加入集群时，它会向每个 chunkserver 询问它的 chunks信息。
  - 操作日志包含了关键元数据改变的历史记录，Master 通过重新执行操作日志来恢复它的文件系统。

 **HDFS **

- 作为GFS的开源版，整体框架，使用场景类似：大文件；商用硬件；追加写 ，少随机写；多顺序读，少随机读； 优化带宽而非延迟 等。

- 主控服务器 当点失效问题：**HA方案(NFS备份)**或 **NameNode 联盟（二级名称节点）**

- [HDFS详解](http://my.oschina.net/crxy/blog/348868)&&[【漫画解读】HDFS存储原理](http://www.36dsj.com/archives/41391)（[英文版](http://www.slideshare.net/jaganadhg/hdfs-10509123)）
- [Hdfs Design(官网)](http://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-hdfs/HdfsDesign.html)
  

HayStack 对象存储系统
Erasure Code