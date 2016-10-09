https://docs.docker.com/engine/extend/plugins_network/
http://blog.csdn.net/xiaolunsanguo/article/details/52036418
https://github.com/docker/libnetwork/blob/master/docs/design.md

//github.com\docker\docker\cmd\dockerd\daemon.go
func (cli *DaemonCli) start() 
    //github.com\docker\docker\cmd\dockerd\daemon.go
	d, err := daemon.NewDaemon(cli.Config, registryService, containerdRemote)
	   //github.com\docker\docker\daemon\daemon.go
	   err := d.restore()
	        //github.com\docker\docker\daemon\daemon_unix.go
	   		daemon.netController, err = daemon.initNetworkController(daemon.configStore, activeSandboxes)
	   		   netOptions, err := daemon.networkOptions(config, activeSandboxes)
	   		   //github.com\docker\libnetwork\controller.go
	   		   controller, err := libnetwork.New(netOptions...)
	   		   // Initialize default network on "null"
	   		   // Initialize default network on "host"
	   		   	//controller.NetworkByName("none"): 判断是否已有“none”
	   		    //github.com\docker\libnetwork\controller.go
	   			// Initialize default network on "host"
				if n, _ := controller.NetworkByName("host"); n == nil {
					if _, err := controller.NewNetwork("host", "host", "", libnetwork.NetworkOptionPersist(true)); err != nil {
						return nil, fmt.Errorf("Error creating default \"host\" network: %v", err)
					}
				}
				if !config.DisableBridge {
					// Initialize default driver "bridge"
					if err := initBridgeDriver(controller, config); err != nil {
						return nil, err
					}
				}
			//重启容器时 调用daemon.waitForNetworks(c)
				for c, notifier := range restartContainers {
					group.Add(1)
					go func(c *container.Container, chNotify chan struct{}) {
						defer group.Done()
						...
						// Make sure networks are available before starting
						// waitForNetworks is used during daemon initialization when starting up containers
						// It ensures that all of a container's networks are available before the daemon tries to start the container.
						// In practice it just makes sure the discovery service is available for containers which use a network that require discovery.

						daemon.waitForNetworks(c)
					}(c, notifier)
				}
				group.Wait()


//NetworkController
// New creates a new instance of network controller.
//github.com\docker\libnetwork\controller.go
func New(cfgOptions ...config.Option) (NetworkController, error)
	c := &controller{
		id:              stringid.GenerateRandomID(),

		cfg:             config.ParseConfigOptions(cfgOptions...),
		sandboxes:       sandboxTable{},
		svcRecords:      make(map[string]svcInfo),
		serviceBindings: make(map[serviceKey]*service),
		agentInitDone:   make(chan struct{}),
	}
    //初始化后端存储
    //github.com\docker\libnetwork\store.go
	//默认是store 是 boltdb
	//github.com\docker\libnetwork\config\config.go
	//var defaultScopes = makeDefaultScopes()
   c.initStores();
   //github.com\docker\libnetwork\drvregistry\drvregistry.go
   // DrvRegistry holds the registry of all network drivers and IPAM drivers that it knows about.
   drvRegistry, err := drvregistry.New(c.getStore(datastore.LocalScope), c.getStore(datastore.GlobalScope), c.RegisterDriver, nil)
   //遍历"bridge"，"host"，"macvlan"， "remote"，"overlay"，"null"注册
   //github.com\docker\libnetwork\drivers_linux.go
   for _, i := range getInitializers() {
   		/*
			func getInitializers() []initializer {
				in := []initializer{
					{bridge.Init, "bridge"},
					{host.Init, "host"},
					{macvlan.Init, "macvlan"},
					{null.Init, "null"},
					 //"github.com/docker/libnetwork/drivers/remote"
					{remote.Init, "remote"},
					{overlay.Init, "overlay"},
				}

				in = append(in, additionalDrivers()...)
				return in
			}
*/
   	drvRegistry.AddDriver(i.ntype, i.fn, dcfg)
   }
   //github.com\docker\libnetwork\drivers_ipam.go
   //初始化IPAM IP地址管理
   initIPAMDrivers(drvRegistry, nil, c.getStore(datastore.GlobalScope))
   		//github.com/docker/libnetwork/ipams 目录下
   		//192.168.XX.XX 等私有IP
   		builtinIpam.Init()
   		//向PLugin 发getCapabilities 接口,然后向store 注册
   		//Init registers a remote ipam when its plugin is activated
   		//github.com\docker\libnetwork\ipams\remote\remote.go
   		remoteIpam.Init()
   		nullIpam.Init()

   c.initDiscovery(c.cfg.Cluster.Watcher)
   c.startExternalKeyListener()

 remote 的 IPAMDriver： 
//github.com\docker\libnetwork\ipams\remote\remote.go
//Init registers a remote ipam when its plugin is activated
//同样分配address
func Init(cb ipamapi.Callback, l, g interface{}) error {
	//github.com\docker\docker\pkg\plugins\plugins.go
	//   extpointHandlers[iface] = fn 还没调用
	plugins.Handle(ipamapi.PluginEndpointType, func(name string, client *plugins.Client) {
		a := newAllocator(name, client)
		if cps, err := a.(*allocator).getCapabilities(); err == nil {
			if err := cb.RegisterIpamDriverWithCapabilities(name, a, cps); err != nil {
				log.Errorf("error registering remote ipam driver %s due to %v", name, err)
			}
		} else {
			log.Infof("remote ipam driver %s does not support capabilities", name)
			log.Debug(err)
			if err := cb.RegisterIpamDriver(name, a); err != nil {
				log.Errorf("error registering remote ipam driver %s due to %v", name, err)
			}
		}
	})
	return nil
}


remote 的 driver:
//github.com\docker\libnetwork\drivers\remote\driver.go
// Init makes sure a remote driver is registered when a network driver
// plugin is activated.
//同样其他接口 调用PLUGIN 接口
func Init(dc driverapi.DriverCallback, config map[string]interface{}) error {
	plugins.Handle(driverapi.NetworkPluginEndpointType, func(name string, client *plugins.Client) {
		// negotiate driver capability with client
		d := newDriver(name, client)
		c, err := d.(*driver).getCapabilities()
		if err != nil {
			log.Errorf("error getting capability for %s due to %v", name, err)
			return
		}
		if err = dc.RegisterDriver(name, d, *c); err != nil {
			log.Errorf("error registering driver for %s due to %v", name, err)
		}
	})
	return nil
}




func (c *controller) NewNetwork(networkType, name string, id string, options ...NetworkOption) (Network, error) 
  	network := &network{
		name:        name,
		networkType: networkType,
		generic:     map[string]interface{}{netlabel.GenericData: make(map[string]string)},
		ipamType:    ipamapi.DefaultIPAM,
		id:          id,
		ctrlr:       c,
		persist:     true,
		drvOnce:     &sync.Once{},
	}
	_, cap, err := network.resolveDriver(networkType, true)
	    c := n.getController()
	    err = c.loadDriver(name)
	       // Plugins pkg performs lazy loading of plugins that acts as remote drivers.
		   // As per the design, this Get call will result in remote driver discovery if there is a corresponding plugin available.
			_, err := plugins.Get(networkType, driverapi.NetworkPluginEndpointType)
	    d, cap = c.drvRegistry.Driver(name)
	// Make sure we have a driver available for this network type
	// before we allocate anything.
	_, err := network.driver(true);
	err = network.ipamAllocate()
	   n.getController().getIPAMDriver(n.ipamType)
	      		// Might be a plugin name. Try loading it
				if err := c.loadIPAMDriver(name)
				    _, err := plugins.Get(name, ipamapi.PluginEndpointType)
	   	err = n.ipamAllocateVersion(4, ipam)
	   	   d.PoolID, d.Pool, d.Meta, err = n.requestPoolHelper(ipam, n.addrSpace, cfg.PreferredPool, cfg.SubPool, n.ipamOptions, ipVer == 6)
	   	       poolID, pool, meta, err := ipam.RequestPool(addressSpace, preferredPool, subPool, options, v6)
	   	       	 //remote:github.com\docker\swarmkit\vendor\github.com\docker\libnetwork\ipams\remote\remote.go
	   	         a.call("RequestPool", req, res)
	err = c.addNetwork(network)
	...
    joinCluster(network)


docker network create --driver weave mynet
//github.com\docker\docker\api\server\router\network\network_routes.go
func (n *networkRouter) postNetworkCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string)
 
 n.clusterProvider.GetNetwork(create.Name)

 nw, err := n.backend.CreateNetwork(create)
 
//github.com\docker\docker\daemon\network.go
 func (daemon *Daemon) createNetwork(create types.NetworkCreateRequest, id string, agent bool) 
    c := daemon.netController
 		//参数
     nwOptions := []libnetwork.NetworkOption{
		libnetwork.NetworkOptionIpam(ipam.Driver, "", v4Conf, v6Conf, ipam.Options),
		libnetwork.NetworkOptionEnableIPv6(create.EnableIPv6),
		libnetwork.NetworkOptionDriverOpts(create.Options),
		libnetwork.NetworkOptionLabels(create.Labels),
	}

	//libnetwork 
	//github.com\docker\libnetwork\controller.go
	n, err := c.NewNetwork(driver, create.Name, id, nwOptions...)


docker run --network=mynet busybox top
//github.com\docker\docker\api\server\router\container\container.go
//创建
router.NewPostRoute("/containers/create", r.postContainersCreate),
func (s *containerRouter) postContainersCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string)
  ccr, err := s.backend.ContainerCreate(types.ContainerCreateConfig{..}, validateHostname)
     daemon.verifyNetworkingConfig(params.NetworkingConfig)
     container, err := daemon.create(params, managed)
        container, err = daemon.newContainer(params.Name, params.Config, imgID, managed)
            //VOLUME 相关
            daemon.createContainerPlatformSpecificSettings(container, params.Config, params.HostConfig)
            //Network
            daemon.updateContainerNetworkSettings(container, endpointsConfigs)

     container.ToDisk()
//启动
//github.com\docker\docker\daemon\start.go
func (daemon *Daemon) ContainerStart(name string, hostConfig *containertypes.HostConfig, validateHostname bool)
   
	// containerStart prepares the container to run by setting up everything the
	// container needs, such as storage and networking, as well as links
	// between containers. The container is left waiting for a signal to
	// begin running.
   daemon.containerStart(container)
      daemon.conditionalMountOnStart(container);
      //github.com\docker\docker\daemon\container_operations.go
      daemon.initializeNetworking(container);
         daemon.allocateNetwork(container)
               //// Cleanup any stale sandbox left over due to ungraceful daemon shutdown
               controller.SandboxDestroy(container.ID);
               daemon.connectToNetwork(container, n, nConf, updateSettings)
						sb := daemon.getNetworkSandbox(container)
						ep, err := n.CreateEndpoint(endpointName, createOptions...)
						err := ep.Join(sb, joinOptions...)
      spec, err := daemon.createSpec(container)
      daemon.containerd.Create(container.ID, *spec, createOptions...)


