package main
/*import (
  _ "github.com/go-sql-driver/mysql"
  "sort"
  "strconv"
  "fmt"
)

func BuildPath(orderList []Order,kitchenList []Kitchen, chosenKitchen int) Path{
  var newPath Path
  newPath.Path_idx = append(newPath.Path_idx,chosenKitchen)

  //in order to minimize the queries of google API we have a strategy
  //to only consider 5 NEAREST UNSERVED ORDERS BASED ON HAVERSINE DISTANCE
  //these 5 orders then will be google queried for its real distance in the map

  
  //find the first order from kitchen 
  nextOrder := GetNearestOrderDistance(kitchenList[chosenKitchen].Loc_lat,kitchenList[chosenKitchen].Loc_lon,newPath,chosenKitchen,kitchenList[chosenKitchen],orderList)
  
  //insert order into new path
  newPath.Path_idx = append(newPath.Path_idx,nextOrder.Idx)
  newPath.Length+=nextOrder.Distance
  newPath.Order_qty+=orderList[nextOrder.Idx].Qty
  orderList[nextOrder.Idx].Used = true  
  kitchenList[chosenKitchen].Order_count--
  

  //complete the path with the closest order and still from the same kitchen
  currentOrder:= nextOrder  
  for true { //break until no order which fits the constraints was found
    nextOrder := GetNearestOrderDistance(orderList[currentOrder.Idx].Latitude,orderList[currentOrder.Idx].Longitude,newPath,chosenKitchen,kitchenList[chosenKitchen],orderList)
           
    //check if the order is not found (?) or the length of the path does not exceed 90000 meters
    if nextOrder.Idx == -1 || (newPath.Length + nextOrder.Distance > 90000){
      break
    }
    //insert order into new path
    newPath.Path_idx = append(newPath.Path_idx,nextOrder.Idx)
    newPath.Length+=nextOrder.Distance
    newPath.Order_qty+=orderList[nextOrder.Idx].Qty
    orderList[nextOrder.Idx].Used = true  
    kitchenList[chosenKitchen].Order_count--
    
    currentOrder = nextOrder
  }
  
  
  return newPath
}

func GreedyBuildPathList () {  
  orderLeft := len(OrderList) 
  for orderLeft>0 {    
    //always make new path from the kitchen serving the most distinct orders
    originKitchen:=nil
    maxOrder:=0
    for i:=range(kitchenList){
      thisKitchen:=&KitchenList[i]
      if thisKitchen.DistinctOrderCount >maxOrder{
        maxOrder = thisKitchen.DistinctOrderCount
        originKitchen = thisKitchen
      }
    }    
    newPath := BuildPath(originKitchen)
    newPath.Id = len(pathList)
    newPath.GiveName()
    pathList = append(pathList,newPath)
    
    orderLeft = 0
    for i:=range(kitchenList){
      orderLeft+= kitchenList[i].Order_count
    }
    fmt.Println(orderLeft)
  }
  return pathList
}

func (p *Path) GiveName(){
  p.Path_name = make([]string,len(p.Path_idx))
  p.Path_name[0] = kitchenList[p.Path_idx[0]].Name
  for i:=1 ; i<len(p.Path_name); i++{
    p.Path_name[i] = strconv.Itoa(orderList[p.Path_idx[i]].Id)
  }
} 

func GetNearestOrderDistance(lat,lon float64 ,path Path, kitchenIdx int, kitchen Kitchen, orderList []Order) Distance {
  var tmpList []Distance
  fmt.Println(kitchen.Name,kitchen.Order_qty,kitchen.Max_capacity)
  //get 5 haversine closest order, and from the same kitchen and fits the path 
  for i:=0;i<len(kitchen.Order_distance) && len(tmpList)<5;i++ {      
    idx := kitchen.Order_distance[i].Idx
    if !orderList[idx].Used && (orderList[idx].Kitchen_distance[0].Idx == kitchenIdx) && (path.Order_qty + orderList[idx].Qty <= 40) {
      tmpList = append(tmpList,Distance{idx,0})
    }      
  }
  
  //calculate their real distance
  for i:=0;i<len(tmpList);i++{
    tmpList[i].Distance = float64(GetGoogleDistance(lat,lon,orderList[tmpList[i].Idx].Latitude,orderList[tmpList[i].Idx].Longitude))
  }
  sort.SliceStable(tmpList,func (lhs,rhs int) bool {
	  return tmpList[lhs].Distance>tmpList[rhs].Distance
	})
	
	//return the closest one
	if len(tmpList)==0 {
	  return Distance{-1,-1}
	} else {
	  return tmpList[0]
	}  
} 

/*func StealNearestOrderDistance(lat,lon float64 ,path Path, orderIdx int, kitchen Kitchen, orderList []Order) Distance {
  var tmpList []Distance
  fmt.Println(kitchen.Name,kitchen.Order_qty,kitchen.Max_capacity)
  //get 5 haversine closest order, and from any kitchen as long as it fits the original kitchen and the path
  for i:=0;i<len(orderList[orderIdx].) && len(tmpList)<5;i++ {      
    idx := kitchen.Order_distance[i].Idx
    if !orderList[idx].Used && (orderList[idx].Kitchen_distance[0].Idx == kitchenIdx) && (path.Order_qty + orderList[idx].Qty <= 40) {
      tmpList = append(tmpList,Distance{idx,0})
    }      
  }
  
  //calculate their real distance
  for i:=0;i<len(tmpList);i++{
    tmpList[i].Distance = float64(GetGoogleDistance(lat,lon,orderList[tmpList[i].Idx].Latitude,orderList[tmpList[i].Idx].Longitude))
  }
  sort.SliceStable(tmpList,func (lhs,rhs int) bool {
	  return tmpList[lhs].Distance>tmpList[rhs].Distance
	})
	
	//return the closest one
	if len(tmpList)==0 {
	  return Distance{-1,-1}
	} else {
	  return tmpList[0]
	}  
}*/
