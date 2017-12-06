package main
import (
  "math"
  "strconv"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "time"
  "log"
)

func GetHaversineDistance(coordOrigin,coordDestination Coordinate) float64 {
	DY := math.Abs(coordOrigin.Latitude - coordDestination.Latitude)/180*math.Pi
	DX := math.Abs(coordOrigin.Longitude - coordDestination.Longitude)/180*math.Pi
	Y1 := coordOrigin.Latitude/180*math.Pi
	Y2 := coordDestination.Latitude/180*math.Pi
	R := 6372800.00000000 //in Meter ! it is the average great-elliptic or great-circle radius
	a := math.Sin(DY/2)*math.Sin(DY/2) + math.Cos(Y1)*math.Cos(Y2)*math.Sin(DX/2)*math.Sin(DX/2)
	c := 2*math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	return R*c
}

func GetGoogleDistance(coordOrigin,coordDestination Coordinate) int{
  
  //convert lat and lon to string for google api query
  stringLatitudeOrigin:= strconv.FormatFloat(coordOrigin.Latitude, 'f', 7, 64)
  stringLongitudeOrigin:= strconv.FormatFloat(coordOrigin.Longitude, 'f', 7, 64)
  stringLatitudeDestination:= strconv.FormatFloat(coordDestination.Latitude, 'f', 7, 64)
  stringLongitudeDestination:= strconv.FormatFloat(coordDestination.Longitude, 'f', 7, 64)     
  
  //fmt.Println(slat1,slon1,slat2,slon2)
  url := "https://maps.googleapis.com/maps/api/distancematrix/json?origins="+stringLatitudeOrigin+","+stringLongitudeOrigin+"&destinations="+stringLatitudeDestination+","+stringLongitudeDestination+"&departure_time=1541202457&traffic_model=best_guess&key="+GOOGLE_API_KEY
  //fmt.Println(url)

  spaceClient := http.Client{
        Timeout: time.Second * 5, // Maximum of 5 secs
  }
  req, err := http.NewRequest(http.MethodGet, url, nil)
  if err != nil {
      log.Fatal(err)
  }

  req.Header.Set("User-Agent", "kulina-matching-beta")

  res, getErr := spaceClient.Do(req)
  if getErr != nil {
      log.Fatal(getErr)
  }

  body, readErr := ioutil.ReadAll(res.Body)
  if readErr != nil {
      log.Fatal(readErr)
  }

  var googleDistance GoogleAPIDistanceResponse
  jsonErr := json.Unmarshal(body, &googleDistance)
  if jsonErr != nil {
      log.Fatal(jsonErr)
  }
  return googleDistance.Rows[0].Elements[0].Distance.Value
}

