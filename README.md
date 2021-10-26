# Weather query
Short and sweet weather API query tool written in Go. It employs both GeoNorge (Kartverket) and YR weather APIs in order to allow fast weather lookup using only a recognizable place name in the Kingdom of Norway.

```
  Usage: ./weather_query <place name>
```

Short and sweet.



### TODO
Handle errors gracefully, e.g. when I query non-existent or unknown place name:
```
  $ go run weather_query.go Tolkien
  panic: runtime error: index out of range [0] with length 0
```
A useful error message is better than a panic.
