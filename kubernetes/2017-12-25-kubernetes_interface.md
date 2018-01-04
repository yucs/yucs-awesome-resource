---
layout:     post
title:      "kubernetes的思考和那些标准"
subtitle:   "kubernetes的思考和那些标准"
date:       2017-12-25
author:     "yucs"
catalog:    true
categories: 
	- Kubernetes

tags:
    - Kubernetes
     
---
 
# 思考
- CNCF生态被业界认可，在容器编排领域，K8s已然是事实标准:
   - 部署方面：kubernetes及社区 已经能够简单快速的搭建和部署 ，这样才能更好的普及。
   - 稳定标准化方面：各种标准化被业界统一认可，核心功能整体架构等 趋于稳定；标准化某种意义就是避免被绑架。
   - 开发定制方面:整体架构的优雅灵活，系统可扩展性好，根据需求在不侵入源码下，可定制性强（client-go项目对开发友好），这样才能更好的解决各公司内部需求，实践落地成功。   
   - 活跃度方面：贡献参与者多，社区活跃度高，各大公司的认可，无论存储，网络，容器方面都提供更多现成的选择，功能也越来越丰富，生态整体向好。


- 在kubernetes和docker两大生态竞争中，毫无疑问kubernetes已胜出,k8s已是大趋势。

- 开源也是场各大公司之间的博弈，无烟的战场，谁掌握了标准的制定，就占领了制高点,一言不合就可以让竞争对手举步维艰。



# [CNCF](https://www.cncf.io/)

[Cloud Native Landscape](https://github.com/cncf/landscape)

[CNCF charter](https://www.cncf.io/about/charter/)

[CNCF 云原生容器生态系统概要](http://dockone.io/article/3006)

   - CNCF是一个开源Linux基金会，它致力于推进云端原生应用和服务的开发.
   - CNCF 的一项重要承诺，就是为基于容器的各类技术的集成确立参考架。
   - CNCF’s community believe there are three core attributes to cloud native computing:
    - Container packaged and distributed.
    - Dynamically scheduled.
    - Micro-services oriented.
   - A cloud native computing system enables computing that builds on these core attributes, and embraces the ideals of:
     - Openness and extensibility.
     - Well-defined APIs at borders of standardized subsystems.
     - Minimal barriers to application lifecycle management.



# 容器标准（OCI）
[OCI 和 runc：容器标准化和 docker](http://cizixs.com/2017/11/05/oci-and-runc)

[docker、oci、runc以及kubernetes梳理](http://www.cnblogs.com/xuxinkun/p/8036832.html)


 - [OCI](https://www.opencontainers.org/about)（Open Container Initiative）
   - 是由多家公司共同成立的项目，并由linux基金会进行管理，致力于container runtime的标准的制定和runc的开发等工作。
   - 使命就是推动容器标准化，容器能运行在任何的硬件和系统上，相关的组件也不必绑定在任何的容器运行时上.
   - 目前主要有两个标准文档：容器运行时标准 [runtime-spec](https://github.com/opencontainers/runtime-spec)和 容器镜像标准[image-spec](https://github.com/opencontainers/image-spec)

- [runc](https://github.com/opencontainers/runc): 是对于OCI标准的一个参考实现.

- [CRI](https://kubernetes.feisky.xyz/plugins/CRI.html)(Container Runtime Interface)
  - kubernetes自己的运行时接口api,通过统一的接口与各个容器引擎之间进行互动。
  -  基于 gRPC 定义了 RuntimeService 和 ImageService，分别用于容器运行时和镜像的管理.
  -  与oci不同，cri与kubernetes的概念更加贴合，并紧密绑定。cri不仅定义了容器的生命周期的管理，还引入了k8s中pod的概念，并定义了管理pod的生命周期.


![K8S_CRI](https://yucs.github.io/picture/K8S_CRI.png)

 
 


# CNI(容器网络接口)

 [CNI](https://github.com/containernetworking/cni)
 
 
 [kubernetes指南:CNI](https://kubernetes.feisky.xyz/network/cni/)
 
 - Container Network Interface (CNI)最早是由CoreOS发起的容器网络规范，是Kubernetes网络插件的基础。其基本思想为：Container Runtime在创建容器时，先创建好network namespace，然后调用CNI插件为这个netns配置网络，其后再启动容器内的进程。现已加入CNCF，成为CNCF主推的网络模型.



# CSI(容器存储接口)

[CSI spce](https://github.com/container-storage-interface/spec/blob/master/spec.md)

[Kubernetes存储介绍系列 ——CSI plugin设计](http://newto.me/k8s-csi-design/)

[design-proposals:csi](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/storage/container-storage-interface.md)

- CSI是Container Storage Interface的简称，旨在能为容器编排引擎和存储系统间建立一套标准的存储调用接口，通过该接口能为容器编排引擎提供存储服务。
- 在CSI之前，K8S里提供存储服务是通过一种称为“in-tree”的方式来提供，这种方式需要将存储提供者的代码逻辑放到K8S的代码库中运行，调用引擎与插件间属于强耦合，持这套标准以后，K8S和存储提供者之间将彻底解耦。





----

markdown文件放在 [github.com/yucs/yucs-awesome-resource](https://github.com/yucs/yucs-awesome-resource) 持续更新，欢迎star ,watch

如有出入请请教，文章持续更新中...
