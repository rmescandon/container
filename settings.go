// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018 Roberto Mier Escandon <rmescandon@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package container

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/mikefarah/yaml.v2"
)

const (
	settingsFile = "settings.yaml"
)

type networkCfg struct {
	Bridge string `yaml:"bridge"`
	Veth   string `yaml:"veth"`
	CIDR   string `yaml:"cidr"`
}

type cfg struct {
	Rootfs  string     `yaml:"rootfs"`
	Network networkCfg `yaml:"network"`
}

func defaultCfg() *cfg {
	return &cfg{
		Network: networkCfg{
			Bridge: "cbr",
			Veth:   "cveth",
			CIDR:   "192.168.150.1/24",
		},
	}
}

// cfgPath specifies the default location of the config file
func cfgPath() string {
	snapCommon := os.Getenv("SNAP_COMMON")
	if len(snapCommon) > 0 {
		return filepath.Join(snapCommon, settingsFile)
	}
	return filepath.Join("/etc/container", settingsFile)
}

func loadCfg() (*cfg, error) {
	bytes, err := ioutil.ReadFile(cfgPath())
	if err != nil {
		return nil, err
	}

	c := defaultCfg()
	err = yaml.Unmarshal(bytes, c)
	return c, err
}
