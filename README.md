## Container

A simple divertimento for understanding how containers work.

Configure go for giving it a go ;)

You can build it with:

    $ go get -t github.com/rmescandon/container
    $ go install -t github.com/rmescandon/container

And launch it with  

    $ $GOBIN/container run /bin/sh <command>

That command will be executed in a container environment

## Rootfs

The container expects to find a rootfs in a specific path, declared in 
the file */etc/container/settings.yaml*. There you can set:

    rootfs: /<a>/<path>

where */<a>/<path>* points to the rootfs to use into the container.

You can create a ubuntu bionic rootfs easily with:

    $ sudo apt install debootstrap
    $ sudo debootstrap --arch amd64 bionic /<a>/<path>

## TBD

For now only mount, pid, user namespaces work. Network is in progress

## Disclaimer

Tested on Ubuntu Bionic (18.04). This software is delivered as is. No support is provided.