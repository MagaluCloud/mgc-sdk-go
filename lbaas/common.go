package lbaas

import (
	"strings"
)

// urlNetworkLoadBalancer constructs the URL path for network load balancer operations
// If lbID is provided, it appends the ID to the base path
// Additional path segments can be provided via extraPath parameter
func urlNetworkLoadBalancer(lbID *string, extraPath ...string) string {
	result := "/v0beta1/network-load-balancers"
	if lbID == nil || *lbID == "" {
		return result
	}
	result += "/" + *lbID
	if len(extraPath) > 0 {
		result += "/" + strings.Join(extraPath, "/")
	}
	return result
}
