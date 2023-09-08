# ⚠️ Deprecation Notice
This tool has been replaced by the Telegraf's `inputs.internet_speed` input plugin. This is kept here for reference, but isn't actively developed. Refer to the Telegraf instructions below to use that instead

# speedtest-influxdb-go
This is a wrapper for [speedtest-go](https://github.com/showwin/speedtest-go) that writes the results to InfluxDB. Based on [Speedtest-for-InfluxDB-and-Grafana](https://github.com/TomRoush/Speedtest-for-InfluxDB-and-Grafana), but rewritten in Go.

## Replacing with Telegraf
Telegraf's `inputs.internet_speed` plugin replicates the functionality of this tool and is more actively maintained. It's recommended to use that instead. The `inputs.internet_speed` plugin does change the metric names, but if that isn't a deal breaker, you can configure Telegraf to perform the speedtests. Telegraf's `--once` flag can be used to replicate the `delay: 0` behavior of this application.

[inputs.internet_speed Documentation](https://github.com/influxdata/telegraf/blob/master/plugins/inputs/internet_speed/README.md)

## Usage
- Build with `go build -o speedtest-influxdb-go`
- Fill in `config.yaml` and keep it in the same directory as the executable
- Run with `./speedtest-influxdb-go`

## Docker Setup
- Build with `docker build -t tomroush/speedtest-influxdb-go:latest .`
- Fill in `config.yaml`. It will be mounted as volume into the container
- Run with `docker run --rm -v $PWD/config.yaml:/config.yaml tomroush/speedtest-influxdb-go:latest`

## Upgrading
This can be used to replace the earlier [Speedtest-for-InfluxDB-and-Grafana](https://github.com/TomRoush/Speedtest-for-InfluxDB-and-Grafana) application written in Python. The Go version is more stable and uses less resources. To upgrade, the previous application's `config.ini` must be translated into a `config.yaml` file. The config values have identical names and meanings. See [config.yaml](config.yaml) for an example file.
