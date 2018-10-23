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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const (
	prompt   = "[container] # "
	selfProc = "/proc/self/exe"
)

func init() {
	// Prevent infinite recursive calls to child
	if os.Args[0] == selfProc {
		if len(os.Args) <= 1 {
			fmt.Println("Please provide more parameters")
			os.Exit(1)
		}

		reexec(os.Args[1:])
		os.Exit(0)
	}
}

/*
Configurable items:
- prompt
- rootfs
- uts mappings
*/

// Run runs command in a new container
func Run(args []string) error {
	if len(args) == 0 {
		return errors.New("You must provide a command to run into the container")
	}

	cmd := exec.Command(selfProc, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
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

func reexec(args []string) {
	// Setup container hostname
	hostname := fmt.Sprintf("[container-%v] # ", randStr(6))
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		fmt.Printf("Could not set the hostname - %s\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{fmt.Sprintf("PS1=%v", hostname)}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error on execution - %s\n", err)
		os.Exit(1)
	}
}
