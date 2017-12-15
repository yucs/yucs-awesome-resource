# kubernetes extend API 
两种方式：
 api接口
- [Kubernetes API 分析 ( Kube-apiserver )](https://www.kubernetes.org.cn/3119.html)

https://kubernetes.io/docs/api-reference/v1.8/#-strong-api-overview-strong-

https://thenewstack.io/extend-kubernetes-1-7-custom-resources/

https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/

https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md

https://coreos.com/blog/introducing-operators.html

https://blog.openshift.com/kubernetes-deep-dive-api-server-part-1/
https://blog.openshift.com/kubernetes-deep-dive-api-server-part-2/
https://blog.openshift.com/kubernetes-deep-dive-api-server-part-3a/






最新1.8 重构过，代码差异比较大：[Kubernetes1.5源码分析(一) apiServer启动分析](http://dockone.io/article/2159)
[apiserver的list-watch代码解读](https://www.kubernetes.org.cn/174.html)


[Kubernetes Informer 详解](https://www.kubernetes.org.cn/2693.html)

[如何用 client-go 拓展 Kubernetes 的 API](http://www.k8smeetup.com/article/VJsZn@nT7)

 [使用 client-go 控制原生及拓展的 Kubernetes API](https://my.oschina.net/caicloud/blog/829365)
 
 
[Kubernetes Scheduler是如何工作的](http://dockone.io/article/2625)
重构比较大:
从Kubernetes 1.7开始，所有需要监控资源变化情况的调用均推荐使用Informer。Informer提供了基于事件通知的只读缓存机制，可以注册资源变化的回调函数，并可以极大减少API的调用。

![pod_create](/picture/pod_create.png)
 


http://blog.csdn.net/WaltonWang/article/list/1

http://cizixs.com/ 
 
 
我们发现为什么部署这些有状态的应用和部署管理它们会比无状态的复杂呢？是因为它们有这些复杂的运维和逻辑在里面 
#controller
https://github.com/kubernetes/community/blob/8decfe4/contributors/devel/controllers.md

```go
for {
  desired := getDesiredState()
  current := getCurrentState()
  makeChanges(desired, current)
}
```



[](https://zhuanlan.zhihu.com/p/27229692?utm_source=wechat_session&utm_medium=social)

[使用 Operator 来扩展 Kubernetes(视频)](https://k8smeetup.maodou.io/course/hFRDJyzkdWXPFanyY)
