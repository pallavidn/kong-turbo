package pkg

import (
	"github.com/golang/glog"
	"github.com/pallavidn/kong-turbo/probe/pkg/conf"
	"github.com/pallavidn/kong-turbo/probe/pkg/discovery"
	"github.com/pallavidn/kong-turbo/probe/pkg/registration"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
	"os"
	"os/signal"
	"syscall"
)

type disconnectFromTurboFunc func()

type KongTAPService struct {
	tapService *service.TAPService
}

func NewKongTAPService(args *conf.KongturboArgs) (*KongTAPService, error) {
	tapService, err := createTAPService(args)

	if err != nil {
		glog.Errorf("Error while building turbo TAP service on target %v", err)
		return nil, err
	}

	return &KongTAPService{tapService}, nil
}

func (p *KongTAPService) Start() {
	glog.V(0).Infof("Starting kong TAP service...")

	// Disconnect from Turbo server when kongturbo is shutdown
	handleExit(func() { p.tapService.DisconnectFromTurbo() })

	// Connect to the Turbo server
	p.tapService.ConnectToTurbo()

	select {}
}

func createTAPService(args *conf.KongturboArgs) (*service.TAPService, error) {
	confPath := args.TurboConf

	conf, err := conf.NewKongTurboServiceSpec(confPath)
	if err != nil {
		glog.Errorf("Error while parsing the service config file %s: %v", confPath, err)
		os.Exit(1)
	}

	glog.V(3).Infof("Read service configuration from %s: %++v", confPath, conf)

	communicator := conf.TurboCommunicationConfig
	targetAddr := conf.KongTurboTargetConf.TargetAddress

	registrationClient := &registration.KongTurboRegistrationClient{}
	discoveryClient := discovery.NewDiscoveryClient(targetAddr)

	return service.NewTAPServiceBuilder().
		WithTurboCommunicator(communicator).
		WithTurboProbe(probe.NewProbeBuilder(registration.TargetType, registration.ProbeCategory).
			WithDiscoveryOptions(probe.FullRediscoveryIntervalSecondsOption(int32(*args.DiscoveryIntervalSec))).
			RegisteredBy(registrationClient).
			WithEntityMetadata(registrationClient).
			DiscoversTarget(targetAddr, discoveryClient)).
		Create()
}

// TODO: Move the handle to turbo-sdk-probe as it should be common logic for similar probes
// handleExit disconnects the tap service from Turbo service when kongturbo is terminated
func handleExit(disconnectFunc disconnectFromTurboFunc) {
	glog.V(4).Infof("*** Handling Kongturbo Termination ***")
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP)

	go func() {
		select {
		case sig := <-sigChan:
			// Close the mediation container including the endpoints. It avoids the
			// invalid endpoints remaining in the server side. See OM-28801.
			glog.V(2).Infof("Signal %s received. Disconnecting from Turbo server...\n", sig)
			disconnectFunc()
		}
	}()
}
