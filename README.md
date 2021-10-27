# Weather query
Short and sweet weather API query tool written in Go. It employs both GeoNorge (Kartverket) and YR weather APIs in order to allow fast weather lookup using only a recognizable place name in the Kingdom of Norway.


## Build
To compile a binary simply run:
```
  go build weather_query.go
```
The compiled version is notably faster (0,168s vs 0,566s). There are also semi-current binaries for Linux, Mac and Windows 64-bit [in the /dist folder](https://github.com/sigg3/weather_query/tree/main/dist).

## Usage

```
  $ ./weather_query 
  Usage: ./weather_query <place>
  Error: Location string argument missing. Use any place in Norway.
```

for example:

```
  $ ./weather_query Tønsberg
  Current weather in:           Tønsberg
  Coordinates (long, lat):      10.40764, 59.26751
  Observation timestamp:        2021-10-26 18:00:00 +0000 UTC
  Temperature (celsius):        7.3°C
  Temperature (fahrenheit):     39.3°F
  Airpressure at sea level:     1009.4 hPa
  Current humidity:             76.3 %
  Wind speed:                   2.8 m/sec
```

Short and sweet. weather_query will pick the first hit if several.



### TODO
* Put repetetive tasks in functions, esp. API query, cf. stackoverflow.com/a/31129967

* Handle errors gracefully, e.g. when I query non-existent or unknown place name:
```
  $ go run weather_query.go Tolkien
  panic: runtime error: index out of range [0] with length 0
```
A useful error message is better than a panic.

* Could add a parenthesis on tstamp (N &lt;time unit> ago)

* Add an array of postal codes

* Use argument flags, to allow for: --usage, --coord (don't lookup place name), --temperature ops.
