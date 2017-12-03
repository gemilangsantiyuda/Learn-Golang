package main
import (
  //"fmt"  
)

const GoogleAPIKey = "AIzaSyAhtIh45bWxhUhAO6vd2w_xv9YKQXv9tAw"
const ChosenDate ="2017-11-20"
var KitchenList []Kitchen  
var OrderList []Order

func main(){  
  //get handle or connection for kulina database
  db := GetDatabaseHandle()
  defer db.Close()
  
  // get kitchens from database and put into kitchen_list  
  KitchenList = GetKitchensFromDatabase(db,ChosenDate)

  //get orders from database and put into order_list
  OrderList = GetOrdersFromDatabase(db,ChosenDate)

  //building haversine distance list of kitchens to orders, and orders to kitchen, and orders to orders
  for i:=range(KitchenList){
    KitchenList[i].BuildOrderDistanceList(OrderList)
  }
  
  for i:=range(OrderList){
    OrderList[i].BuildOrderDistanceList(OrderList)
    OrderList[i].BuildKitchenDistanceList(KitchenList)    
  }
  
  /*//clustering!
  orderList,kitchenList = NaiveClustering(orderList,kitchenList)
  for i:=range(kitchenList){
    fmt.Println(kitchenList[i].Id,kitchenList[i].Min_capacity,kitchenList[i].Max_capacity,kitchenList[i].Order_qty)
  }  
  
  //do the maximum_matching  
  orderList,kitchenList = NaiveOptimizeCluster(orderList,kitchenList,10)
  fmt.Println("-----------------------------")
  for i:=range(kitchenList){
    fmt.Println(kitchenList[i].Id,kitchenList[i].Min_capacity,kitchenList[i].Max_capacity,kitchenList[i].Order_qty)
  }  
     
  fmt.Println(GetGoogleDistance(kitchenList[0].Loc_lat,kitchenList[0].Loc_lon,orderList[0].Latitude,orderList[0].Longitude))
  
  //build paths greedily
  pathList:= BuildPathList(orderList,kitchenList)
  fmt.Println(pathList)
  
  */
  RunKitchensAndOrdersView()
   
}
