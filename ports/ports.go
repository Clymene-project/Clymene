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

package ports

import (
	"strconv"
	"strings"
)

const (
	GateGRPC = 15610

	AgentAdminHTTP = 15690

	GateAdminHTTP = 15690
)

// PortToHostPort converts the port into a host:port address string
func PortToHostPort(port int) string {
	return ":" + strconv.Itoa(port)
}

// GetAddressFromCLIOptions gets listening address based on port (deprecated flags) or host:port (new flags)
func GetAddressFromCLIOptions(port int, hostPort string) string {
	if port != 0 {
		return PortToHostPort(port)
	}
	return FormatHostPort(hostPort)
}

// FormatHostPort returns hostPort in a usable format (host:port) if it wasn't already
func FormatHostPort(hostPort string) string {
	if hostPort == "" {
		return ""
	}

	if strings.Contains(hostPort, ":") {
		return hostPort
	}

	return ":" + hostPort
}
