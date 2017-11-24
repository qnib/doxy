package main

import (
	"os"
	"log"
	"github.com/zpatrick/go-config"
	"github.com/codegangsta/cli"
	"github.com/qnib/doxy/proxy"
	"strings"
)

var (
	dockerSocketFlag = cli.StringFlag{
		Name:  "docker-socket",
		Value: proxy.DOCKER_SOCKET,
		Usage: "Docker host to connect to.",
		EnvVar: "DOXY_DOCKER_SOCKET",
	}
	proxySocketFlag = cli.StringFlag{
		Name:  "proxy-socket",
		Value: proxy.PROXY_SOCKET,
		Usage: "Proxy socket to be created",
		EnvVar: "DOXY_PROXY_SOCKET",
	}
	proxyPatternKey = cli.StringFlag{
		Name:  "pattern-key",
		Value: "default",
		Usage: "pattern key predefined",
		EnvVar: "DOXY_PATTERN_KEY",
	}
	gpuEnabled = cli.BoolFlag{
		Name: "gpu",
		Usage: "Map devices, bind-mounts and environment into each container to allow GPU usage",
		EnvVar: "DOXY_GPU_ENABLED",
	}
	constrainUser = cli.BoolFlag{
		Name: "user-pinning",
		Usage: "Pin user within container to the UID calling the command",
		EnvVar: "DOXY_USER_PINNING_ENABLED",
	}
	debugFlag = cli.BoolFlag{
		Name: "debug",
		Usage: "Print proxy requests",
		EnvVar: "DOXY_DEBUG",
	}
	bindAddFlag = cli.StringFlag{
		Name: "add-binds",
		Usage: "Comma separated list of bind-mounts to add",
		EnvVar: "DOXY_ADDITIONAL_BINDS",
	}
	devMapFlag = cli.StringFlag{
		Name: "device-mappings",
		Usage: "Comma separated list of device mappings",
		EnvVar: "DOXY_DEVICE_MAPPINGS",
	}
	patternFileFlag = cli.StringFlag{
		Name:  "pattern-file",
		Value: proxy.PATTERN_FILE,
		Usage: "File holding line-separated regex-patterns to be allowed (comments allowed, use #)",
		EnvVar: "DOXY_PATTERN_FILE",
	}
	deviceFileFlag = cli.StringFlag{
		Name:  "device-file",
		Value: proxy.DEVICE_FILE,
		Usage: "File holding line-separated devices to be mapped in when in (GPU|HPC) mode (comments allowed, use #)",
		EnvVar: "DOXY_DEVICE_FILE",
	}
)

func EvalOptions(cfg *config.Config) (po []proxy.ProxyOption) {
	proxySock, _ := cfg.String("proxy-socket")
	po = append(po, proxy.WithProxySocket(proxySock))
	dockerSock, _ := cfg.String("docker-socket")
	po = append(po, proxy.WithDockerSocket(dockerSock))
	debug, _ := cfg.Bool("debug")
	po = append(po, proxy.WithDebugValue(debug))
	devMaps, _ := cfg.String("device-mappings")
	gpu, _ := cfg.Bool("gpu")
	po = append(po, proxy.WithGpuValue(gpu))
	po = append(po, proxy.WithDevMappings(strings.Split(devMaps,",")))
	return
}

func EvalPatternOpts(cfg *config.Config) (proxy.ProxyOption) {
	patternsFile, _ := cfg.String("pattern-file")
	reader, err := os.Open(patternsFile)
	defer reader.Close()
	patterns := []string{}
	if err != nil {
		patternsKey, _ := cfg.String("pattern-key")
		if patterns, ok := proxy.PATTERNS[patternsKey]; ok {
			log.Printf("Error reading patterns file '%s', using %s patterns\n", err.Error(), patternsKey)
			return proxy.WithPatterns(patterns)
		}
		log.Printf("Could not find pattern-key '%s'\n", patternsKey)
		os.Exit(1)

	}
	patterns, err  = proxy.ReadLineFile(reader)
	log.Printf("Patterns read from '%s':\n%s\n", patternsFile, strings.Join(patterns,"\n  "))
	return proxy.WithPatterns(patterns)
}

func EvalDevicesOpts(cfg *config.Config) (proxy.ProxyOption) {
	deviceFile, _ := cfg.String("device-file")
	reader, err := os.Open(deviceFile)
	defer reader.Close()
	devices := []string{}
	if err != nil {
		return proxy.WithDevMappings(proxy.DEVICES)
	}
	devices, err  = proxy.ReadLineFile(reader)
	if err != nil {
		log.Printf("Error while reading device file: %s", err.Error())
	}
	return proxy.WithDevMappings(devices)
}

func EvalBindMountOpts(cfg *config.Config) (proxy.ProxyOption) {
	bindStr, _ := cfg.String("add-binds")
	bindMounts := strings.Split(bindStr,",")
	return proxy.WithBindMounts(bindMounts)
}


func RunApp(ctx *cli.Context) {
	log.Printf("[II] Start Version: %s", ctx.App.Version)
	cfg := config.NewConfig([]config.Provider{config.NewCLI(ctx, true)})
	po := EvalOptions(cfg)
	po = append(po, EvalPatternOpts(cfg))
	po = append(po, EvalDevicesOpts(cfg))
	po = append(po, EvalBindMountOpts(cfg))
	p := proxy.NewProxy(po...)
	p.Run()
}

func main() {
	app := cli.NewApp()
	app.Name = "Proxy Docker unix socket to filter out insecure, harmful requests."
	app.Usage = "doxy [options]"
	app.Version = "0.2.0"
	app.Flags = []cli.Flag{
		dockerSocketFlag,
		proxySocketFlag,
		debugFlag,
		gpuEnabled,
		deviceFileFlag,
		patternFileFlag,
		proxyPatternKey,
		bindAddFlag,
	}
	app.Action = RunApp
	app.Run(os.Args)
}