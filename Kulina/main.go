package main
import (
  "fmt"  
)

const GOOGLE_API_KEY = "AIzaSyAhtIh45bWxhUhAO6vd2w_xv9YKQXv9tAw"
const CHOSEN_DATE ="2017-12-05"
const MAX_PATH_LENGTH = 90000.000000 //in meter
const OPTIMIZATION_REPETITION = 10

var KitchenList []Kitchen  
var OrderList []Order
var PathList []Path
var PlaceList []Place

func main(){  
  //get handle or connection for kulina database
  db := GetDatabaseHandle()
  defer db.Close()
  
  // get kitchens from database and put into kitchen_list  
  KitchenList = GetKitchensFromDatabase(db,CHOSEN_DATE)

  //get orders from database and put into order_list
  OrderList = GetOrdersFromDatabase(db,CHOSEN_DATE)
  
  //get places from databases
  //PlaceList = GetPlacesFromDatabase(db,CHOSEN_DATE)
  
  /*isCapacityFit := CheckKitchenCapacityToOrderQty()
  if !isCapacityFit {
    fmt.Println("We have capacity error!")
    return
  }*/
    
  //building haversine distance list of kitchens to Orders
  for i:=range(KitchenList){
    KitchenList[i].BuildOrderDistanceList()
  }  

  //building haversine distance list of orders to kitchens
  for i:=range(OrderList){
    OrderList[i].BuildKitchenDistanceList()
  }  
  
  //clustering!
  GreedyClustering()  
  for i:=range(KitchenList){
    fmt.Println(KitchenList[i].Id,KitchenList[i].Capacity.Min,KitchenList[i].Capacity.Max,KitchenList[i].DistinctOrderCount,KitchenList[i].OrderQty)
  }  
  
  //do the maximum_matching  
  GreedyClusterOptimization(OPTIMIZATION_REPETITION) 
  fmt.Println("-----------------------------")
  for i:=range(KitchenList){
    fmt.Println(KitchenList[i].Id,KitchenList[i].Capacity.Min,KitchenList[i].Capacity.Max,KitchenList[i].DistinctOrderCount,KitchenList[i].OrderQty)
  }   
         
  //make path greedily         
  //GreedyBuildPathList()
    
  //view via localhost the distribution of orders to kitchens in google map
  RunKitchensAndOrdersView()   
}
