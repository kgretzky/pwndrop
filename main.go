package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kgretzky/pwndrop/api"
	"github.com/kgretzky/pwndrop/config"
	"github.com/kgretzky/pwndrop/core"
	"github.com/kgretzky/pwndrop/log"
	"github.com/kgretzky/pwndrop/storage"
	"github.com/kgretzky/pwndrop/utils"

	"github.com/kgretzky/daemon"
)

const SERVICE_NAME = "pwndrop"
const SERVICE_DESCRIPTION = "pwndrop"

var cfg_path = flag.String("config", "", "config file path")
var debug_log = flag.Bool("debug", false, "log debug output")
var disable_autocert = flag.Bool("no-autocert", false, "disable automatic certificate retrieval")
var disable_dns = flag.Bool("no-dns", false, "disable DNS nameserver")
var show_help = flag.Bool("h", false, "show help")

func usage() {
	fmt.Printf("usage: pwndrop [start|stop|install|remove|status] [-config <config_path>] [-debug] [-no-autocert] [-no-dns] [-h]\n\n")
}

func main() {
	var err error
	ch_exit := make(chan bool, 1)

	dmn, err := daemon.New(SERVICE_NAME, SERVICE_DESCRIPTION, "network.target")
	if err != nil {
		log.Error("daemon: %s", err)
		os.Exit(1)
		return
	}
	svc := &core.Service{dmn}

	if len(os.Args) > 1 {
		var ret bool = false
		var svc_cmd bool = false
		cmd := os.Args[1]
		switch cmd {
		case "install":
			ret = svc.Install()
			svc_cmd = true
		case "remove":
			ret = svc.Remove()
			svc_cmd = true
		case "start":
			ret = svc.Start()
			svc_cmd = true
		case "stop":
			ret = svc.Stop()
			svc_cmd = true
		case "status":
			ret = svc.Status()
			svc_cmd = true
		}
		if svc_cmd {
			if ret {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}
	}

	flag.Parse()

	if *show_help {
		usage()
		return
	}
	if *cfg_path == "" {
		*cfg_path = utils.ExecPath("pwndrop.ini")
	}

	log.Info("pwndrop version: %s", config.Version)

	core.Cfg, err = config.NewConfig(*cfg_path)
	if err != nil {
		log.Fatal("config: %v", err)
		os.Exit(1)
		return
	}
	api.SetConfig(core.Cfg)

	if *debug_log {
		log.SetVerbosityLevel(0)
	}

	db_path := filepath.Join(core.Cfg.GetDataDir(), "pwndrop.db")
	log.Info("opening database at: %s", db_path)

	storage.Open(db_path)
	core.Cfg.HandleSetup()
	if err = core.Cfg.Save(); err != nil {
		log.Fatal("config: %v", err)
		os.Exit(1)
		return
	}

	listen_ip := core.Cfg.GetListenIP()
	log.Debug("listen_ip: %s", listen_ip)
	port_http := core.Cfg.GetHttpPort()
	port_https := core.Cfg.GetHttpsPort()

	_, err = core.NewServer(listen_ip, port_http, port_https, !(*disable_autocert), !(*disable_dns), &ch_exit)
	if err != nil {
		log.Fatal("%v", err)
		os.Exit(1)
		return
	}
	select {
	case _ = <-ch_exit:
		log.Fatal("aborting")
		os.Exit(1)
		return
	}
}
