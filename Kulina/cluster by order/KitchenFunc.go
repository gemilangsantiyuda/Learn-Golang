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

func (kitchen *Kitchen) BuildOrderDistanceList(){
  for i:=range(OrderList){
    dist := GetHaversineDistance(kitchen.Coord,OrderList[i].Coord)
    kitchen.OrderDistanceList = append(kitchen.OrderDistanceList, Distance{i,dist})
  }
  sort.SliceStable(kitchen.OrderDistanceList,func(lhs,rhs int) bool {
    return kitchen.OrderDistanceList[lhs].Distance<kitchen.OrderDistanceList[rhs].Distance
  })
}

func (kitchen *Kitchen) BuildPlaceDistanceList(){
  for i:=range(PlaceList){
    dist := GetHaversineDistance(kitchen.Coord,PlaceList[i].Coord)
    kitchen.PlaceDistanceList = append(kitchen.PlaceDistanceList, Distance{i,dist})
  }
  sort.SliceStable(kitchen.PlaceDistanceList,func(lhs,rhs int) bool {
    return kitchen.PlaceDistanceList[lhs].Distance<kitchen.PlaceDistanceList[rhs].Distance
  })
}

func (kitchen *Kitchen) CanServe(order *Order) bool{
  kitchenHasEnoughCapacity:=(kitchen.OrderQty + order.Qty <= kitchen.Capacity.Max)
  isCloseEnough:= (GetHaversineDistance(kitchen.Coord,order.Coord)<=MAX_PATH_LENGTH)
  return kitchenHasEnoughCapacity && isCloseEnough
}

func (kitchen *Kitchen) CanLetGo(order *Order) bool{
  kitchenHasEnoughOrder:=(kitchen.OrderQty-order.Qty>=kitchen.Capacity.Min)
  return kitchenHasEnoughOrder 
}

func (kitchenOrigin *Kitchen) GiveOrderToKitchen(order *Order, kitchenDestination *Kitchen){
  kitchenDestination.OrderQty += order.Qty
  kitchenDestination.DistinctOrderCount++
  kitchenOrigin.DistinctOrderCount--
  kitchenOrigin.OrderQty -= order.Qty
  
  kitchenDestinationIndex:=-1
  for i:= range order.KitchenDistanceList{
    if KitchenList[order.KitchenDistanceList[i].Index].Id == kitchenDestination.Id {
      kitchenDestinationIndex = i
      break
    }
  }
  
  for k:=kitchenDestinationIndex;k>0;k--{
    tmp:=order.KitchenDistanceList[k]
    order.KitchenDistanceList[k] = order.KitchenDistanceList[k-1]
    order.KitchenDistanceList[k-1] = tmp
  }
}
