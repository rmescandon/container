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
Source file including all the configuration for the PID namespace of the container
*/

package container

import (
	"os"
	"path/filepath"
	"syscall"
)

func mountProc(newRoot string) error {
	target := filepath.Join(newRoot, "proc")
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	return syscall.Mount(
		"proc",
		target,
		"proc",
		0,
		"",
	)
}
