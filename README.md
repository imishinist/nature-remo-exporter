# Nature Remo Exporter

Nature Remo Exporter is a Prometheus Exporter for [Nature Remo](https://shop.nature.global/collections/nature-remo).

![Prometheus Exporter Screenshots](/assets/prometheus_exporter.png)

## Getting started

### Creating Access Token

Follow these steps to create an Access Token.

1. Visit the [Nature Remo website](https://home.nature.global).
2. Log in to your Nature Remo account.
3. Click [Generate access token]
4. Copy an Access Token

## Installation

Install from GitHub Releases.

https://github.com/imishinist/nature-remo-exporter/releases

### example

Install to `/tmp/nature-remo-exporter`.

```bash
mkdir -p /tmp/nature-remo-exporter
version=0.1.0 curl -sSL -o- https://github.com/imishinist/nature-remo-exporter/releases/download/v${version}/nature-remo-exporter_Linux_x86_64.tar.gz | tar xzvf - -C /tmp/nature-remo-exporter
```

or install from go install.

```bash
go install github.com/imishinist/nature-remo-exporter@latest
```

## Run

```bash
export REMO_ACCESS_TOKEN=<access token>
nature-remo-exporter --token $REMO_ACCESS_TOKEN
```

## Help

```bash
Nature Remo Exporter is a Prometheus exporter for Nature Remo smart devices.

This tool collects metrics from Nature Remo Cloud API and exposes them in a format
that Prometheus can scrape. It is designed to help monitor and analyze
the performance and data from Nature Remo devices

Usage:
  nature-remo-exporter [flags]

Flags:
  -h, --help                help for nature-remo-exporter
      --interval duration   Interval between metrics refresh (default 30s)
      --port int            Port to listen on (default 9199)
      --token string        Nature Remo access token
```

## Metrics

For details about the available metrics, please refer to the following site:

https://swagger.nature.global/#/default/get_1_devices

| metrics name                   | description              |
|--------------------------------|--------------------------|
| `nature_remo_humidity`         | current humidity         |
| `nature_remo_illumination`     | current illumination     |
| `nature_remo_movement`         | current movement         |
| `nature_remo_movement_counter` | current movement counter |
| `nature_remo_temperature`      | current temperature      |

### Labels

- id
- name
- firmware_version
- mac_address
- bt_mac_address
- serial_number

## Author

- Taisuke Miyazaki ([@imishinist](https://github.com/imishinist))

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
