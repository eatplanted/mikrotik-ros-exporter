# mikrotik-ros-exporter

A simple Prometheus exporter for Mikrotik devices.

This exporter works similar to the [Blackbox Exporter](https://github.com/prometheus/blackbox_exporter) where you can specify different targets in Prometheus to scrape Mikrotik devices through this exporter.

## Configuration

Mikrotik RouterOS Exporter is configured using [configuration file](CONFIGURATION.md) and command line flags.

The configuration file is used to define the various credentials that can be used to query metrics from a Mikrotik device.

To view all available command line flags, run `./mikrotik-ros-exporter -h`.

To specify which [configuration file](CONFIGURATION.md) to load, use the `-config` flag.

In addition, a [sample configuration](examples/config.yml) is also available.

The timeout is automatically determined from the `scrape_timeout` in the [Prometheus configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#configuration-file), slightly reduced to account for network delays. This can be further constrained by the timeout in the configuration file. If neither is specified, the default value is 120 seconds.

## Mikrotik Configuration

For monitoring purposes, it is recommended to create a new user that has API and read-only access to the Router OS device(s):

```
/user group add name=monitoring policy=api,read
/user add name=monitoring group=monitoring password=changeme
```

## Prometheus Configuration

The Mikrotik RouterOS Exporter implements the multi-target exporter pattern, therefore we recommend reading the guide [Understanding and using the multi-target exporter pattern](https://prometheus.io/docs/guides/multi-target-exporter/) to get an overview of the configuration.

The target must be passed to the Mikrotik RouterOS Exporter as a parameter, this can be done with relabelling.

Example config:

```yaml
scrape_configs:
  - job_name: "mikrotik"
    metrics_path: /probe
    params:
      credential: [default]
      skip_tls_verify: [false]
    static_configs:
      - targets:
        - https://10.0.1.1
        - http://10.0.20.1:8081
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9142 # The Mikrotik RouterOS Exporter's real hostname:port.
```
