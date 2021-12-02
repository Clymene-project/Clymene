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

package storage

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// TypeEnvVar is the name of the env var that defines the type of backend used for time series storage.
	TypeEnvVar = "TS_STORAGE_TYPE"

	// DependencyStorageTypeEnvVar is the name of the env var that defines the type of backend used for dependencies storage.
	DependencyStorageTypeEnvVar = "DEPENDENCY_STORAGE_TYPE"

	tsStorageFlag = "--ts-storage.type"
)

// FactoryConfig tells the Factory which types of backends it needs to create for different storage types.
type FactoryConfig struct {
	WriterTypes             []string
	ReaderType              string
	DependenciesStorageType string
}

// FactoryConfigFromEnvAndCLI reads the desired types of storage backends from TS_STORAGE_TYPE and
// DEPENDENCY_STORAGE_TYPE environment variables. Allowed values:
//   * `influxdb` - built-in
//   * `elasticsearch` - built-in
//   * `prometheus` - built-in
//   * `kafka` - built-in
//   * `gateway` - built-in
//   * `cortex` - built-in
//   * `kdb` - built-in
//   * `opentsdb` - built-in
//   * `plugin` - loads a dynamic plugin that implements storage.Factory interface (not supported at the moment)
//
// For backwards compatibility it also parses the args looking for deprecated --ts-storage.type flag.
// If found, it writes a deprecation warning to the log.
func FactoryConfigFromEnvAndCLI(args []string, log io.Writer) FactoryConfig {
	spanStorageType := os.Getenv(TypeEnvVar)
	if spanStorageType == "" {
		// for backwards compatibility check command line for --span-storage.type flag
		spanStorageType = spanStorageTypeFromArgs(args, log)
	}
	if spanStorageType == "" {
		spanStorageType = elasticsearchStorageType
	}
	spanWriterTypes := strings.Split(spanStorageType, ",")
	if len(spanWriterTypes) > 1 {
		fmt.Fprintf(log,
			"WARNING: multiple span storage types have been specified. "+
				"Only the first type (%s) will be used for reading and archiving.\n\n",
			spanWriterTypes[0],
		)
	}
	depStorageType := os.Getenv(DependencyStorageTypeEnvVar)
	if depStorageType == "" {
		depStorageType = spanWriterTypes[0]
	}
	// TODO support explicit configuration for readers
	return FactoryConfig{
		WriterTypes:             spanWriterTypes,
		ReaderType:              spanWriterTypes[0],
		DependenciesStorageType: depStorageType,
	}
}

func spanStorageTypeFromArgs(args []string, log io.Writer) string {
	for i, token := range args {
		if i == 0 {
			continue // skip app name; easier than dealing with +-1 offset
		}
		if !strings.HasPrefix(token, tsStorageFlag) {
			continue
		}
		fmt.Fprintf(
			log,
			"WARNING: found deprecated command line option %s, please use environment variable %s instead\n",
			token,
			TypeEnvVar,
		)
		if token == tsStorageFlag && i < len(args)-1 {
			return args[i+1]
		}
		if strings.HasPrefix(token, tsStorageFlag+"=") {
			return token[(len(tsStorageFlag) + 1):]
		}
		break
	}
	return ""
}
