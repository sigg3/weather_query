# Weather query
Short and sweet weather API query tool written in Go. It employs both GeoNorge (Kartverket) and YR weather APIs in order to allow fast weather lookup using only a recognizable place name in the Kingdom of Norway.

```
  Usage: ./weather_query <place name>
```

for example:

```
  $ go run weather_query.go Tønsberg
  Current weather in:           Tønsberg (long: 10.40764, lat: 59.26751)
  Query time:                   2021-10-26 12:00:00 +0000 UTC
  Temperature (celsius):        10.9°C
  Temperature (fahrenheit):     42.9°F
  Airpressure at sea level:     1007.5 hPa
  Current humidity:             70.3 %
  Wind speed:                   3.5 m/sec
```

Short and sweet.



### TODO
Handle errors gracefully, e.g. when I query non-existent or unknown place name:
```
  $ go run weather_query.go Tolkien
  panic: runtime error: index out of range [0] with length 0
```
A useful error message is better than a panic.
