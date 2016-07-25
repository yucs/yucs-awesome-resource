#跟docker 1.11.0变化：
   [Docker 1.12.0将要发布的新功能](http://liubin.org/blog/2016/06/17/whats-new-in-docker-1-dot-12-dot-0/)
#源码布局，大体流程没大的区别

主要分析 新加入 跟swamkit集成的部分。

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

							n, err := swarmagent.NewNode(...)
							    //github.com\docker\swarmkit\picker\picker.go封装"github.com/docker/swarmkit/api"
								n := &Node{remotes: newPersistentRemotes(stateFile, p...)}


							

							node := &node{	Node:  n,}
					    	c.node = node
							

							n.Start(ctx);
								//func (n *Node) run(ctx context.Context) 
  								go n.run(ctx)
  									    certificates TLS等安全接入 相关处理，有内存数据库，开携程处理等等
  										
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
														m.Run(context.Background()) // todo: store error
					
													}()
													n.manager = m	
													...	
											wg.Done()
												
										}()
										go func() {
											agentErr = n.runAgent(ctx, db, securityConfig.ClientTLSCreds, agentReady)
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
