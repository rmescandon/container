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
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

const prompt = "[container] # "

var child func(args []string)

func init() {
	child = doChild
	// Prevent infinite recursive calls to child
	if os.Args[0] == "child" {
		os.Exit(0)
	}
}

/*
Configurable items:
- prompt
- rootfs
- uts mappings
*/

// Execute executes the command
func (r *RunCmd) Execute(args []string) error {
	return run(args)
}

func (r *ChildCmd) Execute(args []string) error {
	child(args)
	return nil
}

func run(args []string) error {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{fmt.Sprintf("PS1=%v", prompt)}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWUTS |
			// syscall.CLONE_NEWIPC |
			// syscall.CLONE_NEWPID |
			// syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getegid(),
				Size:        1,
			},
		},
	}

	return cmd.Run()
}

func doChild(args []string) {
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{fmt.Sprintf("PS1=%v", prompt)}

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error on container: %v", err)
	}
}
