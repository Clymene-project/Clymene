// Copyright (c) 2017 The Jaeger Authors.
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

package version

var (
	// Version is set during binary building (git revision number)
	Version string
	// BuildTime is set during binary building
	BuildTime string
)

// Info holds build information
type Info struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
}

func Set(version, buildTime string) {
	Version = version
	BuildTime = buildTime
}

// Get creates and initialized Info object
func Get() Info {
	return Info{
		Version:   Version,
		BuildTime: BuildTime,
	}
}
