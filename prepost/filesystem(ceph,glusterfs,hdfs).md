
##HDFS&glusterfs&ceph 个人理解对比

- glusterfs: 无元数据节点，兼容PAXOS接口的文件系统，可以理解无分片（无建条带层情况下），目录层次，ls等获取列表很卡，架构特点像网络协议栈，一层层处理后数据传到下一层。


- ceph:基于RADOS，无元数据节点，基于CRUSH算法，**分片和二级映射**，底层默认存储单位4M。个人觉得最优雅在于：monitor存储的各种表（类似kubernetes,swamkit的desired state一样）, 客户端和OSD 获取这状态（单异常时，monitor会更新，主动或被动通知他们），osd和客户端就会以新表 数据自动迁移，自动修复，自动获取，即客户端和OSD之间无需通信 以规定的协议 来就好了，问题也导致了当有节点掉了，OSD之间自动不受控制的数据修复，很可能导致性能下降，甚至整个集群时间内不可用。

- HDFS:设计目的是：大文件；商用硬件；追加写 ，少随机写；多顺序读，少随机读； 优化带宽而非延迟 等，可见应用场景不适合做docker，KVM等镜像文件这样有随机读写的。有MASTER节点来控制数据迁移，修复等，整体设计相对简单些，控制都由master。


- 都要要回答的问题是：  
  - 添加/删除节点的数据负载均衡的问题
  - 文件数据以 何种方式如何分布c存储在集群节点中,是否有元数据节点
  - 数据的多副本策略 与 一致性
  - IO读写流程


参考：**[杨锦涛：主流开源存储方案孰优孰劣](http://www.infoq.com/cn/interviews/interview-with-yangjintao-talk-open-source-storage-scheme#0-youdao-1-28677-32553cecb956bf88a1550052113e506a)**



## ceph ##
介绍框架：[architecture](http://docs.ceph.com/docs/master/architecture/)&[Ceph架构剖析](https://www.ustack.com/blog/ceph_infra/)




 - **客户端**（无元数据节点，直接根据计算获取读写位置）：
client通过monitor获取CRUSH map等集群信息,用户看到的文件 根据文件位置 进行**分片**（ceph集群默认存储object:4M）：通过librados库根据CRUSH算法（**二级映射**: object（ceph集群里的object）->PG->OSDS; **CRUSH**：根据一棵主机拓扑树，递归算出osd列表）,映射到具体的OSDs 上，客户端 直接与osd通信，读写数据。
- **RADOS**：Ceph Monitor 和 Ceph OSD Daemon集群核心,monitor基于paxos等维护表信息，osds状态等,类似协调系统，数据修复以PG为单位。Block Devices,Object Storage,Filesystem都在这RADOS上封装，
- osd目前通过filestore机制,主备通过plog协议来保证一致性;因写log,多一次IO，写放大问题，下个版本默认将是 Newstore：[ceph存储引擎bluestore解析](http://www.sysnote.org/2016/08/19/ceph-bluestore/)&[Ceph Jewel 版本预览 : 即将到来的新存储BlueStore](http://bbs.ceph.org.cn/article/63)&[BlueStore: a new, faster storage backend for Ceph]()
参考：
 - [Ceph的IO模式分析](http://www.openstack.cn/?p=4270)
 - [Ceph对象存储运维惊魂72小时](http://ceph.org.cn/2016/05/08/ceph%E5%AF%B9%E8%B1%A1%E5%AD%98%E5%82%A8%E8%BF%90%E7%BB%B4%E6%83%8A%E9%AD%8272%E5%B0%8F%E6%97%B6/)


**主要资源**

 - **官网**：[ceph](http://docs.ceph.com/docs/master/)
 - **《learning ceph》**
 - **微信公众号**：**ceph开发每周谈**，ceph社区
 - 网上相关源码分析
 - [http://www.sysnote.org/](http://www.sysnote.org/)跟ceph相关的文章,质量都不错
 - 麦子迈的[http://www.wzxue.com/ceph-storage/](http://www.wzxue.com/ceph-storage/)，现在在公众号**ceph开发每周谈** 发表了。

---
### HDFS&GFS ###
**主要参考**

- **《大数据日知录》第八章 分布式文件系统**
- **google论文：GFS**
- **HADOOP权威指南 第三章 分布式Hadoop文件系统**

####google论文 GFS要点 ####

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

####HDFS####

- 作为GFS的开源版，整体框架，使用场景类似：大文件；商用硬件；追加写 ，少随机写；多顺序读，少随机读； 优化带宽而非延迟 等。

- 主控服务器 当点失效问题：**HA方案(NFS备份)**或 **NameNode 联盟（二级名称节点）**

- [HDFS详解](http://my.oschina.net/crxy/blog/348868)&&[【漫画解读】HDFS存储原理](http://www.36dsj.com/archives/41391)（[英文版](http://www.slideshare.net/jaganadhg/hdfs-10509123)）
- [Hdfs Design(官网)](http://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-hdfs/HdfsDesign.html)
  

