package main

import (
	"os"
	"log"
	"github.com/zpatrick/go-config"
	"github.com/codegangsta/cli"
	"github.com/qnib/doxy/proxy"
)


func RunApp(ctx *cli.Context) {
	log.Printf("[II] Start Version: %s", ctx.App.Version)
	cfg := config.NewConfig([]config.Provider{config.NewCLI(ctx, true)})
	newSock, _ := cfg.String("proxy-socket")
	dockerSock, _ := cfg.String("docker-socket")
	debug, _ := cfg.Bool("debug")
	patternsFile, _ := cfg.String("pattern-file")
	reader, err := os.Open(patternsFile)
	defer reader.Close()
	patterns := []string{}
	if err != nil {
		log.Printf("Error reading patterns file (%s), using default patterns\n", err.Error())
		patterns = []string{
			`^/(v\d\.\d+/)?containers(/\w+)?/(json|stats|top)$`,
			`^/(v\d\.\d+/)?services(/[0-9a-f]+)?$`,
			`^/(v\d\.\d+/)?tasks(/\w+)?$`,
			`^/(v\d\.\d+/)?networks(/\w+)?$`,
			`^/(v\d\.\d+/)?volumes(/\w+)?$`,
			`^/(v\d\.\d+/)?nodes(/\w+)?$`,
			`^/(v\d\.\d+/)?info$`,
			`^/(v\d\.\d+/)?version$`,
			"^/_ping$",
		}
		if debug {
			for i, p := range patterns {
				log.Printf("%-3d: %s\n", i,p)
			}

		}
	} else {
		patterns, err  = proxy.ReadPatterns(reader)
	}
	p := proxy.NewProxy(newSock, dockerSock, debug)
	p.AddPatterns(patterns)
	p.Run()
}

func main() {
	app := cli.NewApp()
	app.Name = "Proxy Docker unix socket to filter out insecure, harmful requests."
	app.Usage = "doxy [options]"
	app.Version = "0.1.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "docker-socket",
			Value: "/var/run/docker.sock",
			Usage: "Docker host to connect to.",
			EnvVar: "DOXY_DOCKER_SOCKET",
		}, cli.StringFlag{
			Name:  "proxy-socket",
			Value: "/tmp/doxy.sock",
			Usage: "Proxy socket to be created",
			EnvVar: "DOXY_PROXY_SOCKET",
		}, cli.BoolFlag{
			Name: "debug",
			Usage: "Print proxy requests",
			EnvVar: "DOXY_DEBUG",
		}, cli.StringFlag{
			Name:  "pattern-file",
			Value: "/etc/doxy.pattern",
			Usage: "File holding line-separated regex-patterns to be allowed (comments allowed, use #)",
			EnvVar: "DOXY_PATTERN_FILE",
	},
	}
	app.Action = RunApp
	app.Run(os.Args)
}