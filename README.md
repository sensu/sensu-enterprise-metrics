# Sensu Enterprise Metrics Collector
TravisCI: [![TravisCI Build Status](https://travis-ci.org/sensu/sensu-enterprise-metrics.svg?branch=master)](https://travis-ci.org/sensu/sensu-enterprise-metrics)

The Sensu Go check plugin for collecting Sensu Enterprise (Classic) metrics.

## Installation

Download the latest version of the sensu-enterprise-metrics from [releases][1],
or create an executable script from this source.

From the local path of the sensu-enterprise-metrics repository:
```
go build -o /usr/local/bin/sensu-enterprise-metrics main.go
```

## Configuration

Example Sensu Go check definition:

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
        "subscriptions":[
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

## Usage Examples

Help:
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

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/sensu/sensu-enterprise-metrics/releases
