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
Source file including all the configuration for the MOUNT namespace of the container
*/

package container

import (
	"os"
	"path/filepath"
	"syscall"
)

func pivotRoot(newRoot string) error {
	oldRoot := filepath.Join(newRoot, "/.pivot_root")

	// bind mount new root to itself - this is a slight hack
	// needed to work around a pivot_root requirement
	if err := syscall.Mount(
		newRoot,
		newRoot,
		"",
		syscall.MS_BIND|syscall.MS_REC,
		"",
	); err != nil {
		return err
	}

	if err := os.MkdirAll(oldRoot, 0700); err != nil {
		return err
	}

	if err := syscall.PivotRoot(newRoot, oldRoot); err != nil {
		return err
	}

	if err := os.Chdir("/"); err != nil {
		return err
	}

	// unmount oldRoot (notice that now it is on /.pivot_root)
	oldRoot = "/.pivot_root"
	return syscall.Unmount(oldRoot, syscall.MNT_DETACH)
}
