package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/swanwish/go-common/utils"
	"github.com/swanwish/go-common/web"
	"github.com/swanwish/logproxyservice/handlers/root"
)

const (
	DEFAULT_PORT int64 = 8080
)

var (
	port int64
)

func parseCmdLineArgs() {
	flag.Int64Var(&port, "port", DEFAULT_PORT, "The port to listen")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	parseCmdLineArgs()

	web.InitHandlers([]web.RouterHandlers{
		&root.Handlers{},
	})

	if port == 0 {
		port = DEFAULT_PORT
	}

	localIps, err := utils.GetLocalIPAddrs()
	if err != nil {
		fmt.Println("Failed to get local ip addresses.")
		return
	}

	fmt.Printf("Service listen on port \x1b[31;1m%d\x1b[0m and server ip addresses are \x1b[31;1m%s\x1b[0m\n", port, strings.Join(localIps, ", "))

	httpAddr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		fmt.Printf("http.ListendAndServer() failed with %#v\n", err)
	}
	fmt.Println("Exited")
}
