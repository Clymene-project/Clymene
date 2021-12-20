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

package app

import (
	"flag"
	"github.com/Clymene-project/Clymene/pkg/config/tlscfg"
	"github.com/Clymene-project/Clymene/ports"
	"github.com/spf13/viper"
)

const (
	gatewayGRPCHostPort = "gateway.grpc-server.host-port"
)

var tlsGRPCFlagsConfig = tlscfg.ServerFlagsConfig{
	Prefix:       "gateway.grpc",
	ShowEnabled:  true,
	ShowClientCA: true,
}

// GatewayOptions holds configuration for gateway
type GatewayOptions struct {
	// gatewayGRPCHostPort is the host:port address that the gateway service listens in on for gRPC requests
	gatewayGRPCHostPort string
	// TLSGRPC configures secure transport for gRPC endpoint to collect spans
	TLSGRPC tlscfg.Options
}

// AddFlags adds flags for gatewayOptions
func AddFlags(flags *flag.FlagSet) {
	flags.String(gatewayGRPCHostPort, ports.PortToHostPort(ports.GatewayGRPC), "The host:port (e.g. 127.0.0.1:15610 or :15610) of the gateway's GRPC server")
	tlsGRPCFlagsConfig.AddFlags(flags)
}

// InitFromViper initializes gatewayOptions with properties from viper
func (cOpts *GatewayOptions) InitFromViper(v *viper.Viper) *GatewayOptions {
	cOpts.gatewayGRPCHostPort = ports.FormatHostPort(v.GetString(gatewayGRPCHostPort))
	cOpts.TLSGRPC = tlsGRPCFlagsConfig.InitFromViper(v)

	return cOpts
}
