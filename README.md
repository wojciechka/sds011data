# Sds011 data

Set of code and utilities to read air particle metrics from [SDS011](https://aqicn.org/sensor/sds011/) sensor.

The code includes a utility that writes raw data from the sensor in a set of files.

The code also includes a tool to upload all data to an InfluxDB instance, possibly also grouping it by specific duration.

## Installing

```bash
go install github.com/wojciechka/sds011data/cmd/sds011-reader
go install github.com/wojciechka/sds011data/cmd/sds011-influxdb-writer
```

## Usage

### Retrieving data from sensor

```bash
$ sds011-reader \
  -device /dev/ttyUSB0 \
	-dataDirectory /path/to/datadir
```

The command above will write data to specified directory, using the file format described below.

The directory `/path/to/datadir` must already exist.

### Uploading data to InfluxDB

```bash
$ sds011-influxdb-writer \
  -dataDirectory /path/to/datadir \
	-stateFile /path/to/datadir/.writer-state \
  -influxToken "(token)" \
  -influxOrg "(org)" \
  -influxBucket "(bucket)" \
  -groupBy "1m" \
  -influxTags "tag1=value1,tag2=value2,..." \
  -wait
```

## File format

The data read from sensor is kept in a directory, using flat file structure, grouped in multiple files.

Data for each specific day is stored using `YYYYMMDD` filename format, such as `20060102`.

Each file keeps data for a specific second as 4 bytes. The format of the data matches the format returned by [SDS011](https://aqicn.org/sensor/sds011/) sensor. The first two bytes contain value of PM2.5 readout in low endian notation, multiplied by 10. The last two bytes contain value of PM10 readout in low endian notation, multiplied by 10.

For example for value of `04 01 fe 00`, `0x0104` indicates current value of PM2.5 - the value being `260`, therefore the real value for PM2.5 result is `26.0`. The raw value for PM10 is `0x00fe`, which is `254` in decimal notation. The real value of the readout for PM10 is `25.4` in this case.

Offset to each data for each second is 4 multiplied by number of seconds elapsed since midnight. An appropriate formula is `offset.Second() + offset.Minute()*60 + offset.Hour()*3600)*4`.

Any data that is missing is filled with the four bytes being `0xffffffff` - as in certain cases a readout of `0.0` for both PM2.5 and PM10 is possible, but values of `6553.5` are very unlikely. Also, as that level would exceed limits by around 100 times, the concern at that point should not be storing of the data, but implications of the values being so high.
