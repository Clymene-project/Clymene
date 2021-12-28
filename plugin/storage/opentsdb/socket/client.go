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

package socket

import (
	"fmt"
	"github.com/Clymene-project/Clymene/plugin/storage/opentsdb/http"
	"github.com/Clymene-project/Clymene/plugin/storage/opentsdb/metricstore/dbmodel"
	"github.com/Clymene-project/Clymene/prompb"
	"go.uber.org/zap"
	"net"
)

type Client struct {
	connections []net.Conn
	converter   *dbmodel.Converter
	hosts       []http.Hosts
	l           *zap.Logger
}
type Options struct {
	Hosts []http.Hosts
}

func (c *Client) SendData(metrics []prompb.TimeSeries) error {
	data, err := c.converter.ConvertTsToOpenTSDBSocket(metrics)
	if err != nil {
		c.l.Error("data convert Error", zap.Error(err))
		return err
	}
	for _, conn := range c.makeConn() {
		_, err = conn.Write(data)
		if err != nil {
			c.l.Error("socket Write Error", zap.Error(err))
		}
	}
	c.closeConn()
	return nil
}

func NewClient(o *Options, converter *dbmodel.Converter, l *zap.Logger) *Client {
	c := &Client{
		converter: converter,
		hosts:     o.Hosts,
		l:         l,
	}
	return c
}

func (c *Client) makeConn() []net.Conn {
	var cons []net.Conn
	for _, h := range c.hosts {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h.Host, h.Port))
		if err != nil {
			c.l.Error("socket connect Error", zap.Error(err))
		}
		cons = append(cons, conn)
	}
	return cons
}

func (c *Client) closeConn() {
	for _, c := range c.connections {
		_ = c.Close()
	}
}
