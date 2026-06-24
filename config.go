// Copyright 2026 Byterio
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memkey

import "time"

// Config defines the config for memkey.
type Config struct {
	// Time before deleting expired keys.
	//
	// Default is 10 * time.Second
	CleanupInterval time.Duration
}

// ConfigDefault is the default config.
var ConfigDefault = Config{
	CleanupInterval: 10 * time.Second,
}

// configDefault is a helper function to set default values.
func configDefault(config ...Config) Config {
	// Return default config if nothing provided.
	if len(config) < 1 {
		return ConfigDefault
	}
	// Override default config.
	cfg := config[0]
	// Set default values.
	if int(cfg.CleanupInterval.Seconds()) <= 0 {
		cfg.CleanupInterval = ConfigDefault.CleanupInterval
	}
	return cfg
}
