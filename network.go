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

/*
Source file including all the configuration for the NETWORK namespace of the container
*/

package container

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func configureNetworkOnHost(pid int) error {
	if err := createBridge("thebridge", "10.20.30.1/24"); err != nil {
		return err
	}

	if err := createVeths("veth"); err != nil {
		return err
	}

	if err := attach("veth0", "thebridge"); err != nil {
		return err
	}

	if err := moveToNetworkNamespace("veth1", pid); err != nil {
		return err
	}

	return nil
}

func configureNetworkOnContainer(pid int) error {
	l, err := netlink.LinkByName("veth1")
	if err != nil {
		return err
	}

	ip, subnet, err := net.ParseCIDR("10.20.30.2/24")
	if err != nil {
		return fmt.Errorf("Could not parse veth2 IP - %s", err)
	}

	addr := &netlink.Addr{IPNet: &net.IPNet{IP: ip, Mask: subnet.Mask}}
	err = netlink.AddrAdd(l, addr)
	if err != nil {
		return errors.New("Unable to assign IP address '10.20.30.2' to veth1")
	}

	if err := netlink.LinkSetUp(l); err != nil {
		return err
	}

	return addDefaultRoute("veth1", "10.20.30.1")
}
