# go-srtm
Shuttle Radar Topography Mission (SRTM) is an international research effort to create a global digital elevation dataset. Since the data files are large I created a small web service in Go so applications can access the data without having the users download the whole database locally.

`go-srtm` only works with raw one arc-second data files (.hgt, 3601x3601). Each file should be 25,934,402 bytes.

## How to run

## API Documentation

### /elevation
This endpoint returns the elevation value for a set of coordinates

**Parameters**:
 * `locations`: Where to look for elevation. Each point is two float numbers (latitude and longitude) separated with
                a comma (`,`). You can query multiple locations in one call by separating the points with `|`, for
                example `/elevation?locations=21.124,37.157|22.984,23.487` will return elevation for (21.124, 37.157)
                and (22.984, 23.487).

**Returns:**: JSON Array of the requested data in meters. If data is not available at a requested point it will return
              `-1`.

### /elevationPath
Returns elevation data every 30m on a path connecting two points.

**Parameters**:
 * `from`: Where does the path start (latitude and longitude, comma separated)
 * `to`: Where does the path end (latitude and longitude, comma separated)
 
**Returns**: JSON Array of elevation data.

### /tiles
Returns the list of available SRTM tiles

**Returns**: JSON Array of available SRTM files (string, no extension). For example: `["N41E021", "N41E022"]`


## Configuration
`go-srtm` will look for `config.ini` in the current path. A sample configuration file:
```ini
address = ""
port = "8080"
dataset = "/srv/srtm"
```

`dataset` should point to the directory containing the SRTM data. You can use `address` to set which IP the service
is listening on.

## License
`go-srtm` is available under the MIT license.
