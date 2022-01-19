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

package tdengine

import (
	"flag"
	"github.com/spf13/viper"
	"time"
)

const (
	configPrefix = "tdengine"

	suffixHostName             = ".hostname"
	suffixServerPort           = ".server-port"
	suffixUser                 = ".user"
	suffixPassword             = ".password"
	suffixDBName               = ".dbname"
	suffixTablePrefix          = ".table-prefix"
	suffixMode                 = ".mode"
	suffixNumOftables          = ".num-of-tables"
	suffixNumOfRecordsPerTable = ".num-of-record-per-table"
	suffixNumOfRecordsPerReq   = ".num-of-record-per-req"
	suffixNumOfThreads         = ".num-of-threads"
	suffixStartTimestamp       = ".start-timestamp"

	suffixMaxConnect  = ".max-connect"
	suffixMaxIdle     = ".max-idle"
	suffixIdleTimeout = ".idle-timeout"

	defaultHostName             = "127.0.0.1"
	defaultServerPort           = 6030
	defaultUser                 = "root"
	defaultPassword             = "taosdata"
	defaultDBName               = "test"
	defaultTablePrefix          = "d"
	defaultMode                 = "r"
	defaultNumOftables          = 2
	defaultNumOfRecordsPerTable = 10
	defaultNumOfRecordsPerReq   = 3
	defaultNumOfThreads         = 1
	defaultStartTimestamp       = "2020-10-01 08:00:00"

	defaultMaxConnect  = 4000
	defaultMaxIdle     = 4000
	defaultIdleTimeout = time.Hour

	defaultSupTblName = "meters"
	defaultKeep       = 365 * 20
	defaultDays       = 30
)

type Options struct {
	hostName             string
	serverPort           int
	user                 string
	password             string
	dbName               string
	supTblName           string
	tablePrefix          string
	mode                 string
	numOftables          int
	numOfRecordsPerTable int
	numOfRecordsPerReq   int
	numOfThreads         int
	startTimestamp       string
	startTs              int64

	maxConnect  int
	maxIdle     int
	idleTimeout time.Duration

	keep int
	days int
}

func (o *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		configPrefix+suffixHostName,
		defaultHostName,
		"The host to connect to TDengine server.",
	)
	flagSet.Int(
		configPrefix+suffixServerPort,
		defaultServerPort,
		"he TCP/IP port number to use for the connection to TDengine server",
	)
	flagSet.String(
		configPrefix+suffixUser,
		defaultUser,
		"The TDengine user name to use when connecting to the server",
	)
	flagSet.String(
		configPrefix+suffixPassword,
		defaultPassword,
		"The password to use when connecting to the server",
	)
	flagSet.String(
		configPrefix+suffixDBName,
		defaultDBName,
		"Destination database",
	)
	flagSet.String(
		configPrefix+suffixTablePrefix,
		defaultTablePrefix,
		"Table prefix name",
	)
	flagSet.String(
		configPrefix+suffixMode,
		defaultMode,
		"mode,r:raw,s:stmt",
	)
	flagSet.Int(
		configPrefix+suffixNumOftables,
		defaultNumOftables,
		"The number of tables.",
	)
	flagSet.Int(
		configPrefix+suffixNumOfRecordsPerTable,
		defaultNumOfRecordsPerTable,
		"The number of records per table",
	)
	flagSet.Int(
		configPrefix+suffixNumOfRecordsPerReq,
		defaultNumOfRecordsPerReq,
		"The number of records per request",
	)
	flagSet.Int(
		configPrefix+suffixNumOfThreads,
		defaultNumOfThreads,
		"The number of threads",
	)
	flagSet.String(
		configPrefix+suffixStartTimestamp,
		defaultStartTimestamp,
		"The start timestamp for one table",
	)
	flagSet.Int(
		configPrefix+suffixMaxConnect,
		defaultMaxConnect,
		"max connections to taosd",
	)
	flagSet.Int(
		configPrefix+suffixMaxIdle,
		defaultMaxIdle,
		"max idle connections to taosd",
	)
	flagSet.Duration(
		configPrefix+suffixIdleTimeout,
		defaultIdleTimeout,
		"Set idle connection timeout",
	)
}

func (o *Options) InitFromViper(v *viper.Viper) {
	o.hostName = v.GetString(configPrefix + suffixHostName)
	o.serverPort = v.GetInt(configPrefix + suffixServerPort)
	o.user = v.GetString(configPrefix + suffixUser)
	o.password = v.GetString(configPrefix + suffixPassword)
	o.dbName = v.GetString(configPrefix + suffixDBName)
	o.tablePrefix = v.GetString(configPrefix + suffixTablePrefix)
	o.mode = v.GetString(configPrefix + suffixMode)
	o.numOftables = v.GetInt(configPrefix + suffixNumOftables)
	o.numOfRecordsPerTable = v.GetInt(configPrefix + suffixNumOfRecordsPerTable)
	o.numOfRecordsPerReq = v.GetInt(configPrefix + suffixNumOfRecordsPerReq)
	o.numOfThreads = v.GetInt(configPrefix + suffixNumOfThreads)
	o.startTimestamp = v.GetString(configPrefix + suffixStartTimestamp)

	startTs, err := time.ParseInLocation("2006-01-02 15:04:05", o.startTimestamp, time.Local)
	if err == nil {
		o.startTs = startTs.UnixNano() / 1e6
	}

	o.maxConnect = v.GetInt(configPrefix + suffixMaxConnect)
	o.maxIdle = v.GetInt(configPrefix + suffixMaxIdle)
	o.idleTimeout = v.GetDuration(configPrefix + suffixIdleTimeout)

	o.supTblName = defaultSupTblName
	o.keep = defaultKeep
	o.days = defaultDays

}
