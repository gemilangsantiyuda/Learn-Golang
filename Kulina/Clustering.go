package main

import (
  //"sort"
)

func GreedyClustering(){
  for i:= range OrderList{
    kitchenIdx:=OrderList[i].KitchenDistanceList[0].Index
    KitchenList[kitchenIdx].DistinctOrderCount++
    KitchenList[kitchenIdx].OrderQty+=int(OrderList[i].Qty)
  }
  ResolveMaximumCapacityViolation()  
  ResolveMinimumCapacityViolation()
} 


func ResolveMaximumCapacityViolation(){
  //For kitchens violating the maximum capacity :
  //Starting from the furthes order to the nearest
  //Give the order which is served by the kitchen
  //To the alternative kitchen nearest to that order and is able to serve that order
  //Do this until the kitchen does not violate the maximum capacity
  for i:= range KitchenList {
    kitchenOrigin:= &KitchenList[i]
    for j:= len(kitchenOrigin.OrderDistanceList)-1; j>=0 &&  kitchenOrigin.OrderQty > kitchenOrigin.Capacity.Max;j--{
      orderToLetGo:=&OrderList[kitchenOrigin.OrderDistanceList[j].Index]
      if orderToLetGo.GetServerKitchenIndex()!=i{
        continue
      }      
      for k:= 1; k<len(KitchenList); k++{
        kitchenDestination:= &KitchenList[orderToLetGo.KitchenDistanceList[k].Index]
        if kitchenDestination.CanServe(orderToLetGo) {
          kitchenOrigin.GiveOrderToKitchen(orderToLetGo,kitchenDestination)          
          break
        }
      }      
    }
  } 
}

func ResolveMinimumCapacityViolation(){  
  for i:= range KitchenList{
    thisKitchen:=&KitchenList[i]    
    if thisKitchen.OrderQty>=thisKitchen.Capacity.Min{
      continue
    }    
    //For kitchens violating the minimum capacity and has more than 50% of the minimum capacity:
    //Starting from the nearest order to the furthest
    //Get the order from the original kitchen serving the order
    //Only when that kitchen is able to let go the forementioned order without violating its minimum capacity
    //Do this until the kitchen does not violate the minimum capacity
    if thisKitchen.OrderQty>=int(float64(thisKitchen.Capacity.Min)*0.5){
      kitchenDestination:=thisKitchen
      for j:= 0; j<len(kitchenDestination.OrderDistanceList) &&  kitchenDestination.OrderQty < kitchenDestination.Capacity.Min;j++{
        orderToTakeOver :=&OrderList[kitchenDestination.OrderDistanceList[j].Index]
        kitchenOrigin   :=&KitchenList[orderToTakeOver.KitchenDistanceList[0].Index]
        if (kitchenOrigin.Id!=kitchenDestination.Id) && kitchenDestination.CanLetGo(orderToTakeOver){
          kitchenOrigin.GiveOrderToKitchen(orderToTakeOver,kitchenDestination)
        }
      }    
    } else {                 
    //if it has less than 50% capacity, distribute the orders it serves to other kitchens
    //nearest to that order, and can serve the order without violating maximum capacity
      kitchenOrigin:=thisKitchen
      for j:=0;j<len(kitchenOrigin.OrderDistanceList) && kitchenOrigin.OrderQty>0;j++{
        orderToLetGo := &OrderList[kitchenOrigin.OrderDistanceList[j].Index]
        if orderToLetGo.GetServerKitchenIndex()!=i{
          continue
        }
        for k:= 1; k<len(KitchenList); k++{
          kitchenDestination:= &KitchenList[orderToLetGo.KitchenDistanceList[k].Index]
          if kitchenDestination.CanServe(orderToLetGo) {
            kitchenOrigin.GiveOrderToKitchen(orderToLetGo,kitchenDestination)          
            break
          }
        }      
      }
    }
  } 
}

/*func NaiveOptimizeClusterGreedy(a int,b int ,OrderList []Order,KitchenList []Kitchen) ([]Order, []Kitchen){
	type OrderTmp struct{
		Idx int
		Lat,Lon float64
		Req int 
	}
	type SwapInfo struct{
		IdxA,IdxB int
		Dist float64
		ReqToA,ReqToB int
		used bool
	}
	var orderA,orderB []OrderTmp
	var swapInfo []SwapInfo
	//var swapChosen []int
	
	for i:= range OrderList{
		if OrderList[i].Kitchen_distance[0].Idx == a {
			orderA = append(orderA, OrderTmp{i,OrderList[i].Latitude,OrderList[i].Longitude,OrderList[i].Qty})
		}
		if OrderList[i].Kitchen_distance[0].Idx == b {
			orderB = append(orderB, OrderTmp{i,OrderList[i].Latitude,OrderList[i].Longitude,OrderList[i].Qty})
		}
	}
	//fmt.Println(orderA)
	//fmt.Println(orderB)
	
  //"profit" = smaller sum of distances
	for i :=range orderA{
		for j:=range orderB {	
		  //determine wether to swap the orders make "profit"
			distOld := OrderList[orderA[i].Idx].Kitchen_distance[0].Distance + OrderList[orderB[j].Idx].Kitchen_distance[0].Distance
			distNew := Haversine(orderA[i].Lat,orderA[i].Lon,KitchenList[b].Loc_lat,KitchenList[b].Loc_lon)
			distNew += Haversine(orderB[j].Lat,orderB[j].Lon,KitchenList[a].Loc_lat,KitchenList[a].Loc_lon)
			dist := distOld-distNew
			req := -orderA[i].Req + orderB[j].Req
			if dist>0 {
				swapInfo = append(swapInfo , SwapInfo{orderA[i].Idx,orderB[j].Idx,dist,req,-req,false})
			}			
			
			
			//determine wether kitchen B to give the order to kitchen A make "profit"
			distOld = OrderList[orderB[j].Idx].Kitchen_distance[0].Distance
			distNew = Haversine(orderB[j].Lat,orderB[j].Lon,KitchenList[a].Loc_lat,KitchenList[a].Loc_lon)
			dist = distOld-distNew
			req = orderB[j].Req
			if dist>0 {
				swapInfo = append(swapInfo , SwapInfo{-1,orderB[j].Idx,dist,req,-req,false})
			}						
		}
		//determine wether kitchen A to give the order to kitchen B make "profit"
		distOld := OrderList[orderA[i].Idx].Kitchen_distance[0].Distance
		distNew := Haversine(orderA[i].Lat,orderA[i].Lon,KitchenList[b].Loc_lat,KitchenList[b].Loc_lon)
		dist := distOld-distNew
		req := -orderA[i].Req
		if dist>0 {
			swapInfo = append(swapInfo , SwapInfo{orderA[i].Idx,-1,dist,req,-req,false})
		}						
	}
	sort.SliceStable(swapInfo,func (lhs,rhs int) bool {
	  return swapInfo[lhs].Dist>swapInfo[rhs].Dist
	})
  //fmt.Println(swapInfo)
	
	for i:=range swapInfo{
	  //if kitchen is giving away order
	  if (swapInfo[i].IdxA==-1 || swapInfo[i].IdxB == -1) && !swapInfo[i].used {         
			custA:=swapInfo[i].IdxA
			custB:=swapInfo[i].IdxB
			kitchenA:=a
			kitchenB:=b		
			reqToA:=swapInfo[i].ReqToA
			reqToB:=swapInfo[i].ReqToB
			if (KitchenList[kitchenA].Order_qty+reqToA >= KitchenList[kitchenA].Min_capacity) && (KitchenList[kitchenA].Order_qty+reqToA <= KitchenList[kitchenA].Max_capacity) && (KitchenList[kitchenB].Order_qty+reqToB >= KitchenList[kitchenB].Min_capacity) && (KitchenList[kitchenB].Order_qty+reqToB <= KitchenList[kitchenB].Max_capacity) {
				
				swapInfo[i].used = true
				KitchenList[kitchenA].Order_qty+=reqToA
				KitchenList[kitchenB].Order_qty+=reqToB
				
				if custA>-1{
				  chosenKitchen := -1
				  for k:= range KitchenList{
					  if OrderList[custA].Kitchen_distance[k].Idx == kitchenB {
						  chosenKitchen = k
						  break
					  }				
				  }
				  for k:= chosenKitchen; k>0;k--{
					  tmp:= OrderList[custA].Kitchen_distance[k]
					  OrderList[custA].Kitchen_distance[k] = OrderList[custA].Kitchen_distance[k-1]
					  OrderList[custA].Kitchen_distance[k-1] = tmp
				  }
				}
				
				if custB>-1{
				  chosenKitchen := -1
				  for k:= range KitchenList{
					  if OrderList[custB].Kitchen_distance[k].Idx == kitchenA {
						  chosenKitchen = k
						  break
					  }				
				  }
				  for k:= chosenKitchen; k>0;k--{
					  tmp:= OrderList[custB].Kitchen_distance[k]
					  OrderList[custB].Kitchen_distance[k] = OrderList[custB].Kitchen_distance[k-1]
					  OrderList[custB].Kitchen_distance[k-1] = tmp
				  }
				}
				
				
				for j:= i+1;j<len(swapInfo);j++{
					if (swapInfo[j].IdxA == custA && custA>-1) || (swapInfo[j].IdxB == custB && custB >-1) {
						swapInfo[j].used=true
					}
				}
			}
	  }
	  
	  //if kitchens are swapping orders
		if !swapInfo[i].used {
			custA:=swapInfo[i].IdxA
			custB:=swapInfo[i].IdxB
			kitchenA:=a
			kitchenB:=b		
			reqToA:=swapInfo[i].ReqToA
			reqToB:=swapInfo[i].ReqToB
			if (KitchenList[kitchenA].Order_qty+reqToA >= KitchenList[kitchenA].Min_capacity) && (KitchenList[kitchenA].Order_qty+reqToA <= KitchenList[kitchenA].Max_capacity) && (KitchenList[kitchenB].Order_qty+reqToB >= KitchenList[kitchenB].Min_capacity) && (KitchenList[kitchenB].Order_qty+reqToB <= KitchenList[kitchenB].Max_capacity) {
				
				swapInfo[i].used = true
				KitchenList[kitchenA].Order_qty+=reqToA
				KitchenList[kitchenB].Order_qty+=reqToB
				
				chosenKitchen := -1
				for k:= range KitchenList{
					if OrderList[custA].Kitchen_distance[k].Idx == kitchenB {
						chosenKitchen = k
						break
					}				
				}
				for k:= chosenKitchen; k>0;k--{
					tmp:= OrderList[custA].Kitchen_distance[k]
					OrderList[custA].Kitchen_distance[k] = OrderList[custA].Kitchen_distance[k-1]
					OrderList[custA].Kitchen_distance[k-1] = tmp
				}
				
				chosenKitchen = -1
				for k:= range KitchenList{
					if OrderList[custB].Kitchen_distance[k].Idx == kitchenA {
						chosenKitchen = k
						break
					}				
				}
				for k:= chosenKitchen; k>0;k--{
					tmp:= OrderList[custB].Kitchen_distance[k]
					OrderList[custB].Kitchen_distance[k] = OrderList[custB].Kitchen_distance[k-1]
					OrderList[custB].Kitchen_distance[k-1] = tmp
				}
				
				for j:= i+1;j<len(swapInfo);j++{
					if swapInfo[j].IdxA == custA || swapInfo[j].IdxB == custB {
						swapInfo[j].used=true
					}
				}
			}
		}
	}	
	return OrderList,KitchenList
}

func NaiveOptimizeCluster(OrderList []Order, KitchenList []Kitchen,cycle int) ([]Order, []Kitchen){
	for k:=1;k<=cycle;k++{
		for i:=range KitchenList{
			for j:=range KitchenList{
				if i<j{
					OrderList,KitchenList = NaiveOptimizeClusterGreedy(i,j,OrderList,KitchenList)
				}
			}
		}
	}
	for i:=range(KitchenList){
	  KitchenList[i].Order_count = 0
	  for j:=range(OrderList){
	    if OrderList[j].Kitchen_distance[0].Idx == i{
	      KitchenList[i].Order_count++
	    }
	  }
	}
	return OrderList,KitchenList
}*/
