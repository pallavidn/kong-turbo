package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/pallavidn/kong-turbo/probe/pkg"
	"github.com/pallavidn/kong-turbo/probe/pkg/conf"
	"fmt"
)

func main() {
	// The default is to log to both of stderr and file
	// These arguments can be overloaded from the command-line args
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "/var/log")
	defer glog.Flush()

	args := conf.NewKongturboArgs(flag.CommandLine)
	flag.Parse()
	fmt.Printf("Kongturbo args %++v\n", args)

	glog.Info("Starting Kongturbo...")
	s, err := pkg.NewKongTAPService(args)

	if err != nil {
		glog.Fatal("Failed creating Kongturbo: %v", err)
	}

	s.Start()
}
