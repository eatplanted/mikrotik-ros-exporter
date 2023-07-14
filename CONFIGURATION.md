# Mikrotik RouterOS Exporter configuration

The file is written in YAML format defined by the scheme described below.

Generic placeholders are defined as follows:

* <string>: a regular string.

See [example.yml](examples/config.yml) for configuration examples.

```yaml
credentials:
  [ <string>: <credential> ... ]
```

## `<credential>`

```yaml
username: <string>
password: <string>
```
