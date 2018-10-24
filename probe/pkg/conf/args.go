package conf

import (
	"flag"
)

const (
	defaultDiscoveryIntervalSec = 600
	DefaultConfPath    = "/etc/kongturbo/turbo.config"
)

type KongturboArgs struct {
	DiscoveryIntervalSec *int
	TurboConf string
}

func NewKongturboArgs(fs *flag.FlagSet) *KongturboArgs {
	p := &KongturboArgs{}

	p.DiscoveryIntervalSec = fs.Int("discovery-interval-sec", defaultDiscoveryIntervalSec, "The discovery interval in seconds")
	fs.StringVar(&p.TurboConf,  "turboconfig", p.TurboConf,  "Path to the config file.")

	return p
}
