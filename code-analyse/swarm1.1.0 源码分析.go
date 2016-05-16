/*
源码分析部分：
  swam manager/join 启动流程
  Discovery底层实现 
  API请求路径
*/
main.go:
 cli.Run()
//cli/cli.go: 
命令行基于github.com/codegangsta/cli 包 框架 不是自己造轮子，比较好读： 

app.Commands = commands //commands 结构体在commonds.go

//1.2.0 Milestone :P1 Replace samalba/dockerclient by docker/engine-api
/*基于SDK : github.com/samalba/dockerclient 
问题： 跟engine-api 关系? 感觉都是SDK 功能一样。。*/


swam manager:
// cli/manager.go:
 func manage(c *cli.Context) 

 uri := getDiscovery(c)  //没参数的话，os.Getenv("SWARM_DISCOVERY")

//基于github.com/docker/libkv/store,对consul,etcd等封装，对于consul,就是获取api client : client, err := api.NewClient(config)//见下面源码分析
 discovery := createDiscovery(uri, c, c.StringSlice("discovery-opt")) 
   discovery, err := discovery.New(uri, hb, 0, getDiscoveryOpt(c))

//scheduler/strategy/strategy.go: init()加入可支持的filter。返回对应的PlacementStrategy。
 s, err := strategy.New(c.String("strategy")) 
 //scheduler/filter/filter.go :init()加入可支持的filter
 //中规中距，[]Filter。具体filters实现Filter 方法
 fs, err := filter.New(names) 

//返回Scheduler struct; scheduler/scheduler.go
sched := scheduler.New(s, fs)


case "swarm":
		cl, err = swarm.NewCluster(sched, tlsConfig, discovery, c.StringSlice("cluster-opt"), engineOpts)
		        	cluster := &Cluster{
						eventHandlers:     cluster.NewEventHandlers(),
						engines:           make(map[string]*cluster.Engine),
						pendingEngines:    make(map[string]*cluster.Engine),
						scheduler:         scheduler,
						TLSConfig:         TLSConfig,
						discovery:         discovery,
						pendingContainers: make(map[string]*pendingContainer),
						overcommitRatio:   0.05,
						engineOpts:        engineOptions,
						createRetry:       0,
					}
					//这里discovery 为consul，实现Backend接口。即github.com/docker/libkv/store/consul/consul.go:Watch"
					discoveryCh, errCh := cluster.discovery.Watch(nil) //改变就会触发cluster.monitorDiscovery的select

					//健康检查
					//一个Entry 对应一个协程进行健康检查vluster/engine.go : func (e *Engine) refreshLoop() 
					go cluster.monitorDiscovery(discoveryCh, errCh)//for{}循环 // An Entry represents a host.  
					// monitorPendingEngines checks if some previous unreachable/invalid engines have been fixed
					go cluster.monitorPendingEngines()// Engine represents a docker engine
					return cluster, nil


server := api.NewServer(hosts, tlsConfig)

 if c.Bool("replication") {
		....
		setupReplication(c, cl, server, discovery, addr, leaderTTL, tlsConfig)
	} else { 
		//路由注册: api/primary.go: api.NewPrimary。要看
		server.SetHandler(api.NewPrimary(cl, tlsConfig, &statusHandler{cl, nil, nil}, c.GlobalBool("debug"), c.Bool("cors")))
	}

  server.ListenAndServe()





//--- 健康检查 ｄｏｃｋｅｒ　，docker daemon 要加--debug!
go cluster.monitorDiscovery(discoveryCh, errCh)
	for{
		select {
				case entries := <-ch:
					   ...
					   c.addEngine(entry.String())
					       engine := cluster.NewEngine(addr, c.overcommitRatio, c.engineOpts)
					       engine.RegisterEventHandler(c);
					       go c.validatePendingEngine(engine)
					       			//哪里一直尝试连接？ 
					            engine.Connect(c.TLSConfig);
					                dockerclient.NewDockerClientTimeout("tcp://"+e.Addr, config, time.Duration(requestTimeout), setTCPUserTimeout)
									return e.ConnectWithClient(c)  
											//获取docker信息, 调用相关docker  api 对应?
												
												// Gather engine specs (CPU, memory, constraints, ...).
									           e.updateSpecs(); 
									           e.StartMonitorEvents()
									           err := e.RefreshImages();
									           e.RefreshVolumes()
											   e.RefreshNetworks()
											   e.emitEvent("engine_connect")
									           ...   

					            // set engine state to healthy, and start refresh loop
								engine.ValidationComplete()
								    //// refreshLoop periodically triggers engine refresh.

								    go e.refreshLoop() //for{} 同样会调用ｄｏｃｋｅｒ　ａｐｉ　获取相关信息
					   
		}
	}








//-----Discovery底层实现
/*
"github.com/docker/libkv/store"
 type Discovery struct {
	backend   store.Backend
	store     store.Store
	heartbeat time.Duration
	ttl       time.Duration
	prefix    string
	path      string
}
*/
createDiscovery(uri, c, c.StringSlice("discovery-opt"))
  hb, err := time.ParseDuration(c.String("heartbeat")) 
  //github.com/docker/docker/pkg/discovery
  discovery.New(uri, hb, 0, getDiscoveryOpt(c))
    if backend, exists := backends[scheme]; exists {   // discovery/kv/kv.go etcd :Init() 注册。
        err := backend.Initialize(uri, heartbeat, ttl, clusterOpts)
        return backend, err
    }

consul:
Init() 注册:   //在 discovery/kv/kv.go  
     consul.Register() //"github.com/docker/libkv/store/consul" 
       libkv.AddStore(store.CONSUL, New):  
         initializers[store] = init
    
    discovery.Register("consul", &Discovery{backend: store.CONSUL}) //  type Backend string  :CONSUL Backend = "consul"
      //Discovery struct 实现  Backend interface   ：discovery/discovery.go
      backends[scheme] = d


//在 discovery/kv/kv.go:
 Initialize(): 
    //这s.backend 为 consuL "github.com/docker/libkv"
    s.store, err = libkv.NewStore(s.backend, addrs, config) 
    //init 就是 Init 注册的 init 即，New函数
    //consul 的话 ，就在github.com/docker/libkv/store/consul/consul.go
     if init, exists := initializers[backend]; exists {  
		return init(addrs, options)   
	} 
	   //init(addrs, options)  即 func New(endpoints []string, options *store.Config) (store.Store, error)：
	   //就是调用 "github.com/hashicorp/consul/api" SDK包 
	   s := &Consul{}
	   config := api.DefaultConfig()
	   client, err := api.NewClient(config)
	   s.client = client








//----------
swarm join 作用:（调用后可以不用了？） 一直存在，为了swarm manager？重启后可以重发现？
// cli/join.go:	
	d, err := discovery.New(dflag, hb, ttl, getDiscoveryOpt(c))

	for {
		
		 d.Register(addr)
		 	//..最后调用consuL的put。
             s.store.Put(path.Join(s.path, addr), []byte(addr), opts)
		time.Sleep(hb)
	}

//验证情景： 1.docker dameon 好，swarm join 坏；2.swarm join 好，docker dameon 坏；3.docker dameon 坏；swarm join 坏；

//cluster 






一个API请求路径：
	api/primary.go:
	  "/volumes/create":                     postVolumesCreate,
	  "/containers/{name:.*}/stop":          proxyContainerAndForceRefresh,

postVolumesCreate(c *context, w http.ResponseWriter, r *http.Request)
   // func (c *Cluster) CreateVolume(request *dockerclient.VolumeCreateRequest) (*cluster.Volume, error)
   //cluster/swarm/cluster.go
 
   c.cluster.CreateVolume(&request) 
   	   1. 全部主机 range c.engines：
           engine.CreateVolume(request)
       2. 部分 单个指定
        config := cluster.BuildContainerConfig(dockerclient.ContainerConfig{Env: []string{"constraint:node==" + parts[0]}})
		nodes, err := c.scheduler.SelectNodesForContainer(c.listNodes(), config)
		c.engines[nodes[0].ID].CreateVolume(request)

 




  


 
