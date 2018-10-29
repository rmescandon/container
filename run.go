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

	// Reexec self process (fork) with cloned namespaces and additional
	// configuration to isolate the command execution
	cmd := exec.Command(selfProc, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Namespaces to clone
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			// syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		// Map container to host users and groups
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

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Could not start command execution - %s", err)
	}

	if err := configureNetworkOnHost(cmd.Process.Pid); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Error waiting for the command execution to finish - %s", err)
	}

	return nil
}

func reexec(args []string) {
	// Setup container hostname
	// (It'll be shown as container prompt)
	hostname := fmt.Sprintf("[container-%v] # ", randStr(6))
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		fmt.Printf("Could not set the hostname - %s\n", err)
		os.Exit(1)
	}

	c, err := loadCfg()
	if err != nil {
		fmt.Printf("Could not read settings - %s\n", err)
		os.Exit(1)
	}

	// Mount /proc.
	// This MUST be done BEFORE PIVOTING. Otherwise it wont be allowed to do it.
	// (now you can check that ps reports only container processes ids
	// and that ls -lah /proc/mounts reports only container mounts but not host's)
	if err := mountProc(c.Rootfs); err != nil {
		fmt.Printf("Could not mount /proc on the new rootfs - %s\n", err)
		os.Exit(1)
	}

	// Pivot root to configured rootfs
	// (now you can check that we have moved to the new path)
	if err := pivotRoot(c.Rootfs); err != nil {
		fmt.Printf("Could not pivot to the new rootfs - %s\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{fmt.Sprintf("PS1=%v", hostname)}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Could not start command execution - %s\n", err)
		os.Exit(1)
	}

	// Configure veth1 interface
	if err := configureNetworkOnContainer(cmd.Process.Pid); err != nil {
		fmt.Printf("Could not configure network on container - %s\n", err)
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for the command execution to finish - %s\n", err)
		os.Exit(1)
	}
}
