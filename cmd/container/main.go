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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rmescandon/container"
)

const helpUsage = "Show this help message"

func main() {
	flag.Usage = func() {
		helpContent := `Usage: container [COMMAND] <args>

  Available commands:
  	run  Runs into a new container the rest of the parameters in a shell environment

`
		fmt.Fprintf(flag.CommandLine.Output(), helpContent)
		flag.PrintDefaults()
	}

	var help bool
	flag.BoolVar(&help, "-help, -h", false, helpUsage)

	flag.Parse()

	if flag.NArg() == 0 || help {
		flag.Usage()
		os.Exit(1)
	}

	container.Run(flag.Args())
}
