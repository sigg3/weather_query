/* Fetch the weather from YR using Norwegian place name.
   For usage info, run the script without any arguments.
   Written by Sigge Smelror (C) 2021, GNU GPL v. 3+

   weather_query is free software: you can redistribute it and/or
   modify it under the terms of the GNU General Public License as
   published by the Free Software Foundation, version 3 or newer.

   weather_query is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   URL: <https://www.gnu.org/licenses/gpl-3.0.txt>

   Bugs/Issues: <https://github.com/sigg3/weather_query/issues>
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type GeoNorgeJson struct {
	Data []struct {
		Name        string `json:"skrivemåte"`
		Coordinates struct {
			Long float64 `json:"øst"`
			Lat  float64 `json:"nord"`
		} `json:"representasjonspunkt"`
	} `json:"navn"`
}

type YrJson struct {
	Data struct {
		Timeseries []struct {
			Time time.Time `json:"time"`
			Data struct {
				Instant struct {
					Details struct {
						AirPressureAtSeaLevel float64 `json:"air_pressure_at_sea_level"`
						AirTemperature        float64 `json:"air_temperature"`
						CloudAreaFraction     float64 `json:"cloud_area_fraction"`
						RelativeHumidity      float64 `json:"relative_humidity"`
						WindFromDirection     float64 `json:"wind_from_direction"`
						WindSpeed             float64 `json:"wind_speed"`
					} `json:"details"`
				} `json:"instant"`
			} `json:"data,omitempty"`
		} `json:"timeseries"`
	} `json:"properties"`
}

// // getJson by Connor Peet @ https://stackoverflow.com/a/31129967
// func getJson(url string, target interface{}) error {
// 	fmt.Println(url)
// 	r, err := xClient.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer r.Body.Close()
// 	fmt.Println(json.NewDecoder(r.Body))
// 	return json.NewDecoder(r.Body).Decode(target)
// }

func main() {
	// Get CLI arg (location) or print usage info
	var uLoc string = "oslo"
	if len(os.Args) > 1 {
		uLoc = os.Args[1]
	} else {
		fmt.Println(strings.Join([]string{"Usage:", os.Args[0], "<place>"}, " "))
		fmt.Println("Error: Location string argument missing. Use any place in Norway.")
		os.Exit(1)
	}

	// Set HTTP request timeout (10 secs)
	var xClient = &http.Client{Timeout: 10 * time.Second}

	// create GET request for coordinates
	const COORD_API string = "https://ws.geonorge.no/stedsnavn/v1/navn?sok="
	const COORD_SUFFIX string = "&utkoordsys=4258&treffPerSide=1&side=1"
	coordApiReq := strings.Join([]string{COORD_API, uLoc, COORD_SUFFIX}, "")

	// Send GET request to COORD_API
	coordApiResp, coordApiErr := xClient.Get(coordApiReq)
	if coordApiErr != nil {
		panic(coordApiErr)
	}

	// Convert to text uint8
	coordTxt, readErr := ioutil.ReadAll(coordApiResp.Body)
	if readErr != nil {
		fmt.Println("Read error GeoNorge")
		panic(readErr)
	}

	// Parse into geo using Location struct
	geo := GeoNorgeJson{}
	geoErr := json.Unmarshal(coordTxt, &geo)
	if geoErr != nil {
		fmt.Println("JSON error GeoNorge")
		panic(geoErr)
	}

	// Pull relevant variables
	var geoName string = geo.Data[0].Name
	var geoLong float64 = geo.Data[0].Coordinates.Long
	var geoLat float64 = geo.Data[0].Coordinates.Lat

	// create GET request for weather query
	const YR_API string = "https://api.met.no/weatherapi/locationforecast/2.0/compact.json?"
	longLatStr := fmt.Sprintf("lat=%v&lon=%v", geoLat, geoLong) // lazy sprintf verb but oh well
	weatherReq := strings.Join([]string{YR_API, longLatStr}, "")

	// Construct GET request to weather_api
	yrApiReq, yrReqErr := http.NewRequest("GET", weatherReq, nil)
	if yrReqErr != nil {
		panic(yrReqErr)
	}

	// Set User Agent header (YR requirement)
	yrApiReq.Header.Set("Accept", "application/json")
	yrApiReq.Header.Set("User-Agent", "GoLangTest/00.1")

	// Send GET request to YR weather API
	yrApiResp, yrApiErr := xClient.Do(yrApiReq)
	if yrApiErr != nil {
		panic(yrApiErr)
	}

	// Convert to text uint8
	yrTxt, readYrErr := ioutil.ReadAll(yrApiResp.Body)
	if readYrErr != nil {
		fmt.Println("Read error YR")
		panic(readYrErr)
	}

	// Parse into YR weather data struct
	yrJson := YrJson{}
	yrJsonErrErr := json.Unmarshal(yrTxt, &yrJson)
	if yrJsonErrErr != nil {
		panic(yrJsonErrErr)
	}

	// Pull relevant variables
	yr := yrJson.Data.Timeseries[0].Data.Instant.Details
	yrTime := yrJson.Data.Timeseries[0].Time

	// Get temperatures
	yrTempC := yr.AirTemperature
	yrTempF := (yrTempC * (9 / 5)) + 32

	// Print to screen (this is just lazy)
	var output [8]string
	output[0] = fmt.Sprintf("Current weather in:           %s", geoName)
	output[1] = fmt.Sprintf("Coordinates (long, lat):      %v, %v", geoLong, geoLat)
	output[2] = fmt.Sprintf("Observation timestamp:        %v", yrTime)
	output[3] = fmt.Sprintf("Temperature (celsius):        %v\u00B0C", yrTempC)
	output[4] = fmt.Sprintf("Temperature (fahrenheit):     %v\u00B0F", yrTempF)
	output[5] = fmt.Sprintf("Airpressure at sea level:     %v hPa", yr.AirPressureAtSeaLevel)
	output[6] = fmt.Sprintf("Current humidity:             %v %%", yr.RelativeHumidity)
	output[7] = fmt.Sprintf("Wind speed:                   %v m/sec", yr.WindSpeed)

	for _, outLine := range output {
		fmt.Println(outLine)
	}

}
