// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package install has the side-effect of registering all builtin
// service discovery config types.
package install

import (
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/aws"          // register aws
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/azure"        // register azure
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/consul"       // register consul
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/digitalocean" // register digitalocean
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/dns"          // register dns
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/eureka"       // register eureka
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/file"         // register file
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/gce"          // register gce
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/hetzner"      // register hetzner
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/http"         // register http
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/kubernetes"   // register kubernetes
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/linode"       // register linode
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/marathon"     // register marathon
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/moby"         // register moby
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/openstack"    // register openstack
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/puppetdb"     // register puppetdb
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/scaleway"     // register scaleway
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/triton"       // register triton
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/uyuni"        // register uyuni
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/xds"          // register xds
	_ "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/zookeeper"    // register zookeeper
)
