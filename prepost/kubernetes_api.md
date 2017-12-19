
#通过CustomResources&&controller扩展kubernetes

根据需求，与其他Pod等公民一样，先用CustomResources扩展添加新Resource,用controller来达到预期状态。

Operator核心就是通过CustomResources扩展添加新Resource组合实现。



# CustomResources

[kubernetes 指南：customresourcedefinition](https://kubernetes.feisky.xyz/concepts/customresourcedefinition.html)

官方相关文档： [custom-resources](https://kubernetes.io/docs/concepts/api-extension/custom-resources), [extend-api-custom-resource-definitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/)


代码级例子 : [extend-kubernetes-1-7-custom-resources](https://thenewstack.io/extend-kubernetes-1-7-custom-resources/)

 - 无需改变代码来扩展 Kubernetes API 的机制，用来管理自定义对象.
 
 - For any new resource, you follow the same methodology:
   - Define the resource schema;
   - Register the resource with the API service and provide proper APIs;
   - Implement a controller which will watch for resource spec changes and make sure your application complies with the desired state.

- [sample-controller](https://github.com/kubernetes/sample-controller)

-  通过工具生成相关代码([code-generator](https://github.com/kubernetes/code-generator))： [code-generation-customresources](https://blog.openshift.com/kubernetes-deep-dive-code-generation-customresources/)



# controller
[A Deep Dive Into Kubernetes Controllers
](https://engineering.bitnami.com/articles/a-deep-dive-into-kubernetes-controllers.html)

[kubewatch-an-example-of-kubernetes-custom-controller](https://engineering.bitnami.com/articles/kubewatch-an-example-of-kubernetes-custom-controller.html)


官方社区给出的开发controller指导： [kubernetes/community:controllers](https://github.com/kubernetes/community/blob/8decfe4/contributors/devel/controllers.md)


 
- Kubernetes runs a group of controllers that take care of routine tasks to ensure the desired state of the cluster matches the observed stat.(each controller is responsible for a particular resource in the Kubernetes world).
 


- 伪代码模型：

```go
for {
  desired := getDesiredState()
  current := getCurrentState()
  makeChanges(desired, current)
}
```

- client-go包的Informer/SharedInformer 
  - Informer/SharedInformer watches for changes on the current state of Kubernetes objects and sends events to Workqueue where events are then popped up by worker(s) to process.（从Kubernetes 1.7开始，所有需要监控资源变化情况的调用均推荐使用Informer。Informer提供了基于事件通知的只读缓存机制，可以注册资源变化的回调函数，并可以极大减少API的调用。）
  - [Kubernetes Informer 详解](https://www.kubernetes.org.cn/2693.html)

- 处理函数：
  - client-go包封装了获取事件变化和针对对异步的队列框架机制，我们只需实现处理逻辑接口：

 ```go
     type ResourceEventHandlerFuncs struct {
	AddFunc    func(obj interface{})
	UpdateFunc func(oldObj, newObj interface{})
	DeleteFunc func(obj interface{})
}
 ```
 
  ![框架图](genenal_pattern_controller.png)



# 其他相关开发资源

- [client-go](https://github.com/kubernetes/client-go)
  - [使用 client-go 控制原生及拓展的 Kubernetes API](https://my.oschina.net/caicloud/blog/829365)
  - [如何用 client-go 拓展 Kubernetes 的 API](http://www.k8smeetup.com/article/VJsZn@nT7) 

- **[使用 Operator 来扩展 Kubernetes(视频)](https://k8smeetup.maodou.io/course/hFRDJyzkdWXPFanyY)**
  

 



# opertor

[Operator：固化到软件中的运维技能](http://blog.fleeto.us/translation/introducing-operators-putting-operational-knowledge-software)

 Operator 是跟应用紧密相关的，所以其中最重要的工作就是把应用自身的运维方法编码成为资源和控制逻辑
 
 我们发现为什么部署这些有状态的应用和部署管理它们会比无状态的复杂呢？是因为它们有这些复杂的运维和逻辑在里面 





https://github.com/spotahome/redis-operator


