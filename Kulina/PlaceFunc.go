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

/*func (place *Place) GetServerKitchenIndex(){
  return place.KitchenDistanceList[0].Index
}*/