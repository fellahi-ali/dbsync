package main

import (
	"dbsync/transfer"
	"flag"
	"fmt"
	_ "github.com/alexbrainman/odbc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/gops/agent"
	_ "github.com/mattn/go-oci8"
	_ "github.com/viant/bgc"
	"github.com/viant/dsc"
	"log"
	"os"
)

var port = flag.Int("port", 8080, "service port")
var debug = flag.Bool("debug", false, "debug flag")

func main() {
	flag.Parse()
	if *debug {
		dsc.Logf = dsc.StdoutLogger
	}
	fmt.Printf("debug: %v\n", *debug)

	go func() {
		if err := agent.Listen(agent.Options{}); err != nil {
			log.Fatal(err)
		}
	}()
	service := transfer.New(nil)
	server := transfer.NewServer(service, *port)
	go server.StopOnSiginals(os.Interrupt)
	fmt.Printf("dstransfer listening on :%d\n", *port)
	server.ListenAndServe()
}
