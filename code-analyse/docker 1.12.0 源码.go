
http://www.containertutorials.com/swarmkit/architecture.html
https://blog.docker.com/2016/07/docker-built-in-orchestration-ready-for-production-docker-1-12-goes-ga/
//dockerd 区别
//github.com\docker\docker\cmd\dockerd\docker.go 

#跟docker 1.11.0变化：
   [Docker 1.12.0将要发布的新功能](http://liubin.org/blog/2016/06/17/whats-new-in-docker-1-dot-12-dot-0/)
#源码布局，大体流程没大的区别

主要分析 集成到docker1.12 的swarmkit。


//swarm init
cli: 调用docker swarm init命令：
 //api/client/swarm/init.go:
 func runInit(dockerCli *client.DockerCli, flags *pflag.FlagSet, opts initOptions) error {
client := dockerCli.Client()
nodeID, err := client.SwarmInit(ctx, req)
               serverResp, err := cli.post(ctx, "/swarm/init", nil, req, nil)
}	

server:
//api/server/router/swarm/cluster.go:
	 //router.NewPostRoute("/swarm/init", sr.initCluster)
				func (sr *swarmRouter) initCluster(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) 
				   nodeID, err := sr.backend.Init(req)
				      //docker/docker/daemon/cluster/cluster.go:
					  // Init initializes new cluster from user provided request.
					  // func (c *Cluster) Init(req types.InitRequest)
			
					  validateAndSanitizeInitRequest(&req)
					 
					  n, err := c.startNewNode(req.ForceNewCluster, req.ListenAddr, "", "", "", false)
							c.config.Backend.IsSwarmCompatible()
							//swarmkit 代码
							//swarmagent "github.com/docker/swarmkit/agent"
							//github.com\docker\swarmkit\agent\node.go

							n, err := swarmagent.NewNode(&swarmagent.NodeConfig{...})
								n := &Node{remotes: newPersistentRemotes(stateFile, p...)}


							

							node := &node{	Node:  n,}
					    	c.node = node
							

							n.Start(ctx);
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
										go func() {
											//\github.com\docker\swarmkit\agent\node.go
											agentErr = n.runAgent(ctx, db, securityConfig.ClientTLSCreds, agentReady)
											      picker := picker.NewPicker(n.remotes, manager.Addr)
												   conn, err := grpc.Dial(manager.Addr,grpc.WithPicker(picker),..)
												   //github.com\docker\swarmkit\agent\agent.go
													// Agent implements the primary node functionality for a member of a swarm
													// cluster. The primary functionality id to run and report on the status of
													// tasks assigned to the node.
												   type Agent struct {
															config *Config

															// The latest node object state from manager
															// for this node known to the agent.
															node *api.Node

															keys []*api.EncryptionKey

															sessionq chan sessionOperation
															worker   Worker

															started chan struct{}
															ready   chan struct{}
															stopped chan struct{} // requests shutdown
															closed  chan struct{} // only closed in run
															err     error         // read only after closed is closed
														}

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



  							c.saveState()
							c.config.Backend.SetClusterProvider(c)

				
//docker node ls
//github.com\docker\docker\api\client\node\list.go
cli: newListCommand(dockerCli *client.DockerCli) 
	runList(dockerCli, opts)
	  client := dockerCli.Client()
	  nodes, err := client.NodeList(,,)
	  //github.com\docker\engine-api\client\node_list.go
	     resp, err := cli.get(ctx, "/nodes", query, nil)
server:
//github.com\docker\docker\api\server\router\swarm\cluster.go
router.NewGetRoute("/nodes", sr.getNodes),
   nodes, err := sr.backend.GetNodes(basictypes.NodeListOptions{Filter: filter})
   //src\github.com\docker\docker\daemon\cluster\cluster.go
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
   func (c *Cluster) GetNodes(options apitypes.NodeListOptions) ([]types.Node, error) 
          r, err := c.client.ListNodes(ctx,&swarmapi.ListNodesRequest{Filters: filters})
          //github.com\docker\swarmkit\api\control.proto
          //该服务在上面分析startNewNode 时swamkit runManager 时注册
          // github.com\docker\swarmkit\manager\manager.go:465
        //  "github.com/docker/swarmkit/manager/controlapi"
          // 初步判断： Server.store就是raft 内存上的，update等操作，都对这进行封装了:
         		//github.com\docker\swarmkit\manager\state\raft\raft.go
                   //n.memoryStore = store.NewMemoryStore(n)
                   // func NewMemoryStore(proposer state.Proposer) *MemoryStore
          		  //
          baseControlAPI := controlapi.NewServer(m.RaftNode.MemoryStore(), m.RaftNode)


			          	// Server is the Cluster API gRPC server.
          				//github.com\docker\swarmkit\manager\controlapi\server.go

							type Server struct {
								store *store.MemoryStore
								raft  *raft.Node
							}

          service Control {
          	rpc ListNodes(ListNodesRequest) returns (ListNodesResponse) {
				option (docker.protobuf.plugin.tls_authorization) = { roles: "swarm-manager" };
	  	    };
		  }

		  // 实际函数：github.com\docker\swarmkit\manager\controlapi\node.go
		  func (s *Server) ListNodes(ctx context.Context, request *api.ListNodesRequest) (*api.ListNodesResponse, error
		  	//数据库呀？
		  	s.store.View(func(tx store.ReadTx) {..})




// docker service create --replicas 1 --name helloworld alpine ping docker.com
		//github.com\docker\docker\api\server\router\swarm\cluster.go
		   router.NewPostRoute("/services/create", sr.createService),
		   //github.com\docker\docker\daemon\cluster\cluster.go
		   func (c *Cluster) CreateService(s types.ServiceSpec, encodedAuth string) (string, error)
                       c.isActiveManager() 
                       populateNetworkID(ctx, c.client, &s)
                       // "github.com/docker/docker/daemon/cluster/convert"
                       serviceSpec, err := convert.ServiceSpecToGRPC(s)
                       //c.client : client         swarmapi.ControlClient
                       r, err := c.client.CreateService(ctx, &swarmapi.CreateServiceRequest{Spec: &serviceSpec})
                       // github.com\docker\docker\daemon\cluster\cluster.go
                        func (s *Server) CreateService(ctx context.Context, request *api.CreateServiceRequest)
    						//见下 swamkit state分析

                       		err := s.store.Update(func(tx store.Tx) error {
									return store.CreateService(tx, service)
							 })
                       		然后GlobalOrchestrator会接收到event,来创建，异步，类似消息队列！
                             //(g *GlobalOrchestrator) Run(ctx context.Context)

//github.com\docker\swarmkit\manager\manager.go
func (m *Manager) Run(parent context.Context) 
leadershipCh, cancel := m.RaftNode.SubscribeLeadership()
go func()
 if newState == raft.IsLeader 
 	//开启服务
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
		//停止服务
		m.Dispatcher.Stop()
		 m.caserver.Stop()
		 m.replicatedOrchestrator.Stop()
		 ...


baseControlAPI := controlapi.NewServer(m.RaftNode.MemoryStore(), m.RaftNode, m.config.SecurityConfig.RootCA())
healthServer := health.NewHealthServer()
...


// Set the raft server as serving for the health server
healthServer.SetServingStatus("Raft", api.HealthCheckResponse_SERVING)
m.RaftNode.JoinAndStart();
go func() {
		err := m.RaftNode.Run(ctx)
		}
err := raft.WaitForLeader(ctx, m.RaftNode); 
c, err := raft.WaitForCluster(ctx, m.RaftNode）			


//state
	   raft目录：raft协议  实现proposer.go:Proposer 接口
	   store目录：基于go-memdb内存，各个文件一个表，init()初始化表；
	  		// MemoryStore is a concurrency-safe, in-memory implementation of the Store interface.
				type MemoryStore struct 封装create,View,Update等操作，
       func (s *MemoryStore) update(proposer state.Proposer, cb func(Tx) error)
       		
		    var tx tx
			tx.init(memDBTx, curVersion)

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

//orchestrator
    开3个
   	m.replicatedOrchestrator = orchestrator.NewReplicatedOrchestrator(s)
	m.globalOrchestrator = orchestrator.NewGlobalOrchestrator(s)
	m.taskReaper = orchestrator.NewTaskReaper(s)

  //github.com\docker\swarmkit\manager\orchestrator\global.go
  func (g *GlobalOrchestrator) Run(ctx context.Context)
    	
	监听queue
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