package main
import (
  "math"
  "net/http"
  "encoding/json"
  "time"
  "io/ioutil"
  "log"
  "strconv"
)

func Haversine(lat1,lon1,lat2,lon2 float64) float64 {
	DY := math.Abs(lat1-lat2)/180*math.Pi
	DX := math.Abs(lon1-lon2)/180*math.Pi
	Y1 := lat1/180*math.Pi
	Y2 := lat2/180*math.Pi
	R := 6371.00000
	a := math.Sin(DY/2)*math.Sin(DY/2) + math.Cos(Y1)*math.Cos(Y2)*math.Sin(DX/2)*math.Sin(DX/2)
	c := 2*math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	return R*c
}

func GetGoogleDistance(lat1,lon1,lat2,lon2 float64) int{
  
  //convert lat and lon to string for google api query
  slat1:= strconv.FormatFloat(lat1, 'f', 7, 64)
  slon1:= strconv.FormatFloat(lon1, 'f', 7, 64)
  slat2:= strconv.FormatFloat(lat2, 'f', 7, 64)
  slon2:= strconv.FormatFloat(lon2, 'f', 7, 64)     
  
  //fmt.Println(slat1,slon1,slat2,slon2)
  url := "https://maps.googleapis.com/maps/api/distancematrix/json?origins="+slat1+","+slon1+"&destinations="+slat2+","+slon2+"&departure_time=1541202457&traffic_model=best_guess&key="+googleDistanceKey
  //fmt.Println(url)

  spaceClient := http.Client{
        Timeout: time.Second * 5, // Maximum of 2 secs
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

  var googleDistance GoogleDistanceResponse
  jsonErr := json.Unmarshal(body, &googleDistance)
  if jsonErr != nil {
      log.Fatal(jsonErr)
  }
  return googleDistance.Rows[0].Elements[0].Distance.Value
}

