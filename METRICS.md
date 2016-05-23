# snap-plugin-collector-cgroups

## Available metrics:

**Info:** First asterisk after stats category is mountpoint path element.

**For example:** /intel/linux/cgroups/blkio_stats/\*/io_merged_recursive/\*/majo  
**Can be:** /intel/linux/cgroups/blkio_stats/**system\_slice**/io_merged_recursive/\*/majo

| Namespace | Description | Unit |
|-----------|-------------|------|
|/intel/linux/cgroups/blkio_stats/\*/io_merged_recursive/\*/major|Number of bios/requests merged into requests belonging to all the descendant cgroups (major)||
|/intel/linux/cgroups/blkio_stats/\*/io_merged_recursive/\*/minor|Number of bios/requests merged into requests belonging to all the descendant cgroups (minor)||
|/intel/linux/cgroups/blkio_stats/\*/io_merged_recursive/\*/op|Number of bios/requests merged into requests belonging to all the descendant cgroups (operation)||
|/intel/linux/cgroups/blkio_stats/\*/io_merged_recursive/\*/value|Number of bios/requests merged into requests belonging to all the descendant cgroups (value)||
|/intel/linux/cgroups/blkio_stats/\*/io_queue_recursive/\*/major|Number of requests queued up at any given instant from all the descendant cgroups (major)||
|/intel/linux/cgroups/blkio_stats/\*/io_queue_recursive/\*/minor|Number of requests queued up at any given instant from all the descendant cgroups (minor)||
|/intel/linux/cgroups/blkio_stats/\*/io_queue_recursive/\*/op|Number of requests queued up at any given instant from all the descendant cgroups (operation)||
|/intel/linux/cgroups/blkio_stats/\*/io_queue_recursive/\*/value|Number of requests queued up at any given instant from all the descendant cgroups (value)||
|/intel/linux/cgroups/blkio_stats/\*/io_service_bytes_recursive/\*/major|Number of bytes transferred to/from the disk from all the descendant cgroups (major)|B
|/intel/linux/cgroups/blkio_stats/\*/io_service_bytes_recursive/\*/minor|Number of bytes transferred to/from the disk from all the descendant cgroups (minor)|B
|/intel/linux/cgroups/blkio_stats/\*/io_service_bytes_recursive/\*/op|Number of bytes transferred to/from the disk from all the descendant cgroups (operation)||
|/intel/linux/cgroups/blkio_stats/\*/io_service_bytes_recursive/\*/value|Number of bytes transferred to/from the disk from all the descendant cgroups (value)||
|/intel/linux/cgroups/blkio_stats/\*/io_service_time_recursive/\*/major|Amount of time between request dispatch and request completion from all the descendant cgroups (major)|ns
|/intel/linux/cgroups/blkio_stats/\*/io_service_time_recursive/\*/minor|Amount of time between request dispatch and request completion from all the descendant cgroups (minor)|ns
|/intel/linux/cgroups/blkio_stats/\*/io_service_time_recursive/\*/op|Amount of time between request dispatch and request completion from all the descendant cgroups (operation)||
|/intel/linux/cgroups/blkio_stats/\*/io_service_time_recursive/\*/value|Amount of time between request dispatch and request completion from all the descendant cgroups (value)||
|/intel/linux/cgroups/blkio_stats/\*/io_serviced_recursive/\*/major|Total number of block I/O requests serviced in that container (major)||
|/intel/linux/cgroups/blkio_stats/\*/io_serviced_recursive/\*/minor|Total number of block I/O requests serviced in that container (minor)||
|/intel/linux/cgroups/blkio_stats/\*/io_serviced_recursive/\*/op|Total number of block I/O requests serviced in that container (operation)||
|/intel/linux/cgroups/blkio_stats/\*/io_serviced_recursive/\*/value|Total number of block I/O requests serviced in that container (value)||
|/intel/linux/cgroups/blkio_stats/\*/io_time_recursive/\*/major|Disk time allocated to all devices from all the descendant cgroups (major)|ns
|/intel/linux/cgroups/blkio_stats/\*/io_time_recursive/\*/minor|Disk time allocated to all devices from all the descendant cgroups (minor)|ns
|/intel/linux/cgroups/blkio_stats/\*/io_time_recursive/\*/op|Disk time allocated to all devices from all the descendant cgroups (operation)||
|/intel/linux/cgroups/blkio_stats/\*/io_time_recursive/\*/value|Disk time allocated to all devices from all the descendant cgroups (value)||
|/intel/linux/cgroups/blkio_stats/\*/io_wait_time_recursive/\*/major|Amount of time the IOs for this cgroup spent waiting in the scheduler queues for service from all the descendant cgroups (major)|ns
|/intel/linux/cgroups/blkio_stats/\*/io_wait_time_recursive/\*/minor|Amount of time the IOs for this cgroup spent waiting in the scheduler queues for service from all the descendant cgroups (minor)|ns
|/intel/linux/cgroups/blkio_stats/\*/io_wait_time_recursive/\*/op|Amount of time the IOs for this cgroup spent waiting in the scheduler queues for service from all the descendant cgroups (operation)||
|/intel/linux/cgroups/blkio_stats/\*/io_wait_time_recursive/\*/value|Amount of time the IOs for this cgroup spent waiting in the scheduler queues for service from all the descendant cgroups (value)||
|/intel/linux/cgroups/blkio_stats/\*/sectors_recursive/\*/major|Number of sectors transferred to/from disk bys from all the descendant cgroups (major)||
|/intel/linux/cgroups/blkio_stats/\*/sectors_recursive/\*/minor|Number of sectors transferred to/from disk bys from all the descendant cgroups (minor)||
|/intel/linux/cgroups/blkio_stats/\*/sectors_recursive/\*/op|Number of sectors transferred to/from disk bys from all the descendant cgroups (operation)||
|/intel/linux/cgroups/blkio_stats/\*/sectors_recursive/\*/value|Number of sectors transferred to/from disk bys from all the descendant cgroups (value)||
|/intel/linux/cgroups/cpu_stats/\*/cpu_usage/percpu_usage/\*/value|CPU time consumed on each CPU by all tasks|ns
|/intel/linux/cgroups/cpu_stats/\*/cpu_usage/total_usage|Total CPU time consumed|ns
|/intel/linux/cgroups/cpu_stats/\*/cpu_usage/usage_in_kernelmode|CPU time consumed by tasks in system (kernel) mode|ns
|/intel/linux/cgroups/cpu_stats/\*/cpu_usage/usage_in_usermode|CPU time consumed by tasks in user mode|ns
|/intel/linux/cgroups/cpu_stats/\*/throttling_data/periods|Number of period intervals that have elapsed||
|/intel/linux/cgroups/cpu_stats/\*/throttling_data/throttled_periods|Number of times tasks in a cgroup have been throttled||
|/intel/linux/cgroups/cpu_stats/\*/throttling_data/throttled_time|Total time duration for which tasks in a cgroup have been throttled|ns
|/intel/linux/cgroups/cpuacct_stats/\*/cpu_usage/percpu_usage/\*/value|CPU time consumed on each CPU by all tasks|ns
|/intel/linux/cgroups/cpuacct_stats/\*/cpu_usage/total_usage|Total CPU time consumed|ns
|/intel/linux/cgroups/cpuacct_stats/\*/cpu_usage/usage_in_kernelmode|CPU time consumed by tasks in system (kernel) mode|ns
|/intel/linux/cgroups/cpuacct_stats/\*/cpu_usage/usage_in_usermode|CPU time consumed by tasks in user mode|ns
|/intel/linux/cgroups/cpuacct_stats/\*/throttling_data/periods|Number of period intervals that have elapsed||
|/intel/linux/cgroups/cpuacct_stats/\*/throttling_data/throttled_periods|Number of times tasks in a cgroup have been throttled||
|/intel/linux/cgroups/cpuacct_stats/\*/throttling_data/throttled_time|Total time duration for which tasks in a cgroup have been throttled|ns
|/intel/linux/cgroups/hugetlb_stats/\*/\*/failcnt|Report the number of allocation failure due to HugeTLB limit||
|/intel/linux/cgroups/hugetlb_stats/\*/\*/max_usage|Report max "hugepagesize" HugeTLB usage recorded|B
|/intel/linux/cgroups/hugetlb_stats/\*/\*/usage|Report current usage for "hugepagesize" HugeTLB|B
|/intel/linux/cgroups/memory_stats/\*/cache|Page cache including tmpfs|B
|/intel/linux/cgroups/memory_stats/\*/kernel_tcp_usage/failcnt|Reports the number of times that the kernel TCP stack memory limit has reached the value set in memory.limit_in_bytes||
|/intel/linux/cgroups/memory_stats/\*/kernel_tcp_usage/limit|Memory usage limit for kernel TCP stack|B
|/intel/linux/cgroups/memory_stats/\*/kernel_tcp_usage/max_usage|Maximum memory usage by kernel TCP stack reported|B
|/intel/linux/cgroups/memory_stats/\*/kernel_tcp_usage/usage|Memory usage by kernel TCP stack|B
|/intel/linux/cgroups/memory_stats/\*/kernel_usage/failcnt|Reports the number of times that the kernel memory limit has reached the value set in memory.limit_in_bytes||
|/intel/linux/cgroups/memory_stats/\*/kernel_usage/limit|Memory usage limit for kernel|B
|/intel/linux/cgroups/memory_stats/\*/kernel_usage/max_usage|Maximum memory usage by kernel reported|B
|/intel/linux/cgroups/memory_stats/\*/kernel_usage/usage|Memory usage by kernel|B
|/intel/linux/cgroups/memory_stats/\*/stats/active_anon|Number of bytes of anonymous and swap cache memory on active LRU list|B
|/intel/linux/cgroups/memory_stats/\*/stats/active_file|Number of bytes of file-backed memory on active LRU list|B
|/intel/linux/cgroups/memory_stats/\*/stats/cache|Number of bytes of page cache memory|B
|/intel/linux/cgroups/memory_stats/\*/stats/dirty|Number of bytes that are waiting to get written back to the disk|B
|/intel/linux/cgroups/memory_stats/\*/stats/hierarchical_memory_limit|Number of bytes of memory limit with regard to hierarchy under which the memory cgroup is|B
|/intel/linux/cgroups/memory_stats/\*/stats/inactive_anon|Amount of buffer or page cache memory that are free and available|B
|/intel/linux/cgroups/memory_stats/\*/stats/inactive_file|Number of bytes of file-backed memory on inactive LRU list|B
|/intel/linux/cgroups/memory_stats/\*/stats/mapped_file|Number of bytes of mapped file (includes tmpfs/shmem)|B
|/intel/linux/cgroups/memory_stats/\*/stats/pgfault|Number of page faults per second||
|/intel/linux/cgroups/memory_stats/\*/stats/pgmajfault|Number of major page faults per second||
|/intel/linux/cgroups/memory_stats/\*/stats/pgpgin|Number of charging events to the memory cgroup||
|/intel/linux/cgroups/memory_stats/\*/stats/pgpgout|Number of uncharging events to the memory cgroup||
|/intel/linux/cgroups/memory_stats/\*/stats/rss|Number of bytes of anonymous and swap cache memory|B
|/intel/linux/cgroups/memory_stats/\*/stats/rss_huge|Number of bytes of anonymous transparent hugepages|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_active_anon|Total number of bytes of anonymous and swap cache memory on active LRU list in all children's active_anon|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_active_file|Total number of bytes of file-backed memory on active LRU list in all children's active_file|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_cache|Total number of bytes of page cache memory in all children's cache|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_dirty|Total number of bytes that are waiting to get written back to the disk in all children's dirty|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_inactive_anon|The total amount of buffer or page cache memory that are free and available in all children's total_inactive|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_inactive_file|Total number of bytes of file-backed memory on inactive LRU list in all children's inactive_file|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_mapped_file|Total number of bytes of mapped file (includes tmpfs/shmem) in all children's mapped_file|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_pgfault|Total number of page faults per second in all children's pgfault||
|/intel/linux/cgroups/memory_stats/\*/stats/total_pgmajfault|Total number of page major faults which happened since the creation of the cgroup in all children's pgmajfault||
|/intel/linux/cgroups/memory_stats/\*/stats/total_pgpgin|Total number of charging events to the memory cgroup in all children's pgpin||
|/intel/linux/cgroups/memory_stats/\*/stats/total_pgpgout|Number of uncharging events to the memory cgroup in all children's pgpout||
|/intel/linux/cgroups/memory_stats/\*/stats/total_rss|Total number of bytes of anonymous and swap cache memory in all children's rss|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_rss_huge|Total number of bytes of anonymous transparent hugepages in all children's rss_huge|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_unevictable|Total number of bytes of memory that cannot be reclaimed (mlocked etc) in all children's unevictable|B
|/intel/linux/cgroups/memory_stats/\*/stats/total_writeback|Total number of bytes of file/anon cache that are queued for syncing to disk in all children's writeback|B
|/intel/linux/cgroups/memory_stats/\*/stats/unevictable|Number of bytes of memory that cannot be reclaimed (mlocked etc)|B
|/intel/linux/cgroups/memory_stats/\*/stats/writeback|Number of bytes of file/anon cache that are queued for syncing to disk|B
|/intel/linux/cgroups/memory_stats/\*/swap_usage/failcnt|Reports the number of times the swap space limit has reached the value set in memorysw.limit_in_bytes||
|/intel/linux/cgroups/memory_stats/\*/swap_usage/limit|Reports the limit of swap space usage by porcesses in the cgroup||
|/intel/linux/cgroups/memory_stats/\*/swap_usage/max_usage|Reports the maximum swap space used by processes in the cgroup||
|/intel/linux/cgroups/memory_stats/\*/swap_usage/usage|Reports the total swap space usage by processes in the cgroup||
|/intel/linux/cgroups/memory_stats/\*/usage/failcnt|Reports the number of times that the memory limit has reached the value set in memory.limit_in_bytes||
|/intel/linux/cgroups/memory_stats/\*/usage/limit|Memory usage limit for processes in the cgroup|B
|/intel/linux/cgroups/memory_stats/\*/usage/max_usage|Maximum memory usage by processes in the cgroup reported|B
|/intel/linux/cgroups/memory_stats/\*/usage/usage|Memory usage by processes in the cgroup|B
|/intel/linux/cgroups/pids_stats/\*/current|Number of processes currently in the cgroup||
|/intel/linux/cgroups/pids_stats/\*/limit|Limit of processes currently in the cgroup||
