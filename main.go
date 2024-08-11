package main

import (
	"fileagent/config"
	"fileagent/server"
	"flag"
	"fmt"
	"github.com/kardianos/service"
	"log"
	"runtime"
	"strconv"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func getArgs() int64 {

	portstr := config.Conf.Port
	if portstr == "" {
		portstr = "8888"
	}

	port, err := strconv.ParseInt(portstr, 10, 64)
	if err != nil {
		log.Fatalln("cannot parse port:", portstr)
	}

	return port
}

var (
	install     = flag.Bool("install", false, "Install the service")
	uninstall   = flag.Bool("uninstall", false, "Uninstall the service")
	start       = flag.Bool("start", false, "Start the service")
	stop        = flag.Bool("stop", false, "Stop the service")
	versionFlag = flag.Bool("version", false, "version")
)
var version = "1.0.0"

func main() {

	flag.Parse()

	svcConfig := &service.Config{
		Name:        "File-Agent",
		DisplayName: "File-Agent Service",
		Description: "File-Agent Service.",
	}

	program := &ServiceProgram{}
	s, err := service.New(program, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 根据命令行参数执行对应操作
	if *install {
		err := s.Install()
		if err != nil {
			log.Fatalf("Failed to install service: %v", err)
		}
		fmt.Println("Service installed successfully.")
	} else if *uninstall {
		err := s.Uninstall()
		if err != nil {
			log.Fatalf("Failed to uninstall service: %v", err)
		}
		fmt.Println("Service uninstalled successfully.")
	} else if *start {
		err := s.Start()
		if err != nil {
			log.Fatalf("Failed to start service: %v", err)
		}
		fmt.Println("Service started successfully.")
	} else if *stop {
		err := s.Stop()
		if err != nil {
			log.Fatalf("Failed to stop service: %v", err)
		}
		fmt.Println("Service stopped successfully.")
	} else if *versionFlag {
		fmt.Printf("File-Agent version: %s\n", version)
	} else {
		// 如果没有传递任何命令行参数，则作为服务运行
		err = s.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// ServiceProgram 实现了服务的逻辑
type ServiceProgram struct{}

func (p *ServiceProgram) Start(s service.Service) error {
	// 你可以在这里启动你服务的主要逻辑
	port := getArgs()
	server.StartHttp(port)
	return nil
}

func (p *ServiceProgram) Stop(s service.Service) error {
	// 你可以在这里实现服务的清理逻辑
	log.Println("Service is stopping...")
	return nil
}
