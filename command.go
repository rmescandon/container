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

// Cmd holds the cli arguments
type Cmd struct {
	Run   RunCmd    `command:"run" description:"Runs into a new container the rest of the parameters in a shell environment"`
	Child ReexecCmd `command:"reexec" description:"Runs child into a new container" hidden:"true"`
}

// RunCmd command to prepare a shell session into a container
type RunCmd struct{}

// ReexecCmd command re-executing a Run with container namespaces already set up
type ReexecCmd struct{}
