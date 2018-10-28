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
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func createBridge(name, cidr string) error {
	b := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
	}

	if err := netlink.LinkAdd(b); err != nil {
		return fmt.Errorf("Could not add the bridge - %s", err)
	}

	ip, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("Could not parse input CIDR - %s", err)
	}

	addr := &netlink.Addr{IPNet: &net.IPNet{IP: ip, Mask: subnet.Mask}}
	if err := netlink.AddrAdd(b, addr); err != nil {
		return fmt.Errorf("Could not add address to the bridge - %s", err)
	}

	if err := netlink.LinkSetUp(b); err != nil {
		return fmt.Errorf("Could not setup bridge - %s", err)
	}

	return nil
}

func createVeths(name string) error {
	veth0 := fmt.Sprintf("%s0", name)
	veth1 := fmt.Sprintf("%s1", name)

	// Return them if already exists
	_, err := net.InterfaceByName(veth0)
	if err == nil {
		_, err = net.InterfaceByName(veth1)
		return err
	}

	v := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name: veth0,
		},
		PeerName: veth1,
	}

	if err := netlink.LinkAdd(v); err != nil {
		return fmt.Errorf("Could not add the veth devs - %s", err)
	}

	if err := netlink.LinkSetUp(v); err != nil {
		return fmt.Errorf("Could not setup veth devs - %s", err)
	}

	return nil
}

func attach(veth0, bridge string) error {
	v, err := netlink.LinkByName(veth0)
	if err != nil {
		return fmt.Errorf("Could not get veth by name - %s", err)
	}

	b, err := netlink.LinkByName(bridge)
	if err != nil {
		return fmt.Errorf("Could not get bridge by name - %s", err)
	}

	return netlink.LinkSetMaster(v, b.(*netlink.Bridge))
}

func moveToNetworkNamespace(veth1 string, pid int) error {
	v, err := netlink.LinkByName(veth1)
	if err != nil {
		return err
	}

	return netlink.LinkSetNsPid(v, pid)
}

func addDefaultRoute(veth1, bridgeIP string) error {
	v, err := netlink.LinkByName(veth1)
	if err != nil {
		return err
	}

	ip, _, err := net.ParseCIDR(bridgeIP)
	if err != nil {
		return fmt.Errorf("Could not parse bridge IP - %s", err)
	}

	route := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: v.Attrs().Index,
		Gw:        ip,
	}

	return netlink.RouteAdd(route)
}
