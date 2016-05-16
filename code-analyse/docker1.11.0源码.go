//要是我，对于 每个模块，我会怎么写.
//看代码：自上而下的思考逻辑， 跟而自己写代码，则是自下而上的思考.


//感觉自己造轮子的pkg 比较多。。
/*
  源码分析部分：
     command api 流程 
     docker  daemon 启动流程
     创建container 流程
     docker event 逻辑
     volume plugin
*/

启动 docker/docker：
clientCli := client.NewDockerCli(stdin, stdout, stderr, clientFlags)//api/client/cli.go
c := cli.New(clientCli, daemonCli) // docker/daemon.go : daemonCli cli.Handler = NewDaemonCli()
 c.Run(flag.Args()...);   
 
 //command api 流程:
  cli/cli.go -> api/client : CmdXXX -> engine-api (SDK 包装API)->  api/server -> s.backend.func ->daemon struct 动作(deamon目录)。 

//daemon 启动流程：
  docker/daemon.go : 
 func (cli *DaemonCli) CmdDaemon(args ...string)

      api, err := apiserver.New(serverConfig)
            r.runContainerdDaemon() //会exec.command():docker-containerd 
            r.apiClient = containerd.NewAPIClient(conn)
     	//遍历开服务
       srv, err := s.newServer(addr.Proto, addr.Addr)

       containerdRemote, err := libcontainerd.New(cli.getLibcontainerdRoot(), cli.getPlatformRemoteOptions()...)
  
      // NewDaemon sets up everything for the daemon to be able to service
      // requests from the webserver. 
	    // daemon/daemon.go: 
       d, err := daemon.NewDaemon(cli.Config, registryService,containerdRemote)
                          // Configure the volumes driver
                           volStore, err := configureVolumes(config, rootUID, rootGID)
                            
                            // Discovery is only enabled when the daemon is launched with an address to advertise.  When
                          // initialized, the daemon is registered and we can store the discovery backend as its read-only
                           err := d.initDiscovery(config)
                           //netork
                           d.netController, err = d.initNetworkController(config)

                           //  libcontainerd/remote_linux.go
                          d.containerd, err = containerdRemote.Client(d)   
   
     api.InitRouters(d) //router/server.go

     setupConfigReloadTrap(*configFile, cli.flags, reload)

     // The serve API routine never exits unless an error occurs
	// We need to start it as a goroutine and wait on it so
	// daemon doesn't exit
	   serveAPIWait := make(chan error)
	   go api.Wait(serveAPIWait) 
	     s.serveAPI()
	        s.initRouterSwapper()


//该结构体 router_swapper.go 实现http.Handler       
func (s *Server) initRouterSwapper() {
	s.routerSwapper = &routerSwapper{
		router: s.createMux(),
	}
}
// createMux initializes the main router the server uses.
func (s *Server) createMux() *mux.Router {
	m := mux.NewRouter()
	if utils.IsDebugEnabled() {
		profilerSetup(m, "/debug/")
	}

	logrus.Debugf("Registering routers")
	for _, apiRouter := range s.routers {
		for _, r := range apiRouter.Routes() {
			f := s.makeHTTPHandler(r.Handler())

			logrus.Debugf("Registering %s, %s", r.Method(), r.Path())
			m.Path(versionMatcher + r.Path()).Methods(r.Method()).Handler(f)
			m.Path(r.Path()).Methods(r.Method()).Handler(f)
		}
	}

	return m
}

api/server/router/container 都会注册：
type containerRouter struct {
	backend Backend
	routes  []router.Route
}

// Backend is all the methods that need to be implemented to provide container specific functionality.
type Backend interface {
	execBackend
	copyBackend
	stateBackend
	monitorBackend
	attachBackend
}
Daemon目录下的 daemon.go:Daemon struct实现了Backend接口。
其他network,graph 一样。 


//with containerd /runc 



 
 //docker event 逻辑: 
  docker daemon 端：
   订阅模式 Publisher/subscribers : "github.com/docker/docker/pkg/pubsub" :  sync.WaitGroup + goroute+timeout 方式 ,可以借鉴

 daemon/daemon: 
  func NewDaemon(config *Config, registryService *registry.Service) 
      eventsService := events.New()

/daemon/events: 
 func New() *Events {
	return &Events{
		events: make([]eventtypes.Message, 0, eventsLimit),
		pub:    pubsub.NewPublisher(100*time.Millisecond, bufferSize),
	}
  }

 server api  注册 ：

 api/server/router/system.go ：router.NewGetRoute("/events", r.getEvents),

func (s *systemRouter) getEvents(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
    ....
	output := ioutils.NewWriteFlusher(w)
	defer output.Close()
	output.Flush()

	enc := json.NewEncoder(output)

	buffered, l := s.backend.SubscribeToEvents(since, sinceNano, ef)
	defer s.backend.UnsubscribeFromEvents(l)
     ... 

	var closeNotify <-chan bool
	if closeNotifier, ok := w.(http.CloseNotifier); ok {
		closeNotify = closeNotifier.CloseNotify()
	}

	for {
		select {
		case ev := <-l:
			jev, ok := ev.(events.Message)
			if !ok {
				logrus.Warnf("unexpected event message: %q", ev)
				continue
			}
			if err := enc.Encode(jev); err != nil {
				return err
			}
		case <-timer.C:
			return nil
		case <-closeNotify:
			logrus.Debug("Client disconnected, stop sending events")
			return nil
		}
	}
}

docker cli 端：
   api/events.go: 
   func (cli *DockerCli) CmdEvents(args ...string) error 
      responseBody, err := cli.client.Events(context.Background(), options)
 	  defer responseBody.Close()
      streamEvents(responseBody, cli.out)
        decodeEvents(): 死函数 for{} GET 请求

        结构体咋跟 官网文档 不一样？？（curl  0.0.0.0:2375/events?since=1455606562）










   //volume plugin:

type Daemon struct {
      volumes                   *store.VolumeStore
}

//daemon开启：
//需要时获取 docker/daemon.go: ??



 //docker/daemon.go 启动时plugin 初始化：
 func (cli *DaemonCli) CmdDaemon(args ...string) 
      daemon.NewDaemon(cli.Config, registryService)
               configureVolumes()

 func configureVolumes(config *Config, rootUID, rootGID int) (*store.VolumeStore, error) {
  volumesDriver, err := local.New(config.Root, rootUID, rootGID)
  if err != nil {
    return nil, err
  }
 //"github.com/docker/docker/volume/drivers": extpoint.go
  volumedrivers.Register(volumesDriver, volumesDriver.Name())
  return store.New(), nil
}
    

// docker run --rm -ti --volume-driver lvs  -v DBASS_DAT:/mnd  ubuntu bin/bash 
 //create api  plugin相关:
    // hostConfig :&{[DBASS_DAT:/mnd]
     func (daemon *Daemon) create(params types.ContainerCreateConfig)
           
            daemon.setHostConfig(container, params.HostConfig); 
                     daemon.registerMountPoints(container, hostConfig);
                          //bind:&{ /mnd true DBASS_DAT lvs <nil>  rprivate false}
                         //volume/store/store.go
                          v, err := daemon.volumes.CreateWithRef(bind.Name, bind.Driver, container.ID, nil)  
   								v, exists := s.getNamed(name); exists
   								//volume/drivers/extpoint.go
   								vd, err := volumedrivers.GetDriver(driverName)
   											Lookup(name) 
   											   //github.com/docker/docker/pkg/plugins
   											   pl, err := plugins.Get(name, extName)
   											   d := NewVolumeDriver(name, pl.Client)

   								if v, _ := vd.Get(name); v != nil {
										return v, nil
								}
								
								return vd.Create(name, opts)


按需连接  ：
docker目录下：  grep 'github.com/docker/docker/pkg/plugins' * -nr。就几个地方。
以postVolumesCreate为例： // 就是做成一次次的http请求.
   daemon.volumes.Create(name, driverName, opts) //(s *VolumeStore) Create(name, driverName string, opts map[string]string) (volume.Volume, error)
            volumedrivers.GetDriver(driverName) //volumeDriver 为plugin接口
                Lookup(name)    // Lookup(name string) (volume.Driver, error)
                  plugins.Get(name, extName)
                  NewVolumeDriver(name, pl.Client)
            vd.Create(name, opts)
}



//创建container：
//docker 1.10.X
api/server/router/container.go：
func (r *containerRouter) initRoutes()
  ...
  router.NewPostRoute("/containers/create", r.postContainersCreate), 
  ...
 

 func (s *containerRouter) postContainersCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error
  
  //也是用json.NewDecoder(src) ,Be aware this function is not checking whether the resulted structs are nil,
  // it's your business to do so
  //传输也是用包裹ContainerConfigWrapper 复合结构体
 
  config, hostConfig, networkingConfig, err := runconfig.DecodeContainerConfig(r.Body) 
      //函数中 var w ContainerConfigWrapper  //由engine-api/types/container中定义，验证参数等等
      
      //engine-api/types/container hostConfig说明： HostConfig the non-portable Config structure of a container.
      // Here, "non-portable" means "dependent of the host we are running on".
      // Portable information *should* appear in Config.
      API接口的Mounts 结构体在哪啊？！

  //两个地方config包装:client ->(ContainerConfigWrapper) server ->(ContainerCreateConfig) dameon 
  ccr, err := s.backend.ContainerCreate(types.ContainerCreateConfig{
    Name:             name,
    Config:           config,
    HostConfig:       hostConfig,
    NetworkingConfig: networkingConfig,
    AdjustCPUShares:  adjustCPUShares,
  })

//daemon/create.go
func (daemon *Daemon) ContainerCreate(params types.ContainerCreateConfig) (types.ContainerCreateResponse, error)
    验证参数
     // verifyContainerSettings performs validation of the hostconfig and config
     // structures.
     daemon.verifyContainerSettings
     // Checks if the client set configurations for more than one network while creating a container
     daemon.verifyNetworkingConfig

     // adaptContainerSettings is called during container creation to modify any
     // settings necessary in the HostConfig structure.
     daemon.adaptContainerSettings

  // Create creates a new container from the given configuration with a given name.
     daemon.create(params)
  

  // Create creates a new container from the given configuration with a given name.
func (daemon *Daemon) create(params types.ContainerCreateConfig) (retC *container.Container, retErr error)
   
   //有些复杂，涉及到其他，volume...ipc,dist 其他东西
   .....
   创建失败，资源回收等出差处理 通过 defer func 判断retErr 方式。分配完一个资源，一个defer.本生也是栈方式。
   .....

   创建完：
   daemon.LogContainerEvent(container, "create")




//ContainerStart
//docker 1.11  with containerd runc
func (daemon *Daemon) ContainerStart
         daemon.containerStart(container)
            daemon.createSpec(container)
            daemon.containerd.Create()   //  client obj  func (clnt *client) Create  libcontainerd/client_linux.go
                 container := clnt.newContainer(filepath.Join(dir, containerID), options...) &container{} obj //包含 process
                 container.start()  //libcontainerd/container_linux.go
                   resp, err := ctr.client.remote.apiClient.CreateContainer(context.Background(), r)




/*
docker 1.10.3
 // containerStart prepares the container to run by setting up everything the
// container needs, such as storage and networking, as well as links
// between containers. The container is left waiting for a signal to
// begin running.
      func (daemon *Daemon) ContainerStart(name string, hostConfig *containertypes.HostConfig) error 
	         
	        
	         // defer func() {
	         //  	 daemon.Cleanup(container)
	         //  	 daemon.LogContainerEventWithAttributes(container, "die", attributes)
	         //  	      container.UnmountVolumes(false, daemon.LogVolumeEvent)
	         //        	 volumeMount.Volume.Unmount()
	      			//      volumeEventLog(volumeMount.Volume.Name(), "unmount", attributes)
          //    }()	

      		mounts, err := daemon.setupMounts(container)
      				for _, m := range container.MountPoints {
 						daemon.lazyInitializeVolume(container.ID, m)
      							daemon.volumes.GetWithRef(m.Name, m.Driver, containerID)
      					path, err := m.Setup()
      					       if m.Volume != nil {
							 	return m.Volume.Mount()
							   }
      					daemon.LogVolumeEvent(m.Volume.Name(), "mount", attributes)
      				}


      		 err := daemon.waitForStart(container);
      		     return container.StartMonitor(daemon, container.HostConfig.RestartPolicy)
  							return container.monitor.wait() 
  							   	select {
									case <-m.startSignal:
										//"github.com/docker/docker/pkg/promise"
										//m.start :container/monitor.go
									case err := <-promise.Go(m.start):
										return err
	
								  }

*/

// //container/monitor.go
//   func (m *containerMonitor) start() error
  		
//   		defer func() {
//   			//收到 containerStop 后执行。
//   			m.Close()
//   			   m.supervisor.Cleanup(m.container)
//   			}()

//   		for{

//   		}



//  func (daemon *Daemon) containerStop(container *container.Container, seconds int)
//          发信号。





 