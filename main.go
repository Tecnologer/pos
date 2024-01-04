package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/tecnologer/pos/server"
	"os"
)

var (
	port    = flag.Int("port", 8080, "Port to listen on")
	verbose = flag.Bool("verbose", false, "Print verbose logs")
)

func main() {
	flag.Parse()

	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logfile, err := os.OpenFile("sale.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.SetOutput(logfile)

	server.Run(*port)
}
