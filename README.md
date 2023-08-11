# windows_exporter

![Build Status](https://github.com/prometheus-community/windows_exporter/workflows/windows_exporter%20CI/CD/badge.svg)

A Prometheus exporter for Windows machines.

**0.22.3 / 2023-08-11**

**[ADD] 增加对winlogbeat采集事件状态的监控**

- HELP windows_winlogbeat_published_events Winlogbeat monitoring log total published events.
- TYPE windows_winlogbeat_published_events counter
- windows_winlogbeat_published_events 25


**0.22.2 / 2023-07-27**

**[BUGfix] 优化获取filebeat openfiles信息的方式**


**0.22.1 / 2023-07-25**

**[ADD] 增加对filebeat进程状态的监控**

- HELP node_filebeat_up Value is 1 if filebeat process is 'up', 0 otherwise.
- TYPE node_filebeat_up gauge
- windows_process_status_up{process_name="filebeat"} 1

**[ADD] 增加对filebeat打开文件状态的监控**

- HELP windows_filebeat_openfiles Filebeat monitoring log harvester openfiles running.
- TYPE windows_filebeat_openfiles counter
- windows_filebeat_openfiles 2


**[ADD] 增加对winlogbeat进程状态的监控**

- HELP node_rsyslog_up Value is 1 if rsyslog process is 'up', 0 otherwise.
- TYPE node_rsyslog_up gauge
- windows_process_status_up{process_name="winlogbeat"} 1


**[DEL] 移除以go_开头的33个监控项,具体如下：**

- go_gc_duration_seconds{quantile="0"}
- go_gc_duration_seconds{quantile="0.25"}
- go_gc_duration_seconds{quantile="0.5"}
- go_gc_duration_seconds{quantile="0.75"}
- go_gc_duration_seconds{quantile="1"}
- go_gc_duration_seconds_sum
- go_gc_duration_seconds_count
- go_goroutines
- go_info{version="go1.19.1"}
- go_memstats_alloc_bytes
- go_memstats_alloc_bytes_total
- go_memstats_buck_hash_sys_bytes
- go_memstats_frees_total
- go_memstats_gc_sys_bytes
- go_memstats_heap_alloc_bytes
- go_memstats_heap_idle_bytes
- go_memstats_heap_inuse_bytes
- go_memstats_heap_objects
- go_memstats_heap_released_bytes
- go_memstats_heap_sys_bytes
- go_memstats_last_gc_time_seconds
- go_memstats_lookups_total
- go_memstats_mallocs_total
- go_memstats_mcache_inuse_bytes
- go_memstats_mcache_sys_bytes
- go_memstats_mspan_inuse_bytes
- go_memstats_mspan_sys_bytes
- go_memstats_next_gc_bytes
- go_memstats_other_sys_bytes
- go_memstats_stack_inuse_bytes
- go_memstats_stack_sys_bytes
- go_memstats_sys_bytes
- go_threads


## Flags

windows_exporter accepts flags to configure certain behaviours. The ones configuring the global behaviour of the exporter are listed below, while collector-specific ones are documented in the respective collector documentation above.

Flag     | Description | Default value
---------|-------------|--------------------
`--web.listen-address` | host:port for exporter. | `:9182`
`--telemetry.path` | URL path for surfacing collected metrics. | `/metrics`
`--telemetry.max-requests` | Maximum number of concurrent requests. 0 to disable. | `5`
`--collectors.enabled` | Comma-separated list of collectors to use. Use `[defaults]` as a placeholder which gets expanded containing all the collectors enabled by default." | `[defaults]`
`--collectors.print` | If true, print available collectors and exit. |
`--scrape.timeout-margin` | Seconds to subtract from the timeout allowed by the client. Tune to allow for overhead or high loads. | `0.5`
`--web.config.file` | A [web config][web_config] for setting up TLS and Auth | None


## Supported versions

windows_exporter supports Windows Server versions 2008R2 and later, and desktop Windows version 7 and later.

## Usage

promu.exe build -v   
candle.exe -nologo -arch "x64" -ext WixFirewallExtension -ext WixUtilExtension -out .\windows_exporter.wixobj -dVersion="0.22.1" .\windows_exporter.wxs
light.exe -nologo -spdb -ext WixFirewallExtension -ext WixUtilExtension -out ".\windows_exporter-0.22.1-amd64.msi" .\windows_exporter.wixobj

The prometheus metrics will be exposed on [localhost:9182](http://localhost:9182)


### Using a configuration file

YAML configuration files can be specified with the `--config.file` flag. e.g. `.\windows_exporter.exe --config.file=config.yml`. If you are using the absolute path, make sure to quote the path, e.g. `.\windows_exporter.exe --config.file="C:\Program Files\windows_exporter\config.yml"`

```yaml
collectors:
  enabled: cpu,cs,net,service
collector:
  service:
    services-where: "Name='windows_exporter'"
log:
  level: warn
```

An example configuration file can be found [here](docs/example_config.yml).


