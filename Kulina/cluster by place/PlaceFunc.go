package main

import (
  "database/sql"
_ "github.com/go-sql-driver/mysql"
  "sort"
)

func (place *Place) ScanFromSQL(row *sql.Rows){
  var placeSQL PlaceSQLContainer
  err:=row.Scan(&placeSQL.Id,&placeSQL.Qty) 
  checkErr(err)
  place.Id  = placeSQL.Id.String
  place.Qty = int(placeSQL.Qty.Int64)
}

func (place *Place) BuildKitchenDistanceList(){
  for i:=range(KitchenList){
    dist := GetHaversineDistance(place.Coord,KitchenList[i].Coord)
    place.KitchenDistanceList = append(place.KitchenDistanceList, Distance{i,dist})
  }
  sort.SliceStable(place.KitchenDistanceList,func(lhs,rhs int) bool {
    return place.KitchenDistanceList[lhs].Distance<place.KitchenDistanceList[rhs].Distance
  })
}

func (place *Place) BuildPlaceDistanceList(){
  for i:=range(PlaceList){
    placeDestination:=&PlaceList[i]
    if place!=placeDestination{
      dist := GetHaversineDistance(place.Coord,placeDestination.Coord)
      place.PlaceDistanceList = append(place.PlaceDistanceList, Distance{i,dist})
    }    
  }
  sort.SliceStable(place.PlaceDistanceList,func(lhs,rhs int) bool {
    return place.PlaceDistanceList[lhs].Distance<place.PlaceDistanceList[rhs].Distance
  })
}
