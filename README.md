[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/sensu/sensu-enterprise-metrics)
[ ![Build Status](https://travis-ci.org/sensu/sensu-enterprise-metrics.svg?branch=master)](https://travis-ci.org/sensu/sensu-enterprise-metrics)

# Sensu Enterprise Metrics Collector

- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Sensu Go](#sensu-go)
    - [Asset registration](#asset-registration)
    - [Asset definition](#asset-definition)
    - [Check definition](#check-definition)
  - [Sensu Core](#sensu-core)
    - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

This check plugin collects Sensu Enterprise (Classic) metrics.

## Usage examples

### Help

```
$ sensu-enterprise-metrics -h
Usage of sensu-enterprise-metrics:
  -ca-cert string
    	Path to CA certificate
  -host string
    	Sensu Enterprise API host. (default "localhost")
  -insecure-skip-verify
    	Don't verify TLS hostnames
  -latest
    	Only return the latest point per Enterprise metric.
  -password string
    	Sensu Enterprise API password.
  -port int
    	Sensu Enterprise API port. (default 4567)
  -scheme string
    	Sensu Enterprise URL scheme (http or https) (default "http")
  -timeout int
    	Sensu Enterprise API request timeout (in seconds). (default 15)
  -user string
    	Sensu Enterprise API user.
```

## Configuration
### Sensu Go
#### Asset registration

Assets are the best way to make use of this plugin. If you're not using an asset, please consider doing so! If you're using sensuctl 5.13 or later, you can use the following command to add the asset: 

`sensuctl asset add sensu/sensu-enterprise-metrics`

If you're using an earlier version of sensuctl, you can download the asset definition from [this project's Bonsai asset index page](https://bonsai.sensu.io/assets/sensu/sensu-enterprise-metrics).

To install from the local path of the sensu-enterprise-metrics repository:

`go build -o /usr/local/bin/sensu-enterprise-metrics main.go`

#### Asset definition

```yaml
---
type: Asset
api_version: core/v2
metadata:
  name: sensu-enterprise-metrics
  labels: 
  annotations:
    io.sensu.bonsai.url: https://bonsai.sensu.io/assets/sensu/sensu-enterprise-metrics
    io.sensu.bonsai.api_url: https://bonsai.sensu.io/api/v1/assets/sensu/sensu-enterprise-metrics
    io.sensu.bonsai.tier: Community
    io.sensu.bonsai.version: 0.0.1
    io.sensu.bonsai.namespace: sensu
    io.sensu.bonsai.name: sensu-enterprise-metrics
    io.sensu.bonsai.tags: ''
spec:
  builds:
  - url: https://assets.bonsai.sensu.io/6968c1b55fd42b61cb91f3ec7e83e68e69923955/sensu-enterprise-metrics_0.0.1_linux_amd64.tar.gz
    sha512: ec30a68bcf1ecd50e7d97520916519c95dfec5135c8cfa8c105cf26e88e2a213062e4c7553e3895f5915cebb6a64e392f18e373abac4901d7096eb53743c3cf7
    filters:
    - entity.system.os == 'linux'
    - entity.system.arch == 'amd64'
```

#### Check definition
```yaml

---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-enterprise-metrics
spec:
  command: sensu-enterprise-metrics -user username -password secret -latest
  output_metric_format: graphite_plaintext
  output_metric_handlers:
  - influxdb
  runtime_assets:
  - sensu/sensu-prometheus-collector
  subscriptions:
  - sensu-enterprise
  publish: true
  interval: 10
```

### Sensu Core

#### Check definition
```json
{
  "api_version": "core/v2",
  "type": "CheckConfig",
  "metadata": {
    "namespace": "default",
    "name": "sensu-enterprise-metrics"
  },
  "spec": {
    "command": "sensu-enterprise-metrics -user username -password secret -latest",
    "subscriptions": [
      "sensu-enterprise"
    ],
    "publish": true,
    "interval": 10,
    "output_metric_format": "graphite_plaintext",
    "output_metric_handlers": [
      "influxdb"
    ]
  }
}
```

## Installation from source

### Sensu Go

See the instructions above for [asset registration](#asset-registration).

### Sensu Core

Install and setup plugins on [Sensu Core](https://docs.sensu.io/sensu-core/latest/installation/installing-plugins/).

## Additional notes

None.

## Contributing

See [CONTRIBUTING.md](https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md) for information about contributing to this plugin.
