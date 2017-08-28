package main

import (
	"os"
	"log"
	"github.com/zpatrick/go-config"
	"github.com/codegangsta/cli"
	"github.com/qnib/doxy/proxy"
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
	debugFlag = cli.BoolFlag{
		Name: "debug",
		Usage: "Print proxy requests",
		EnvVar: "DOXY_DEBUG",
	}
	patternFileFlag = cli.StringFlag{
		Name:  "pattern-file",
		Value: proxy.PATTERN_FILE,
		Usage: "File holding line-separated regex-patterns to be allowed (comments allowed, use #)",
		EnvVar: "DOXY_PATTERN_FILE",
	}
)

func EvalOptions(cfg *config.Config) (po []proxy.ProxyOption) {
	proxySock, _ := cfg.String("proxy-socket")
	po = append(po, proxy.WithProxySocket(proxySock))
	dockerSock, _ := cfg.String("docker-socket")
	po = append(po, proxy.WithDockerSocket(dockerSock))
	debug, _ := cfg.Bool("debug")
	po = append(po, proxy.WithDebugValue(debug))
	return
}

func EvalPatternOpts(cfg *config.Config) (proxy.ProxyOption) {
	patternsFile, _ := cfg.String("pattern-file")
	reader, err := os.Open(patternsFile)
	defer reader.Close()
	patterns := []string{}
	if err != nil {
		log.Printf("Error reading patterns file (%s), using default patterns\n", err.Error())
		return proxy.WithPatterns(proxy.DEFAULT_PATTERNS)
	}
	patterns, err  = proxy.ReadPatterns(reader)
	return proxy.WithPatterns(patterns)
}

func RunApp(ctx *cli.Context) {
	log.Printf("[II] Start Version: %s", ctx.App.Version)
	cfg := config.NewConfig([]config.Provider{config.NewCLI(ctx, true)})
	po := EvalOptions(cfg)
	po = append(po, EvalPatternOpts(cfg))
	p := proxy.NewProxy(po...)
	p.Run()
}

func main() {
	app := cli.NewApp()
	app.Name = "Proxy Docker unix socket to filter out insecure, harmful requests."
	app.Usage = "doxy [options]"
	app.Version = "0.1.1"
	app.Flags = []cli.Flag{
		dockerSocketFlag,
		proxySocketFlag,
		debugFlag,
		patternFileFlag,
	}
	app.Action = RunApp
	app.Run(os.Args)
}