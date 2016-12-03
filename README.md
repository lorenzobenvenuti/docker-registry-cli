A simple command line client for Docker Registry V2.

# Usage

### List repositories

```bash
$ docker-registry-cli -r http://host:port -u user -p password repositories
```

prints

```
lorenzo/ubuntu
lorenzo/busybox
```

### List images

```bash
$ docker-registry-cli -r http://host:port -u user -p password images
```

prints

```
lorenzo/ubuntu:1.0
lorenzo/busybox:1.0
lorenzo/busybox:1.1
```

### Search images

```bash
$ docker-registry-cli -r http://host:port -u user -p password search busy
```

prints

```
lorenzo/busybox:1.0
lorenzo/busybox:1.1
```

### Delete manifest

```bash
$ docker-registry-cli -r http://host:port -u user -p password delete lorenzo/ubuntu:1.0
```

deletes the manifest corresponding to the specified image; this requires that the registry has the `REGISTRY_STORAGE_DELETE_ENABLED` environment variable set to true.
This command deletes the manifest, not the images layers; a robust implementation would require to scan all the images in order to collect the "orphan" layers (a layer could be part of different images). This feature has actually been implemented in version 2.4 of the registry, see [garbage collection](https://docs.docker.com/registry/garbage-collection/). In short, after deleting a manifest, you can free space bu invoking

```bash
$ registry garbage-collect /path/to/config.yml
```

For a container:

```bash
$ docker exec mycontainer /bin/registry garbage-collect /etc/docker/registry/config.yml
```

# TODO

* Global configuration for registry url, username and password
  * Configuration file
  * Environment variables
  * Command (`docker-registry-cli config "registry=http://host:port"`)
* Support regular expressions in search
* Implement garbage-collection (hopefully it will be included in Docker Registry API!?)
* Tree view; something like
```
lorenzo
├─ ubuntu
│  └─ 1.0
└─ busybox
   ├─ 1.0
   └─ 1.1
```
* Use the autocompletion feature of the [kingpin](https://github.com/alecthomas/kingpin) command line parser
