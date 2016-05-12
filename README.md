# snap-plugin-collectoir-cgroups
Snap plugin for collecting cgroups metrics from /sys/fs/cgroup filesystem using [libcontainer](https://github.com/opencontainers/runc/tree/master/libcontainer) library for gatchering cgroup stats.

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

 Plugin collects specified metrics in-band on OS level.

### System Requirements

 - Linux system
 - Snap daemon started with root permissions

### Installation
#### Download cgroups plugin binary:
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
This plugin has the ability to gather all rootfs cgroups metrics:

| Namespace | Description |
|-----------|-------------|
| /intel/linux/cgroups/blkio_stats/*| Block I/O subsystem metrics |
| /intel/linux/cgroups/cpu_stats/*| CPU subsystem metrics |
| /intel/linux/cgroups/hugetlb_stats/*| HugeTLB subsystem metrics |
| /intel/linux/cgroups/memory_stats/*| Memory subsystem metrics |
| /intel/linux/cgroups/pids_stats/*| Process IDs subsystem metrics |

### Examples
Example of using snap cgroups collector and getting cgroups data remotely via REST API in real time.

First of all, it is needed to create task manifest, to tell cgroups collector plugin what data you want to collect. You can find example task manifest file in examples/tasks/ dir:

    {
        "version": 1,
        "schedule": {
            "type": "simple",
            "interval": "5s"
        },
        "workflow": {
            "collect": {
                "metrics": {
                    "/intel/linux/cgroups/blkio_stats/io_serviced_recursive/0/minor": {},
                    "/intel/linux/cgroups/blkio_stats/io_serviced_recursive/0/op": {},
                    "/intel/linux/cgroups/blkio_stats/io_serviced_recursive/0/value": {},
                    "/intel/linux/cgroups/blkio_stats/io_serviced_recursive/1/major": {},
                    "/intel/linux/cgroups/blkio_stats/io_serviced_recursive/1/minor": {},
                    "/intel/linux/cgroups/blkio_stats/io_serviced_recursive/1/op": {},
                    "/intel/linux/cgroups/blkio_stats/io_serviced_recursive/1/value": {},
                    "/intel/linux/cgroups/cpu_stats/cpu_usage/total_usage": {},
                    "/intel/linux/cgroups/cpu_stats/cpu_usage/usage_in_kernelmode": {},
                    "/intel/linux/cgroups/cpu_stats/cpu_usage/usage_in_usermode": {},
                    "/intel/linux/cgroups/cpu_stats/throttling_data/periods": {},
                    "/intel/linux/cgroups/cpu_stats/throttling_data/throttled_periods": {},
                    "/intel/linux/cgroups/cpu_stats/throttling_data/throttled_time": {},
                    "/intel/linux/cgroups/hugetlb_stats/2MB/failcnt": {},
                    "/intel/linux/cgroups/hugetlb_stats/2MB/max_usage": {},
                    "/intel/linux/cgroups/hugetlb_stats/2MB/usage": {}
                    "/intel/linux/cgroups/memory_stats/cache": {},
                    "/intel/linux/cgroups/memory_stats/kernel_tcp_usage/failcnt": {},
                    "/intel/linux/cgroups/memory_stats/kernel_tcp_usage/limit": {},
                    "/intel/linux/cgroups/memory_stats/kernel_tcp_usage/max_usage": {},
                    "/intel/linux/cgroups/memory_stats/kernel_tcp_usage/usage": {},
                    "/intel/linux/cgroups/memory_stats/kernel_usage/failcnt": {},
                    "/intel/linux/cgroups/memory_stats/kernel_usage/limit": {},
                    "/intel/linux/cgroups/memory_stats/kernel_usage/max_usage": {},
                    "/intel/linux/cgroups/memory_stats/kernel_usage/usage": {},
                    "/intel/linux/cgroups/pids_stats/init_scope/current": {},
                    "/intel/linux/cgroups/pids_stats/init_scope/limit": {},
                    "/intel/linux/cgroups/pids_stats/system_slice/current": {},
                    "/intel/linux/cgroups/pids_stats/system_slice/limit": {},
                }
            }
        }
    }

This task manifest collects some stats from all cgroup subsystems (blkio, cpu, hugetlb, memory, pids) continously with 5 second interval.

#### How to start it?
If you've got snap properly installed, you need to start snap daemon with root permissions and plugin trust mode:

    $ sudo snapd -t 0

Then in another terminal window you need to load snap plugin using snapctl tool:

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

    NAMESPACE                                                                DATA                    TIMESTAMP
    /intel/linux/cgroups/blkio_stats/io_serviced_recursive/0/minor           0                       2016-05-12 11:50:35.020103525 +0200 CEST
    /intel/linux/cgroups/blkio_stats/io_serviced_recursive/0/op              Read                    2016-05-12 11:50:35.107393787 +0200 CEST
    /intel/linux/cgroups/blkio_stats/io_serviced_recursive/0/value           0                       2016-05-12 11:50:35.038974999 +0200 CEST
    /intel/linux/cgroups/blkio_stats/io_serviced_recursive/1/major           11                      2016-05-12 11:50:35.075105042 +0200 CEST
    /intel/linux/cgroups/blkio_stats/io_serviced_recursive/1/minor           0                       2016-05-12 11:50:34.992783854 +0200 CEST
    /intel/linux/cgroups/blkio_stats/io_serviced_recursive/1/op              Write                   2016-05-12 11:50:35.087683987 +0200 CEST
    /intel/linux/cgroups/blkio_stats/io_serviced_recursive/1/value           0                       2016-05-12 11:50:35.05338193 +0200 CEST
    /intel/linux/cgroups/cpu_stats/cpu_usage/total_usage                     0                       2016-05-12 11:50:34.999210413 +0200 CEST
    /intel/linux/cgroups/cpu_stats/cpu_usage/usage_in_kernelmode             0                       2016-05-12 11:50:35.021984412 +0200 CEST
    /intel/linux/cgroups/cpu_stats/cpu_usage/usage_in_usermode               0                       2016-05-12 11:50:35.001145321 +0200 CEST
    /intel/linux/cgroups/cpu_stats/throttling_data/periods                   0                       2016-05-12 11:50:35.024332746 +0200 CEST
    /intel/linux/cgroups/cpu_stats/throttling_data/throttled_periods         0                       2016-05-12 11:50:35.003082448 +0200 CEST
    /intel/linux/cgroups/cpu_stats/throttling_data/throttled_time            0                       2016-05-12 11:50:35.091528061 +0200 CEST
    /intel/linux/cgroups/hugetlb_stats/2MB/failcnt                           0                       2016-05-12 11:50:35.058841571 +0200 CEST
    /intel/linux/cgroups/hugetlb_stats/2MB/max_usage                         0                       2016-05-12 11:50:35.093596618 +0200 CEST
    /intel/linux/cgroups/hugetlb_stats/2MB/usage                             0                       2016-05-12 11:50:35.07846079 +0200 CEST
    /intel/linux/cgroups/memory_stats/cache                                  7.3334784e+07           2016-05-12 11:50:35.060562164 +0200 CEST
    /intel/linux/cgroups/memory_stats/kernel_tcp_usage/failcnt               0                       2016-05-12 11:50:34.998067249 +0200 CEST
    /intel/linux/cgroups/memory_stats/kernel_tcp_usage/limit                 2.251799813685247e+15   2016-05-12 11:50:35.110497746 +0200 CEST
    /intel/linux/cgroups/memory_stats/kernel_tcp_usage/max_usage             0                       2016-05-12 11:50:35.006567079 +0200 CEST
    /intel/linux/cgroups/memory_stats/kernel_tcp_usage/usage                 24576                   2016-05-12 11:50:35.026982888 +0200 CEST
    /intel/linux/cgroups/memory_stats/kernel_usage/failcnt                   0                       2016-05-12 11:50:35.112223307 +0200 CEST
    /intel/linux/cgroups/memory_stats/kernel_usage/limit                     9.223372036854772e+18   2016-05-12 11:50:35.095054631 +0200 CEST
    /intel/linux/cgroups/memory_stats/kernel_usage/max_usage                 0                       2016-05-12 11:50:35.096889677 +0200 CEST
    /intel/linux/cgroups/memory_stats/kernel_usage/usage                     0                       2016-05-12 11:50:35.062796217 +0200 CEST
    /intel/linux/cgroups/pids_stats/init_scope/current                       1                       2016-05-12 11:50:35.079772278 +0200 CEST
    /intel/linux/cgroups/pids_stats/init_scope/limit                         512                     2016-05-12 11:50:35.098615219 +0200 CEST
    /intel/linux/cgroups/pids_stats/system_slice/current                     130                     2016-05-12 11:50:35.024025875 +0200 CEST
    /intel/linux/cgroups/pids_stats/system_slice/limit                       0                       2016-05-12 11:50:35.064794947 +0200 CEST

Now if your task works and returns data, you can get the data via REST API:

    curl -L http://localhost:8181/v1/tasks/\<task id\>

That's it!

### Roadmap
There isn't a current roadmap for this plugin. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Authors

* [Mateusz Kleina](https://github.com/mkleina)
* [Marcin Olszewski](https://github.com/marcintao)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.