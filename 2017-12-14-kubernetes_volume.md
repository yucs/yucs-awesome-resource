
---
layout:     post
title:      " Kubernetes之存储学习整理"
subtitle:   " Kubernetes之volume学习整理"
date:       2017-12-14
author:     "yucs"
catalog:    true
categories: 
	- Kubernetes

tags:
    - Kubernetes
      
---



#  概要


- 存储选型思考
 -  一般应用服务：应用级本身不做数据的冗余，为了数据的安全性，而且这类读写延迟高些也能接受（读写IO路径长，多副本机制，都会增加读写延迟），开源的主流使用ceph（默认采用三副本，设计优雅，理念也是自动化）
  -  数据类服务：本身为了高可用而使用多副本冗余机制，通常对性能和延时有比较高的要求
     -  简单方案可以采用如hostpath等本地存储方案，妥协点是数据无法迁移（当然，一般数据类系统 添加删除节点时，本身有负载均衡功能，所以可以通过  删除节点，添加新节点 这种“迁移”方式，迁移过程就是对服务有可能所影响）
     -  使用网络块存储（block device）（性能高的SAN存储）: 跟平台解耦，灵活迁移，代价就是延时有些高，性能有些低（像couchbase 这类内存Nosql，数据在内存，通过异步刷新数据到磁盘 ，对磁盘读写延迟一些可以接受的）

- 开发存储插件
 - 背景：X银行往往有自己的高性能的SAN存储系统，需要进行对接，块设备挂载本地后使用LVM。
      - 基于[FlexVolume](https://github.com/kubernetes/community/blob/master/contributors/devel/flexvolume.md):  [lvm](https://github.com/kubernetes/kubernetes/tree/master/examples/volumes/flexvolume) 根据需求二次定制就好了。
      - （可选）参考[external-storage](https://github.com/kubernetes-incubator/external-storage):[hostPath demo](https://github.com/kubernetes-incubator/external-storage/tree/master/docs/demo/hostpath-provisioner)


  - 参考TiDB关于本地存储的解决方案部分：[黄东旭DTCC2017演讲实录：When TiDB Meets Kubernetes](https://zhuanlan.zhihu.com/p/27229692?utm_source=wechat_session&utm_medium=social)

# 概念
[kubernetes指南：存储](https://feisky.gitbooks.io/kubernetes/concepts/volume.html)

[官方文档](https://kubernetes.io/docs/concepts/storage/volumes/)

[IBM开源技术微讲堂:Kubernetes的存储管理](http://ibm.biz/opentech-ma)（好）


- volume :
 
  - Kubernetes Volume的生命周期与Pod绑定,容器挂掉后Kubelet再次重启容器时，Volume的数据依然还在.

  - 而Pod删除时，Volume才会清理。数据是否丢失取决于具体的Volume类型，比如emptyDir的数据会丢失，而PV的数据则不会丢. (官方文档: which is erased when a Pod is removed, the contents of a cephfs volume(其他网络存储一样，即PV) are preserved and the volume is merely unmounted.)
  - 限制：
     -  声明POD时，暴露出存储细节， 一般用户视角来说，可能不关心，有一定的耦合。
     -  不包含对第三方存储的管理：在声明POD前，对于第三方存储，要先创建好对应的volume,删除POD也需要手动删除volume资源。
 
- PV(Persistent Volumes)，PersistentVolumeClaim (PVC),StorageClass:
  - 为了解决volume的限制，更加方便自动化。
  - 概念：
	  - PersistentVolume（PV）是集群之中的一块网络存储。跟 Node 一样，也是集群的资源。相对于 Volume 会有独立于 Pod 的生命周期（有 PV controller来实现PV/PVC的生命周期）。
	  - 而PersistentVolumeClaim (PVC) 是对 PV 的请求,pod声明使用它。（ 从Storage Admin与用户的角度看PV与PVC :Admin创建和维护PV; 用户只需要使用PVC(size & access mode).
	  - StorageClass来动态创建PV，不仅节省了管理员的时间，还可以封装不同类型的存储供PVC选用。就是封装了对第三方网络存储的管理操作，这样就不用手动创建volume或者手动声明一个PV。（所以，最灵活做法声明POD使用PVC,而PVC使用StorageClass）。


   

# 原理  
[IBM开源技术微讲堂:Kubernetes的存储管理](http://ibm.biz/opentech-ma) 

[Kuberenetes 存储架构总体介绍](http://newto.me/k8s-storage-architecture/)

[Kubernetes存储介绍系列 ——CSI plugin设计](http://newto.me/k8s-csi-design/)

[Kubernetes存储介绍系列 —— AttachDetachController1](http://newto.me/k8s-adcontroller-caches/)

[Kubernetes 存储功能和源码深度解析（一）](http://dockone.io/article/2082)

[Kubernetes 存储功能和源码深度解析（二）](http://dockone.io/article/2087)

![K8S_volume_arch](https://yucs.github.io/picture/K8S_volume_arch.png) 
![k8s_volume_arch2](https://yucs.github.io/picture/k8s_volume_arch2.png)

- Volume Plugins 
 - 存储提供的扩展接口, 包含了各类存储提供者的plugin实现。
 - 实现自定义的Plugins 可以通过FlexVolume(K8s 1.8版本，目前算是过度方案)
 - kubernetes 1.9以后可能推荐CSI（Container Storage Interface）用方式来实现。
     - 支持这套标准以后，K8S和存储提供者之间将彻底解耦，终极目标是将存储的所有的部件作为sidecar container运行在K8S上（当前K8s 1.8版本设计还没有完全做到，需要一个兼容的发展周期），而不再作为K8S部件运行在host上。
- Volume Manager 
   - 运行在kubelet 里让存储Ready的部件，主要是mount/unmount（attach/detach可选）
   - pod调度到这个node上后才会有卷的相应操作，所以它的触发端是kubelet（严格讲是kubelet里的pod manager），根据Pod Manager里pod spec里申明的存储来触发卷的挂载操作
  - Kubelet会监听到调度到该节点上的pod声明，会把pod缓存到Pod Manager中，VolumeManager通过Pod Manager获取PV/PVC的状态，并进行分析出具体的attach/detach、mount/umount, 操作然后调用plugin进行相应的业务处理
- PV/PVC Controller 
  - 运行在Master上的部件，主要做provision/delete
  - PV Controller和K8S其它组件一样监听API Server中的资源更新，对于卷管理主要是监听PV，PVC， SC三类资源，当监听到这些资源的创建、删除、修改时，PV Controller经过判断是需要做创建、删除、绑定、回收等动作。

-  Attach/Detach Controller
  - 运行在Master上，主要做一些块设备（block device）的attach/detach（eg:rbd,cinder块设备需要在mount之前先挂载到主机上，看源码看哪那些实现了Attah接口)
 - 非必须controller: 为了在attach卷上支持plugin headless形态，Controller Manager提供配置可以禁用。  
 - 它的核心职责就是当API Server中，有卷声明的pod与node间的关系发生变化时，需要决定是通过调用存储插件将这个pod关联的卷attach到对应node的主机（或者虚拟机）上，还是将卷从node上detach掉.
 
- K8s挂载卷的基本过程   -  用户创建Pod包含一个PVC   -  Pod被分配到节点NodeA   -  Kubelet等待Volume Manager准备设备   -  PV controller调用相应Volume Plugin(in-tree或者out-of-tree)创建持久化卷并在系统中创建 PV对象以及其与PVC的绑定(Provision)   - Attach/Detach controller或者Volume Manager通过Volume Plugin实现块设备挂载(Attach)  -  Volume Manager等待设备挂载完成，将卷挂载到节点指定目录(mount)  - /var/lib/kubelet/plugins/kubernetes.io/aws-ebs/mounts/vol-xxxxxxxxxxxxxxxxx  -  Kubelet在被告知设备准备好后启动Pod中的容器，利用Docker –v等参数将已经挂载到本地 的卷映射到容器中(volume mapping)

- PV & PVC状态图
  PV的状态图：
    ![pv](https://yucs.github.io/picture/PV_status.png)
   PVC的状态图:
    ![pvc](https://yucs.github.io/picture/pvc_status.png)


# 源码分析
Volume Manager : [kubernetes数据卷管理源码分析](http://www.voidcn.com/article/p-dtoyjptm-bog.html)

 - kubelet管理volume的方式基于两个不同的状态：
   - DesiredStateOfWorld：预期中，pod对volume的使用情况，简称预期状态。当pod.yaml定制好volume，并提交成功，预期状态就已经确定.
   - ActualStateOfWorld：实际中，pod对voluem的使用情况，简称实际状态。实际状态是kubelet的后台线程监控的结果.
 
 - vm.desiredStateOfWorldPopulator.Run方法根据从apiserver同步到的pod信息，来更新DesiredStateOfWorld。另外一个方法vm.reconciler.Run，是预期状态和实际状态的协调者，它负责将实际状态调整成与预期状态。预期状态的更新实现，以及协调者具体如何协调.
 

PVcontroller: [Kubernetes 存储功能和源码深度解析（二）](http://dockone.io/article/2087)


# 开发资源
- volume plugin
 - 基于[FlexVolume](https://github.com/kubernetes/community/blob/master/contributors/devel/flexvolume.md):  [lvm](https://github.com/kubernetes/kubernetes/tree/master/examples/volumes/flexvolume)
 - 基于[CSI](https://github.com/container-storage-interface/spec/blob/master/spec.md): 等待 k8s release 版本支持后

- Out-of-Tree Provisioner
 - 由于在Pod中声明volume有局限性，要更灵活的化，就需要pv controller等进行生命周期的管理。
 - [external-storage](https://github.com/kubernetes-incubator/external-storage):[hostPath demo](https://github.com/kubernetes-incubator/external-storage/tree/master/docs/demo/hostpath-provisioner)


-----
本文出处：https://yucs.github.io/2017/12/14/2017-12-14-kubernetes_volume/
如有出入请请教，文章持续更新中...