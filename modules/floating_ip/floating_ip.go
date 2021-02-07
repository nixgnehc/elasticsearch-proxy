/*
Copyright Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package floating_ip

import (
	log "github.com/cihub/seelog"
	"github.com/nixgnehc/infini-framework/core/config"
	"github.com/nixgnehc/infini-framework/core/env"
	"github.com/nixgnehc/infini-framework/core/net"
)

type FloatingIPConfig struct {
	Enabled   bool   `config:"enabled"`
	IP        string `config:"ip"`
	Netmask   string `config:"netmask"`
	Interface string `config:"interface"`
	Priority  int    `config:"priority"`
}

type FloatingIPPlugin struct {
}

func (this FloatingIPPlugin) Name() string {
	return "Floating_IP"
}

var (
	floatingIPConfig = FloatingIPConfig{
		Enabled:  false,
		Netmask:  "255.255.255.0",
		Priority: 1,
	}
)

func (module FloatingIPPlugin) Setup(cfg *config.Config) {
	cfg.Unpack(&floatingIPConfig)
}

func (module FloatingIPPlugin) Start() error {
	log.Info("setup floating IP, root privilege are required")
	err := net.SetupAlias(floatingIPConfig.Interface, floatingIPConfig.IP, floatingIPConfig.Netmask)
	if err != nil {
		panic(err)
	}

	apiConfig := &config.APIConfig{}
	env.ParseConfig("api", apiConfig)
	log.Infof("high availability address: %s://%s:%s", apiConfig.GetSchema(), floatingIPConfig.IP, apiConfig.NetworkConfig.GetBindingPort())
	return nil
}

func (module FloatingIPPlugin) Stop() error {
	net.DisableAlias(floatingIPConfig.Interface, floatingIPConfig.IP, floatingIPConfig.Netmask)
	return nil
}
