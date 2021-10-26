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

type Location struct {
	Data []struct {
		Name        string `json:"skrivemåte"`
		Coordinates struct {
			Long float64 `json:"øst"`
			Lat  float64 `json:"nord"`
		} `json:"representasjonspunkt"`
	} `json:"navn"`
}

type YR struct {
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
	var userLoc string = "oslo"
	if len(os.Args) > 1 {
		userLoc = os.Args[1]
	} else {
		fmt.Println(strings.Join([]string{"Usage:", os.Args[0], "<place>"}, " "))
		fmt.Println("Error: Location string argument missing.")
		os.Exit(1)
	}

	// Set HTTP request timeout (10 secs)
	var xClient = &http.Client{Timeout: 10 * time.Second}

	// create GET request for coordinates
	const coord_api string = "https://ws.geonorge.no/stedsnavn/v1/navn?sok="
	const coord_api_post string = "&utkoordsys=4258&treffPerSide=1&side=1"
	coord_api_request := strings.Join([]string{coord_api, userLoc, coord_api_post}, "")

	// Send GET request to coord_api
	coord_resp, err := xClient.Get(coord_api_request)
	if err != nil {
		panic(err)
	}

	// Convert to text uint8
	coord_body, readErr := ioutil.ReadAll(coord_resp.Body)
	if readErr != nil {
		fmt.Println("Read error GeoNorge")
		panic(readErr)
	}

	// Parse into Location struct
	geo := Location{}
	jsonErr := json.Unmarshal(coord_body, &geo)
	if jsonErr != nil {
		fmt.Println("JSON error GeoNorge")
		panic(jsonErr)
	}

	// Pull relevant variables
	var place_name string = geo.Data[0].Name
	var place_long float64 = geo.Data[0].Coordinates.Long
	var place_lat float64 = geo.Data[0].Coordinates.Lat

	// create GET request for weather query
	const weather_api string = "https://api.met.no/weatherapi/locationforecast/2.0/compact.json?"
	long_lat_str := fmt.Sprintf("lat=%v&lon=%v", place_lat, place_long) // lazy sprintf verb but oh well
	weather_request := strings.Join([]string{weather_api, long_lat_str}, "")

	// Construct GET request to weather_api
	yr_api_request, yrReqErr := http.NewRequest("GET", weather_request, nil)
	if yrReqErr != nil {
		panic(yrReqErr)
	}

	// Set User Agent header (YR requirement)
	yr_api_request.Header.Set("Accept", "application/json")
	yr_api_request.Header.Set("User-Agent", "GoLangTest/00.1")

	// Send GET request to YR weather API
	yr_resp, yrRespErr := xClient.Do(yr_api_request)
	if yrRespErr != nil {
		panic(yrRespErr)
	}

	// Convert to text uint8
	yr_raw, readYrErr := ioutil.ReadAll(yr_resp.Body)
	if readYrErr != nil {
		fmt.Println("Read error YR")
		panic(readYrErr)
	}

	// Parse into YR weather data struct
	yr_json := YR{}
	json2Err := json.Unmarshal(yr_raw, &yr_json)
	if json2Err != nil {
		fmt.Println("JSON error YR")
		// fmt.Println("yr_raw:")
		// fmt.Println(yr_raw)
		fmt.Println("yr_json:")
		fmt.Println(yr_json)
		panic(json2Err)
	}

	// Pull relevant variables
	yr := yr_json.Data.Timeseries[0].Data.Instant.Details
	yr_time := yr_json.Data.Timeseries[0].Time
	// var pressure_hPa int = yr.AirPressureAtSeaLevel
	// var temp_celsius float64 = yr.AirTemperature
	// var humidity_percent float64 = yr.RelativeHumidity
	// var wind_meters_sec int = yr.WindSpeed

	// Get fahrenheit values for our visitors
	// Using (0°C × 9/5) + 32 = 32°F
	temp_celsius := yr.AirTemperature
	temp_fahrenheit := (temp_celsius * (9 / 5)) + 32

	// Print to screen (this is just lazy)
	var output [8]string
	output[0] = fmt.Sprintf("Current weather in:           %s", place_name)
	output[1] = fmt.Sprintf("Coordinates (long, lat):      %v, %v", place_long, place_lat)
	output[2] = fmt.Sprintf("Observation timestamp:        %v", yr_time)
	output[3] = fmt.Sprintf("Temperature (celsius):        %v\u00B0C", temp_celsius)
	output[4] = fmt.Sprintf("Temperature (fahrenheit):     %v\u00B0F", temp_fahrenheit)
	output[5] = fmt.Sprintf("Airpressure at sea level:     %v hPa", yr.AirPressureAtSeaLevel)
	output[6] = fmt.Sprintf("Current humidity:             %v %%", yr.RelativeHumidity)
	output[7] = fmt.Sprintf("Wind speed:                   %v m/sec", yr.WindSpeed)

	for _, showline := range output {
		fmt.Println(showline)
	}

}
