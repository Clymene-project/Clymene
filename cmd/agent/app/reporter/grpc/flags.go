/*
 * Copyright (c) 2021 The Clymene Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grpc

import (
	"flag"
	"github.com/Clymene-project/Clymene/pkg/config/tlscfg"
	"github.com/Clymene-project/Clymene/ports"
	"strings"

	"github.com/spf13/viper"
)

const (
	gRPCPrefix        = "reporter.grpc"
	gatewayHostPort   = gRPCPrefix + ".host-port"
	retry             = gRPCPrefix + ".retry.max"
	defaultMaxRetry   = 3
	discoveryMinPeers = gRPCPrefix + ".discovery.min-peers"
)

var tlsFlagsConfig = tlscfg.ClientFlagsConfig{
	Prefix:         gRPCPrefix,
	ShowEnabled:    true,
	ShowServerName: true,
}

// AddFlags adds flags for Options.
func AddFlags(flags *flag.FlagSet) {
	flags.Uint(retry, defaultMaxRetry, "Sets the maximum number of retries for a call")
	flags.Int(discoveryMinPeers, 3, "Max number of collectors to which the agent will try to connect at any given time")
	flags.String(gatewayHostPort, "localhost"+ports.PortToHostPort(ports.GatewayGRPC), "Comma-separated string representing host:port of a static list of gateways to connect to directly")
	tlsFlagsConfig.AddFlags(flags)
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (b *ConnBuilder) InitFromViper(v *viper.Viper) *ConnBuilder {
	hostPorts := v.GetString(gatewayHostPort)
	if hostPorts != "" {
		b.GatewayHostPorts = strings.Split(hostPorts, ",")
	}
	b.MaxRetry = uint(v.GetInt(retry))
	b.TLS = tlsFlagsConfig.InitFromViper(v)
	b.DiscoveryMinPeers = v.GetInt(discoveryMinPeers)
	return b
}
