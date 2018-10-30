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
	if exists(name) {
		return nil
	}

	b := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
	}

	if err := netlink.LinkAdd(b); err != nil {
		return fmt.Errorf("Could not add the bridge - %s", err)
	}

	return configureLinkFromCIDR(name, cidr)
}

func createVethPair(name string) error {
	veth0 := fmt.Sprintf("%s0", name)
	veth1 := fmt.Sprintf("%s1", name)

	// assume that if veth0 exists, also does veth1
	if exists(veth0) {
		return nil
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
		return fmt.Errorf("Could not get veth %s by name - %s", veth0, err)
	}

	b, err := netlink.LinkByName(bridge)
	if err != nil {
		return fmt.Errorf("Could not get bridge %s by name - %s", bridge, err)
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

func configureLinkFromCIDR(name, cidr string) error {
	ip, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("Could not parse %s CIDR - %s", name, err)
	}

	netip := &net.IPNet{IP: ip, Mask: subnet.Mask}
	return configureLinkFromIP(name, netip)
}

func configureLinkFromIP(name string, ipnet *net.IPNet) error {
	l, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	addr := &netlink.Addr{IPNet: ipnet}
	err = netlink.AddrAdd(l, addr)
	if err != nil {
		return fmt.Errorf("Unable to assign address %s to %s", ipnet.String(), name)
	}

	return netlink.LinkSetUp(l)
}

func addDefaultRoute(veth string, bridgeIP net.IP) error {
	v, err := netlink.LinkByName(veth)
	if err != nil {
		return err
	}

	route := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: v.Attrs().Index,
		Gw:        bridgeIP,
	}

	return netlink.RouteAdd(route)
}

func exists(name string) bool {
	_, err := net.InterfaceByName(name)
	return err == nil
}
