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
)

func configureNetworkOnHost(bridge, veth, cidr string, pid int) error {
	if err := createBridge(bridge, cidr); err != nil {
		return err
	}

	if err := createVethPair(veth); err != nil {
		return err
	}

	veth0 := fmt.Sprintf("%s0", veth)
	if err := attach(veth0, bridge); err != nil {
		return err
	}

	veth1 := fmt.Sprintf("%s1", veth)
	return moveToNetworkNamespace(veth1, pid)
}

func configureNetworkOnContainer(veth, cidr string, pid int) error {
	ip, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("Could not parse CIDR - %s", err)
	}

	// move to the next IP, that will be the one for veth1
	next := nextIP(ip)
	ipnet := &net.IPNet{IP: next, Mask: subnet.Mask}

	veth1 := fmt.Sprintf("%s1", veth)
	if err := configureLinkFromIP(veth1, ipnet); err != nil {
		return err
	}

	return addDefaultRoute(veth1, ip)
}

func nextIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	dup = dup.To4()
	dup[3]++
	return dup
}
