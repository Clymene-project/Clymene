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
	"errors"
	"flag"
	"fmt"
	"github.com/Clymene-project/Clymene/pkg/config/tlscfg"
	"github.com/Clymene-project/Clymene/pkg/discovery"
	"github.com/Clymene-project/Clymene/pkg/discovery/grpcresolver"
	"github.com/Clymene-project/Clymene/ports"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	"strconv"
	"strings"
	"time"
)

const (
	gRPCPrefix        = "gateway.grpc-client"
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

// Options Struct to hold configurations
type Options struct {
	// GatewayHostPorts is list of host:port Clymene Gates.
	GatewayHostPorts []string `yaml:"gatewayHostPorts"`

	MaxRetry uint
	TLS      tlscfg.Options

	DiscoveryMinPeers int
	Notifier          discovery.Notifier
	Discoverer        discovery.Discoverer
}

// CreateConnection creates the gRPC connection
func (b *Options) CreateConnection(logger *zap.Logger) (*grpc.ClientConn, error) {
	var dialOptions []grpc.DialOption
	var dialTarget string
	if b.TLS.Enabled { // user requested a secure connection
		logger.Info("Agent requested secure grpc connection to gate(s)")
		tlsConf, err := b.TLS.Config(logger)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %w", err)
		}

		creds := credentials.NewTLS(tlsConf)
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
	} else { // insecure connection
		logger.Info("Agent requested insecure grpc connection to gate(s)")
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	if b.Notifier != nil && b.Discoverer != nil {
		logger.Info("Using external discovery service with roundrobin load balancer")
		grpcResolver := grpcresolver.New(b.Notifier, b.Discoverer, logger, b.DiscoveryMinPeers)
		dialTarget = grpcResolver.Scheme() + ":///round_robin"
	} else {
		if b.GatewayHostPorts == nil {
			return nil, errors.New("at least one GatewayHostPorts hostPort address is required when resolver is not available")
		}
		if len(b.GatewayHostPorts) > 1 {
			scheme := strconv.FormatInt(time.Now().UnixNano(), 36)
			r := manual.NewBuilderWithScheme(scheme)
			var resolvedAddrs []resolver.Address
			for _, addr := range b.GatewayHostPorts {
				resolvedAddrs = append(resolvedAddrs, resolver.Address{Addr: addr})
			}
			r.InitialState(resolver.State{Addresses: resolvedAddrs})
			dialTarget = r.Scheme() + ":///round_robin"
			logger.Info("Agent is connecting to a static list of GatewayHostPorts", zap.String("dialTarget", dialTarget), zap.String("GatewayHostPorts hosts", strings.Join(b.GatewayHostPorts, ",")))
		} else {
			dialTarget = b.GatewayHostPorts[0]
		}
	}
	dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(grpcresolver.GRPCServiceConfig))
	dialOptions = append(dialOptions, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(b.MaxRetry))))
	conn, err := grpc.Dial(dialTarget, dialOptions...)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// AddFlags adds flags for Options.
func AddFlags(flags *flag.FlagSet) {
	flags.Uint(retry, defaultMaxRetry, "Sets the maximum number of retries for a call")
	flags.Int(discoveryMinPeers, 3, "Max number of collectors to which the agent will try to connect at any given time")
	flags.String(gatewayHostPort, "localhost"+ports.PortToHostPort(ports.GatewayGRPC), "Comma-separated string representing host:port of a static list of gateways to connect to directly")
	tlsFlagsConfig.AddFlags(flags)
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (b *Options) InitFromViper(v *viper.Viper) *Options {
	hostPorts := v.GetString(gatewayHostPort)
	if hostPorts != "" {
		b.GatewayHostPorts = strings.Split(hostPorts, ",")
	}
	b.MaxRetry = uint(v.GetInt(retry))
	b.TLS = tlsFlagsConfig.InitFromViper(v)
	b.DiscoveryMinPeers = v.GetInt(discoveryMinPeers)
	return b
}
