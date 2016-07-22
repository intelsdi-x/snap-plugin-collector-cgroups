# snap-plugin-collector-cgroups

Snap plugin for collecting cgroups metrics using [libcontainer](https://github.com/opencontainers/runc/tree/master/libcontainer) library for cgroup statistics.

*WARNING*: This plugin requires root privileges to run. Use with caution.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements
* Snap daemon (`snapd`) started with root permissions
* [golang 1.5+](https://golang.org/dl/) - needed only for building

### Operating systems
All OSs currently supported by Snap:
* Linux/amd64

### Installation
#### Download plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-cgroups
Clone repo into `$GOPATH/src/github/intelsdi-x/`:

    $ git clone https://github.com/<yourGithubID>/snap-plugin-collector-cgroups

Build the plugin by running make in repo:

    $ make

This builds the plugin in `/build/sysdevfs`

## Documentation

### Collected Metrics
This plugin has the ability to gather all rootfs cgroups metrics from /sys/fs/cgroup. See [METRICS.md](METRICS.md) for full metrics list.

### Examples
Here is an example of using Snap cgroups collector and getting cgroups data remotely via REST API in real time.

A Task Manifest is required to select what data you want to collect. You can find an example Task Manifest file in [examples/tasks/](examples/tasks/) directory:

    {
        "version": 1,
        "schedule": {
            "type": "simple",
            "interval": "5s"
        },
        "workflow": {
            "collect": {
                "metrics": {
                    "/intel/linux/cgroups/hugetlb_stats/*/*/failcnt": {},
                    "/intel/linux/cgroups/hugetlb_stats/*/*/max_usage": {},
                    "/intel/linux/cgroups/hugetlb_stats/*/*/usage": {}
                }
            }
        }
    }


This Task Manifest collects a portion of statistics from all cgroup subsystems (blkio, cpu, hugetlb, memory, pids) continously every 5 seconds.

### Examples
If you've got Snap installed, you need to start snap daemon with root permissions and plugin trust mode:

    $ sudo snapd -t 0

Then, in another terminal window, load the plugin:

    $ snapctl plugin load snap-plugin-collector-cgroups
    Plugin loaded
    Name: cgroups
    Version: 2
    Type: collector
    Signed: false
    Loaded Time: Thu, 12 May 2016 11:36:40 CEST

To see all available metrics:

    $ snapctl metric list

To load example task (examples/tasks/cgroups.json):

    $ snapctl task create -t cgroups.json
    Using task manifest to create task
    Task created
    ID: 99541dbf-0395-4858-a05d-f62937fc6b82
    Name: Task-99541dbf-0395-4858-a05d-f62937fc6b82
    State: Running

You can preview task responses using:

    $ snapctl task watch 99541dbf-0395-4858-a05d-f62937fc6b82

Example response:

    Watching Task (012bd95e-bfc5-493f-875f-4c85dd7a2f61):
    NAMESPACE                                                DATA    TIMESTAMP
    /intel/linux/cgroups/hugetlb_stats/root/2MB/failcnt      0       2016-05-24 15:34:45.455812812 +0200 CEST
    /intel/linux/cgroups/hugetlb_stats/root/2MB/max_usage    0       2016-05-24 15:34:45.455933449 +0200 CEST
    /intel/linux/cgroups/hugetlb_stats/root/2MB/usage        0       2016-05-24 15:34:45.456032188 +0200 CEST

Now if your task works and returns data, you can get the data via REST API:

    curl -L http://localhost:8181/v1/tasks/\<task id\>

That's it!

### Roadmap
There isn't a current roadmap for this plugin. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-cgroups/issues) and feel free to then submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-cgroups/pulls).


## Community Support
This repository is one of **many** plugins in **Snap**, the open telemetry framework. See the full project at http://github.com/intelsdi-x/snap. To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Mateusz Kleina](https://github.com/mkleina)
* Author: [Marcin Olszewski](https://github.com/marcintao)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
