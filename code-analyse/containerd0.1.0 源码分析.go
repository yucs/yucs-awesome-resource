func daemon(address, stateDir string, concurrency int, runtimeName string, runtimeArgs []string) error 
    sv, err := supervisor.New(stateDir, runtimeName, runtimeArgs)  
                     go s.exitHandler()
					 go s.oomHandler()
     machine, err := CollectMachineInformation()
     monitor, err := NewMonitor()
                    go m.start()

     // concurrency 框架
     wg := &sync.WaitGroup{}
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		w := supervisor.NewWorker(sv, wg)
		go w.Start()
	}
	if err := sv.Start(); err != nil {
		      //一个goroutine 见 func (s *Supervisor) Start() 
		return err
	}

	server, err := startServer(address, sv) 
	          s := grpc.NewServer()
			  types.RegisterAPIServer(s, server.NewServer(sv))


			


startTasks 跟 Task 区别：
task 加入job chanl 队列 ,startTasks执行队列
//一个worker 
	func (w *worker) Start() {
	 defer w.wg.Done()
	 //startTasks := make(chan *startTask, 10) chan 会死循环在这。
	 for t := range w.s.startTasks {
	 	  //runtime\container_linux.go
	 	  process, err := t.Container.Start(t.Checkpoint, runtime.NewStdio(t.Stdin, t.Stdout, t.Stderr))
	 	  err := w.s.monitor.MonitorOOM(t.Container); err != nil && err != runtime.ErrContainerExited
	 	  err := w.s.monitorProcess(process); err != nil 
	 	  w.s.notifySubscribers(Event{..})
	 }
	}

//
	func (s *Supervisor) Start() error {
	logrus.WithFields(logrus.Fields{
		"stateDir":    s.stateDir,
		"runtime":     s.runtime,
		"runtimeArgs": s.runtimeArgs,
		"memory":      s.machine.Memory,
		"cpus":        s.machine.Cpus,
	}).Debug("containerd: supervisor running")
	go func() {
		for i := range s.tasks {
			s.handleTask(i)  // (s *Supervisor) handleTask(i Task) 
		}
	}()
	return nil
}

//*****supervisor_linux.go
func (s *Supervisor) handleTask(i Task) {
	var err error
	switch t := i.(type) {
	case *AddProcessTask:
		err = s.addProcess(t)
	case *CreateCheckpointTask:
		err = s.createCheckpoint(t)
	case *DeleteCheckpointTask:
		err = s.deleteCheckpoint(t)
	case *StartTask:
		err = s.start(t)
	case *DeleteTask:
		err = s.delete(t)
	case *ExitTask:
		err = s.exit(t)
	case *ExecExitTask:
		err = s.execExit(t)
	case *GetContainersTask:
		err = s.getContainers(t)
	case *SignalTask:
		err = s.signal(t)
	case *StatsTask:
		err = s.stats(t)
	case *UpdateTask:
		err = s.updateContainer(t)
	case *UpdateProcessTask:
		err = s.updateProcess(t)
	case *OOMTask:
		err = s.oom(t)
	default:
		err = ErrUnknownTask
	}
	if err != errDeferedResponse {
		i.ErrorCh() <- err
		close(i.ErrorCh())
	}
}


//API start 路口
// api/grpc/server_linux.go 
func (s *apiServer) CreateContainer
	e := &supervisor.StartTask{}
	e.StartResponse = make(chan supervisor.StartResponse, 1)
	createContainerConfigCheckpoint(e, c)  //s是否设置e.checkpoint
	s.sv.SendTask(e) //调用supervisor,放入chanl


supervisor.handleTask:
 从supervisor.tasks拉 chanl:
   case *StartTask:
		err = s.start(t)


//supervisor/create.go:
func (s *Supervisor) start(t *StartTask) error {
	start := time.Now()
	//runtime/container.go:
	container, err := runtime.New(s.stateDir, t.ID, t.BundlePath, s.runtime, s.runtimeArgs, t.Labels)

	s.containers[t.ID] = &containerInfo{
		container: container,
	}

	//supervisor\metrics.go
	ContainersCounter.Inc(1)
	
	task := &startTask{
		Err:           t.ErrorCh(),
		Container:     container,
		StartResponse: t.StartResponse,
		..
	}
	
	task.setTaskCheckpoint(t)

	s.startTasks <- task // 将执行 func (w *worker) Start() 
	ContainerCreateTimer.UpdateSince(start)
	return errDeferedResponse
}

 
func (w *worker) Start() 调用 
//runtime\container_linux.go
 var shimBinary = os.Args[0] + "-shim"  // 即 docker-containerd-shim  ,只有Start 跟  (c *container) Exec 调用 
   func (c *container) Start：
   
	cmd := exec.Command(shimBinary,
		c.id, c.bundle, c.runtime,
	)
	spec, err := c.readSpec()//f, err := os.Open(filepath.Join(c.bundle, "config.json"))
	config := &processConfig{
		checkpoint:  checkpoint,
		root:        processRoot,
		id:          InitProcessID,
		c:           c,
		stdio:       s,
		spec:        spec,
		processSpec: specs.ProcessSpec(spec.Process),
	}
	p, err := newProcess(config)

	 err := c.startCmd(InitProcessID, cmd, p)
	       err := cmd.Start();   
	       err := waitForStart(p, cmd)
	             错误信息通过读文件。。 containerd-shim 跟 containerd 交互 通过文件方式。IPC  file  pipe
	             .... 


//docker-containerd-shim
 main:
  cwd, err := os.Getwd()
   f, err := os.OpenFile(filepath.Join(cwd, "shim-log.json"), os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0666)
   err := start()

  
  func start() error {
  // open the exit pipe
	f, err := os.OpenFile("exit", syscall.O_WRONLY, 0)
	control, err := os.OpenFile("control", syscall.O_RDWR, 0)
	p, err := newProcess(flag.Arg(0), flag.Arg(1), flag.Arg(2))

	err := p.start();
  }

//containerd-shim\process.go
func (p *process) start() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	logPath := filepath.Join(cwd, "log.json")
	args := append([]string{
		"--log", logPath,
		"--log-format", "json",
	}, p.state.RuntimeArgs...)
	... //其他 配置信息
	
	args = append(args,
		"-d",
		"--pid-file", filepath.Join(cwd, "pid"),
		p.id,
	)
	cmd := exec.Command(p.runtime, args...)//调用Runc 
	
	cmd.Dir = p.bundle
	err := cmd.Start()
	 err := cmd.Wait()
	data, err := ioutil.ReadFile("pid")
	pid, err := strconv.Atoi(string(data))

	p.containerPid = pid
}

