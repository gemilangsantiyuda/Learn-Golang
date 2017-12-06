package main

import (
  "database/sql"
_ "github.com/go-sql-driver/mysql"
  "sort"  
)

func (order *Order) ScanFromSQL(row *sql.Rows){
  var orderSQL OrderSQLContainer
  err:=row.Scan(&orderSQL.Id,&orderSQL.Qty,&orderSQL.Coord.Latitude,&orderSQL.Coord.Longitude,&orderSQL.PlaceId)  
  checkErr(err)
  order.Id=int(orderSQL.Id.Int64)
  order.Qty=int(orderSQL.Qty.Int64)
  order.Coord= Coordinate{orderSQL.Coord.Latitude.Float64, orderSQL.Coord.Longitude.Float64}
  order.PlaceId = orderSQL.PlaceId.String
}

func (order *Order) BuildOrderDistanceList(orderList []Order){
  for i:=range(orderList){
    if orderList[i].Id != order.Id {
      dist := GetHaversineDistance(order.Coord,orderList[i].Coord)
      order.OrderDistanceList = append(order.OrderDistanceList, Distance{i,dist})  
    }    
  }
  sort.SliceStable(order.OrderDistanceList,func(lhs,rhs int) bool {
    return order.OrderDistanceList[lhs].Distance<order.OrderDistanceList[rhs].Distance
  })
}

func (order *Order) BuildKitchenDistanceList(){
  for i:=range(KitchenList){
    dist := GetHaversineDistance(order.Coord,KitchenList[i].Coord)
    order.KitchenDistanceList = append(order.KitchenDistanceList, Distance{i,dist})  
  }
  sort.SliceStable(order.KitchenDistanceList,func(lhs,rhs int) bool {
    return order.KitchenDistanceList[lhs].Distance<order.KitchenDistanceList[rhs].Distance
  })
}

func (order *Order) GetServerKitchenIndex() int{
  return order.KitchenDistanceList[0].Index
}
