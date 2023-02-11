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

package gateway

import (
	"flag"
	"github.com/Clymene-project/Clymene/plugin/storage/gateway/grpc"
	"github.com/Clymene-project/Clymene/plugin/storage/gateway/http"
	"github.com/spf13/viper"
)

type Options struct {
	ServiceType string
	grpcOptions grpc.Options
	httpOptions http.Options
}

const (
	configPrefix      = "gateway"
	suffixServiceType = ".service-type"

	defaultServiceType = "grpc"
)

func (o *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		configPrefix+suffixServiceType,
		defaultServiceType,
		"Setting the type of gateway server (grpc or http)",
	)
	grpc.AddFlags(flagSet)
	http.AddFlags(flagSet)
}
func (o *Options) InitFromViper(v *viper.Viper) {
	o.ServiceType = v.GetString(configPrefix + suffixServiceType)
	o.grpcOptions.InitFromViper(v)
	o.httpOptions.InitFromViper(v)
}
