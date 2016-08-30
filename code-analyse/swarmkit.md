## 参考 ##
- [Docker 1.12.0将要发布的新功能](http://liubin.org/blog/2016/06/17/whats-new-in-docker-1-dot-12-dot-0/)

- [docker-built-in-orchestration](https://blog.docker.com/2016/07/docker-built-in-orchestration-ready-for-production-docker-1-12-goes-ga/)

- [ swarmkit architecture](http://www.containertutorials.com/swarmkit/architecture.html)


### **看代码说明：** ###
   -  **类似python ，理解就是 python中通过缩进代表{}，而这里代表函数,即缩进的该行上上行的被调用子函数**
		

		//eg: 该例子就是func1调用subfunc1，subfunc2，而subfunc1，subfunc2是同层次的函数
		func1(1,2)
   		   subfunc1(1,2)
	       subfunc2(3,4)


#### swamkit代码整体框架

操作逻辑：先持久化到raft，再提交到内存数据库（raft 实现 Proposer 接口），产生变化事件，通过go-events包订阅者模式发送到队列中，manager节点 会开启相关管理协层（orchestrator，Allocator）订阅相关事件消费处理，可以理解为内部消息队列。 

eg:[docker-built-in-orchestration](https://blog.docker.com/2016/07/docker-built-in-orchestration-ready-for-production-docker-1-12-goes-ga/)中创建流程图，API层创建服务，数据存入Raft和数据库后，立即返回（见下源码分析）,即异步处理；orchestrator goroutine（manager 启动，Watch相关Queue，订阅者模式）会收到该event,然后处理，也会产生新的event：

**manager/state目录**

	   raft目录：raft协议  实现proposer.go:Proposer 接口
	   store目录：基于go-memdb内存，各个文件一个表，init()初始化表；
	  		   // MemoryStore is a concurrency-safe, in-memory implementation of the Store interface.
				type MemoryStore struct 封装create,View,Update等操作，
       //更新操作逻辑：先持久化到raft，再提交内存（raft 实现 Proposer 接口），产生变化事件，
	   //通过go-events包订阅者模式发送到队列中，相关管理模块订阅相关事件消费处理如orchestrator。
       func (s *MemoryStore) update(proposer state.Proposer, cb func(Tx) error)
       		
		    var tx tx
			tx.init(memDBTx, curVersion)
			//内存数据库操作
			err := cb(&tx)
       		// (tx *tx) update(table string, o Object)等封装都会
       		/// tx.changelist = append(tx.changelist，。。）
       		 sa, err = tx.changelistStoreActions()

       		 //先持久化到raft，再提交内存	
        	 proposer.ProposeValue(context.Background(), sa, func() {
						memDBTx.Commit()
					})
        	 for _, c := range tx.changelist {
				s.queue.Publish(c)
			}
			if len(tx.changelist) != 0 {
				s.queue.Publish(state.EventCommit{})
			}
        watch目录：
        watch:"github.com/docker/go-events"  

**manager/orchestrator目录**

	    orchestrator开3个管理协层
	   m.replicatedOrchestrator = orchestrator.NewReplicatedOrchestrator(s)
	   m.globalOrchestrator = orchestrator.NewGlobalOrchestrator(s)
	   m.taskReaper = orchestrator.NewTaskReaper(s)
	  //以global 为例子
      //github.com\docker\swarmkit\manager\orchestrator\global.go
	  func (g *GlobalOrchestrator) Run(ctx context.Context)
	    	
		监听订阅相关queue event
		// Watch changes to services and tasks
		queue := g.store.WatchQueue()
		watcher, cancel := queue.Watch()
	
		//获取集群，节点，服务
		store.FindClusters(readTx, store.ByName("default"))
		store.FindNodes(readTx, store.All)
		store.FindServices(readTx, store.All)
	
		for {
			select {
			case event := <-watcher:
				  switch v := event.(type) {
				    case state.EventUpdateCluster:
				    case state.EventCreateService:
				    	... 
				  }
			}
		}

######docker 中cluster相关结构体：
docker 1.12 新的Cluster结构体嵌入swarmagent.Node的结构体，集成 swamkit ,通过调用swarmagent.NewNode(...)，是swamkit的使用者：

    //github.com\docker\docker\daemon\cluster\cluster.go		   
	// Cluster provides capabilities to participate in a cluster as a worker or a
	// manager.
		type Cluster struct {
			sync.RWMutex
			*node
			root        string
			config      Config
			configEvent chan struct{} // todo: make this array and goroutine safe
			listenAddr  string
			stop        bool
			err         error
			cancelDelay func()
		}

	type node struct {
		*swarmagent.Node    
		done           chan struct{}
		ready          bool
		conn           *grpc.ClientConn
		client         swarmapi.ControlClient
		reconnectDelay time.Duration
	}
	github.com\docker\swarmkit\agent\node.go:
	含有chan 参数 基本就是用来同步通信的。
	// Node implements the primary node functionality for a member of a swarm
	// cluster. Node handles workloads and may also run as a manager.
	type Node struct {
		sync.RWMutex
		config               *NodeConfig
		remotes              *persistentRemotes
		role                 string
		conn                 *grpc.ClientConn
		connCond             *sync.Cond
		nodeID               string
		nodeMembership       api.NodeSpec_Membership
		started              chan struct{}
		stopped              chan struct{}
		ready                chan struct{} // closed when agent has completed registration and manager(if enabled) is ready to receive control requests
		certificateRequested chan struct{} // closed when certificate issue request has been sent by node
		closed               chan struct{}
		err                  error
		agent                *Agent
		manager              *manager.Manager
		roleChangeReq        chan api.NodeRole // used to send role updates from the dispatcher api on promotion/demotion
		managerRoleCh        chan struct{}
	}

#####dockerd 启动时跟 cluster相关 流程：

	daemonCli = NewDaemonCli()
	err = daemonCli.start()
	func (cli *DaemonCli) start() (err error)
		 api := apiserver.New(serverConfig)
		 cli.api = api
		 registryService := registry.NewService(cli.Config.ServiceOptions)
		 containerdRemote, err := libcontainerd.New(cli.getLibcontainerdRoot(), cli.getPlatformRemoteOptions()...)
		 pluginInit(cli.Config, containerdRemote, registryService)

		//	"github.com/docker/docker/daemon"
	  //"github.com/docker/docker/daemon/cluster"
		 d, err := daemon.NewDaemon(cli.Config, registryService, containerdRemote)
		 c, err := cluster.New(cluster.Config{
				Root:    cli.Config.Root,
				Name:    name,
				Backend: d,
			})
		      st, err := c.loadState() 
			  //见下cluster startNewNode 流程
		      n, err := c.startNewNode(false, st.ListenAddr, "", "", "", false)
		 cli.initMiddlewares(api, serverConfig)
		 initRouter(api, d, c)
				routers := []router.Router{
					container.NewRouter(d, decoder),
					image.NewRouter(d, decoder),
					systemrouter.NewRouter(d, c),
					volume.NewRouter(d),
					build.NewRouter(dockerfile.NewBuildManager(d)),
					swarmrouter.NewRouter(c),
				}
		 go api.Wait(serveAPIWait)

####docker node ls 流程
	docker client:
	//github.com\docker\docker\api\client\node\list.go
     newListCommand(dockerCli *client.DockerCli) 
		runList(dockerCli, opts)
		  client := dockerCli.Client()
		  nodes, err := client.NodeList(,,)
		  //github.com\docker\engine-api\client\node_list.go
		     resp, err := cli.get(ctx, "/nodes", query, nil)
	
    dockerd daemon:
	//github.com\docker\docker\api\server\router\swarm\cluster.go
	router.NewGetRoute("/nodes", sr.getNodes),
	   nodes, err := sr.backend.GetNodes(basictypes.NodeListOptions{Filter: filter})
	   //src\github.com\docker\docker\daemon\cluster\cluster.go
	   func (c *Cluster) GetNodes(options apitypes.NodeListOptions) ([]types.Node, error) 
	          r, err := c.client.ListNodes(ctx,&swarmapi.ListNodesRequest{Filters: filters})
	        /*  
	          //ListNodes 接口：github.com\docker\swarmkit\api\control.proto
	          //该服务在上面分析 swamkit manager 开启流程时注册
	          // github.com\docker\swarmkit\manager\manager.go:465
	          //"github.com/docker/swarmkit/manager/controlapi"
	          baseControlAPI := controlapi.NewServer(m.RaftNode.MemoryStore(), m.RaftNode)
	          service Control {
	          	rpc ListNodes(ListNodesRequest) returns (ListNodesResponse) {
					option (docker.protobuf.plugin.tls_authorization) = { roles: "swarm-manager" };
		  	    };
			  }*/
			  // 实际函数：github.com\docker\swarmkit\manager\controlapi\node.go
			  func (s *Server) ListNodes(ctx context.Context, request *api.ListNodesRequest) (*api.ListNodesResponse, error
			  	s.store.View(func(tx store.ReadTx) {..})

// docker service create 流程

		//github.com\docker\docker\api\server\router\swarm\cluster.go
		   router.NewPostRoute("/services/create", sr.createService),
		   //github.com\docker\docker\daemon\cluster\cluster.go
		   func (c *Cluster) CreateService(s types.ServiceSpec, encodedAuth string) (string, error)
                       c.isActiveManager() 
                       populateNetworkID(ctx, c.client, &s)
                       // "github.com/docker/docker/daemon/cluster/convert"
					   //转化为grpc接口的结构体
                       serviceSpec, err := convert.ServiceSpecToGRPC(s)
                       //c.client : client         swarmapi.ControlClient
					   //调用swamkit grpc api 接口
					   //github.com\docker\swarmkit\manager\controlapi\service.go
                       r, err := c.client.CreateService(ctx, &swarmapi.CreateServiceRequest{Spec: &serviceSpec})
                       
					   // github.com\docker\docker\daemon\cluster\cluster.go
                        func (s *Server) CreateService(ctx context.Context, request *api.CreateServiceRequest)
    						//见上swamkit state分析 manager/state目录
                       		err := s.store.Update(func(tx store.Tx) error {
									return store.CreateService(tx, service)
							 })
        然后Orchestrator目录 会接收到event,来创建，异步，类似消息队列！
                       
	

####swarm init流程：
 - cli: 调用docker swarm init命令：
   - docker/api/client/swarm/init.go：
   		
			 runInit(dockerCli *client.DockerCli, flags *pflag.FlagSet, opts initOptions)
			   client := dockerCli.Client()
			   nodeID, err := client.SwarmInit(ctx, req)
                  serverResp, err := cli.post(ctx, "/swarm/init", nil, req, nil)
 - dockerd daemon server:
  - docker/api/server/router/swarm/cluster.go:
 
			 //router.NewPostRoute("/swarm/init", sr.initCluster)
				func (sr *swarmRouter) initCluster(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) 
				   nodeID, err := sr.backend.Init(req)
				      //docker/docker/daemon/cluster/cluster.go:
					  // Init initializes new cluster from user provided request.
					  // func (c *Cluster) Init(req types.InitRequest)
			
					  validateAndSanitizeInitRequest(&req)
					  //见下cluster startNewNode 流程
					  n, err := c.startNewNode(req.ForceNewCluster, req.ListenAddr, "", "", "", false)


####swarm join 流程:
docker swarm join --secret <SECRET> <MANAGER-IP>:<PORT>

       //github.com\docker\docker\api\server\router\swarm\cluster.go
		router.NewPostRoute("/swarm/join", sr.joinCluster),
		func (c *Cluster) Join(req types.JoinRequest)
		   //注意跟init 参数不同。看swam init
		   //见下cluster startNewNode 流程流程
		   n, err := c.startNewNode(false, req.ListenAddr, req.RemoteAddrs[0], req.Secret, req.CACertHash, req.Manager)
		   certificateRequested := n.CertificateRequested()
		   错误处理

####cluster startNewNode 流程：
可见dockerd,swam join/init最终都会调用startNewNode函数，主要参数不同：

	//github.com\docker\docker\daemon\cluster\cluster.go
	func (c *Cluster) startNewNode(forceNewCluster bool, localAddr, remoteAddr, listenAddr, advertiseAddr, joinAddr, joinToken string)
			c.config.Backend.IsSwarmCompatible()
			//swarmkit 代码
			//swarmagent "github.com/docker/swarmkit/agent"
			//github.com\docker\swarmkit\agent\node.go
			n, err := swarmagent.NewNode(&swarmagent.NodeConfig{...})
				//该Node为SWAMKIT定义的
				n := &Node{remotes: newPersistentRemotes(stateFile, p...)
			//该node 为 dockerd 定义的 github.com\docker\swarmkit\agent\node.go
			node := &node{	Node:  n,}
	    	c.node = node
			//见下 swamkit Node start 流程
			n.Start(ctx);
#### swamkit Node start 流程：
	//github.com\docker\swarmkit\agent\node.go
	n.Start(ctx):
	  //func (n *Node) run(ctx context.Context) 
		go n.run(ctx)
			    certificates TLS等安全接入 相关处理，有内存数据库，开携程处理等等
				db, err := bolt.Open(taskDBPath, 0666, nil)
				n.loadCertificates()
			wg.Add(2)
			go func() {
				//开启manager 节点
				managerErr = n.runManager(ctx, securityConfig, managerReady) // store err and loop
					remoteAddr, _ := n.remotes.Select(n.nodeID)
					    //"github.com/docker/swarmkit/manager"
						m, err := manager.New(&manager.Config{ForceNewCluster: n.config.ForceNewCluster,JoinRaft:  remoteAddr.Addr,})	
						   	    创建文件
						   	    RaftNode := raft.NewNode(context.TODO(), newNodeOpts)
					
								m := &Manager{
									config:      config,
									listeners:   listeners,
									caserver:    ca.NewServer(RaftNode.MemoryStore(), config.SecurityConfig),
									Dispatcher:  dispatcher.New(RaftNode, dispatcherConfig),
									server:      grpc.NewServer(opts...),
									localserver: grpc.NewServer(opts...),
									RaftNode:    RaftNode,
									stopped:     make(chan struct{}),
								}
								
						go func() {
							// Run starts all manager sub-systems and the gRPC server at the configured
							// address.
							//分析见下面
							m.Run(context.Background()) // todo: store error
										
						}()
						n.manager = m	
						...	
				wg.Done()
					
			}()
			//开启 work节点
			go func() {
				//github.com\docker\swarmkit\agent\node.go
				agentErr = n.runAgent(ctx, db, securityConfig.ClientTLSCreds, agentReady)
				      picker := picker.NewPicker(n.remotes, manager.Addr)
					   conn, err := grpc.Dial(manager.Addr,grpc.WithPicker(picker),..)
					   //结构体： github.com\docker\swarmkit\agent\agent.go
					   agent, err := New(&Config{
										Hostname:         n.config.Hostname,
										Managers:         n.remotes,
										Executor:         n.config.Executor,
										DB:               db,
										Conn:             conn,
										NotifyRoleChange: n.roleChangeReq,
									})
					   	//github.com\docker\swarmkit\agent\agent.go
					    agent.Start(ctx)
					    	err := a.worker.Init(ctx)
					    	a.worker.Listen(ctx, reporter)
				wg.Done()
		
			}()
	
			go func() {
				//agent 先起来
				<-agentReady
				if role == ca.ManagerRole {
					<-managerReady
				}
				close(n.ready)
			}()
			wg.Wait()
	
####swamkit manager 开启流程:
m.Run(context.Background()):

	   //github.com\docker\swarmkit\manager\manager.go
		func (m *Manager) Run(parent context.Context) 
			leadershipCh, cancel := m.RaftNode.SubscribeLeadership()
			go func()
			 if newState == raft.IsLeader 
			 	//开启功能goroutine
				s := m.RaftNode.MemoryStore()
				s.Update(func(tx store.Tx) error {
					store.CreateCluster(tx, &api.Cluster{...})
					// Add Node entry for ourself, if one doesn't exist already.
					store.CreateNode(tx, &api.Node{..})
				})
				m.replicatedOrchestrator = orchestrator.NewReplicatedOrchestrator(s)
				m.globalOrchestrator = orchestrator.NewGlobalOrchestrator(s)
				m.taskReaper = orchestrator.NewTaskReaper(s)
				m.scheduler = scheduler.New(s)
				m.keyManager = keymanager.New(m.RaftNode.MemoryStore(), keymanager.DefaultConfig()
			
				m.allocator, err = allocator.New(s)
			
				go func(keyManager *keymanager.KeyManager) {
					if err := keyManager.Run(ctx); err != nil {
						log.G(ctx).WithError(err).Error("keymanager failed with an error")
					}
				}(m.keyManager)
			
				go func(d *dispatcher.Dispatcher) {
					if err := d.Run(ctx); err != nil {
						log.G(ctx).WithError(err).Error("Dispatcher exited with an error")
					}
				}(m.Dispatcher)
			
				go func(server *ca.Server) {
					if err := server.Run(ctx); err != nil {
						log.G(ctx).WithError(err).Error("CA signer exited with an error")
					}
				}(m.caserver)
			else if newState == raft.IsFollower
					//停止功能goroutine
					m.Dispatcher.Stop()
					 m.caserver.Stop()
					 m.replicatedOrchestrator.Stop()
					 ...
			
			//grpc注册服务
			baseControlAPI := controlapi.NewServer(m.RaftNode.MemoryStore(), m.RaftNode, m.config.SecurityConfig.RootCA())
			healthServer := health.NewHealthServer()
			...
			authenticatedControlAPI := api.NewAuthenticatedWrapperControlServer(baseControlAPI, authorize)
			...
			
			// Set the raft server as serving for the health server
			healthServer.SetServingStatus("Raft", api.HealthCheckResponse_SERVING)
			m.RaftNode.JoinAndStart();
			go func() {
					err := m.RaftNode.Run(ctx)
			}
			err := raft.WaitForLeader(ctx, m.RaftNode); 
			c, err := raft.WaitForCluster(ctx, m.RaftNode）			

