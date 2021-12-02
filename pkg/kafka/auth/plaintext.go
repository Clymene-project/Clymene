// Copyright (c) 2019 The Jaeger Authors.
//
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

package auth

import (
	"github.com/Shopify/sarama"
)

// PlainTextConfig describes the configuration properties needed for SASL/PLAIN with kafka
type PlainTextConfig struct {
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func setPlainTextConfiguration(config *PlainTextConfig, saramaConfig *sarama.Config) {
	saramaConfig.Net.SASL.Enable = true
	saramaConfig.Net.SASL.User = config.UserName
	saramaConfig.Net.SASL.Password = config.Password
}
