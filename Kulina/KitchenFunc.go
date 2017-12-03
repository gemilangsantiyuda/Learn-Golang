package main

import (
  "database/sql"
_ "github.com/go-sql-driver/mysql"
  "sort"  
)

func (kitchen *Kitchen) ScanFromSQL(row *sql.Rows){
  var kitchenSQL KitchenSQLContainer
  err:=row.Scan(&kitchenSQL.Id,&kitchenSQL.Name,&kitchenSQL.Capacity.Min,&kitchenSQL.Capacity.Max,&kitchenSQL.Capacity.Optimum,&kitchenSQL.Coord.Latitude,&kitchenSQL.Coord.Longitude) 
  checkErr(err)
  kitchen.Id=int(kitchenSQL.Id.Int64)
  kitchen.Name=kitchenSQL.Name.String
  kitchen.Capacity.Min=int(kitchenSQL.Capacity.Min.Int64)
  kitchen.Capacity.Max=int(kitchenSQL.Capacity.Max.Int64)
  kitchen.Capacity.Optimum=int(kitchenSQL.Capacity.Optimum.Int64)
  kitchen.Coord= Coordinate{kitchenSQL.Coord.Latitude.Float64, kitchenSQL.Coord.Longitude.Float64}
}

func (kitchen *Kitchen) BuildOrderDistanceList(orderList []Order){
  for i:=range(orderList){
    dist := GetHaversineDistance(kitchen.Coord,orderList[i].Coord)
    kitchen.OrderDistanceList = append(kitchen.OrderDistanceList, Distance{i,dist})
  }
  sort.SliceStable(kitchen.OrderDistanceList,func(lhs,rhs int) bool {
    return kitchen.OrderDistanceList[lhs].Distance<kitchen.OrderDistanceList[rhs].Distance
  })
}
