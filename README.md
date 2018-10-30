# Container

A simple divertimento for understanding how containers work.
This tool allows to execute a command in a very basic container environment.


## Rootfs

Container expects to find a rootfs in some local path. One easy
way to create a rootfs is by using **debootstrap**:

    $ sudo apt update
    $ sudo apt install debootstrap


For example, You can create a ubuntu bionic rootfs easily with:

    $ sudo debootstrap --arch amd64 bionic /<a>/<path>

where */a/path* points to a path in your local disk where
placing the rootfs to be used into the container.


The container searches for the rootfs in the path declared in
the file **/etc/container/settings.yaml**. You must create such
file and provide the setttings:

    rootfs: /<a>/<path>

## Build

Configure go for giving it a go ;)

    $ sudo apt install golang
    $ echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    $ echo 'export PATH=${PATH}:${GOPATH}/bin' >> ~/.bashrc
    $ source ~/.bashrc

Checkout and build container tool with:

    $ go get -t github.com/rmescandon/container
    $ go install -t github.com/rmescandon/container

## Run

Launch container tool with:

    $ sudo $GOBIN/container run <list_of_commands_to_execute>

for example, listing the files in the rootfs:

    $ sudo $GOBIN/container run /bin/sh ls -la

or simply entering into the container for a further cli executions:

    $ sudo $GOBIN/container run /bin/sh

All those commands will be executed in a container environment, isolated (more or less)
from host.

You can check that

* container rootfs is at the settings configured path and host rootfs is not visible from within the container
* `ls -lah /proc/mounts` reports only container mounts but not host's
* `ip link` shows container interfaces but not the host's
* `id` into the container is the root one, but does not have root permissions over the host
* container hostname is different from the host one

!!!Note:
    It is needed _sudo_ ing the command execution because setuid trick is still not
    implemented and all the networking configuration requires host root privileges.
    but in future this won't be required. Sorry for the temporary inconvenience

## Settings

Not only rootfs can be configured in the **/etc/container/settings.yaml** file but
the bridge name, virtual devices or the bridge CIDR, like:

    rootfs: /data/rootfs/bionic
    network:
        bridge: thebridge
        veth : theveth
        cidr: 10.20.30.40/24

## TBD

For now only mount, pid, user namespaces work.
Network is in progress. It is possible to ping host ips, but not
reaching internet from the container

## Disclaimer

Tested on Ubuntu Bionic (18.04). This software is delivered as is. No support is provided.