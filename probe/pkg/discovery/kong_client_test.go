package discovery

import "testing"

const (
	KONG_ADMIN = "http://a197f0290a80811e8819f069e7b4672c-1886170686.us-west-2.elb.amazonaws.com:8001"
	KONG_PROXY = "http://a1962c9b2a80811e8819f069e7b4672c-1649750350.us-west-2.elb.amazonaws.com:8000"
)

func TestNewKongClient(t *testing.T) {

	kongClient, _ := NewKongHttpClient(KONG_ADMIN)
	kongClient.GetServices()

	kongClient.GetRoutes()

}
