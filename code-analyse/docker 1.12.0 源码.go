#跟docker 1.11.0变化：
   [Docker 1.12.0将要发布的新功能](http://liubin.org/blog/2016/06/17/whats-new-in-docker-1-dot-12-dot-0/)
#源码布局，大体流程没大的区别

主要分析 新加入 跟swamkit集成的部分。

//dockerd 区别
//github.com\docker\docker\cmd\dockerd\docker.go 
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
							//swarmkit 客户端？

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
											<-agentReady
											if role == ca.ManagerRole {
												<-managerReady
											}
											close(n.ready)
										}()
										wg.Wait()



  							c.saveState()
							c.config.Backend.SetClusterProvider(c)

					 		开启3个携程：	
					 		return node, nil
					  select {
						case <-n.Ready():
						case <-n.done:
					   } 

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



//docker swarm join --secret <SECRET> <MANAGER-IP>:<PORT>
       //github.com\docker\docker\api\server\router\swarm\cluster.go
		router.NewPostRoute("/swarm/join", sr.joinCluster),
		func (c *Cluster) Join(req types.JoinRequest)
		   //注意跟init 参数不同。看swam init
		   n, err := c.startNewNode(false, req.ListenAddr, req.RemoteAddrs[0], req.Secret, req.CACertHash, req.Manager)
		   certificateRequested := n.CertificateRequested()
		   错误处理

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
                       			//对表的操作？ 实际动作？
                        	//update 封装的，会对raft进行操作,其实就是raft的内存状态吧？
                       		err := s.store.Update(func(tx store.Tx) error {
									return store.CreateService(tx, service)
							 })
    

