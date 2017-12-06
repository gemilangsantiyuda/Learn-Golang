package main

import (
  "sort"
)

func NaiveClustering(orderList []Order, kitchenList []Kitchen) ([]Order,[]Kitchen){
  for i:= range orderList{
    kitchenIdx:= orderList[i].Kitchen_distance[0].Idx
    kitchenList[kitchenIdx].Order_count++
    kitchenList[kitchenIdx].Order_qty+=int(orderList[i].Qty)
  }
  orderList,kitchenList = NormalMoreMax(orderList,kitchenList)
  orderList,kitchenList = NormalLessMin(orderList,kitchenList)
	return orderList,kitchenList
} 

func NormalMoreMax(orderList []Order, kitchenList []Kitchen) ([]Order,[]Kitchen){
  for i:= range kitchenList{
    for j:= len(kitchenList[i].Order_distance)-1; j>=0 &&  kitchenList[i].Order_qty > kitchenList[i].Max_capacity;j--{
      orderIdx:=kitchenList[i].Order_distance[j].Idx
      if orderList[orderIdx].Kitchen_distance[0].Idx!=i{
        continue
      }
      chosenKitchen := -1
      for k:= 1; k<len(kitchenList); k++{
        kitchenSubIdx:= orderList[orderIdx].Kitchen_distance[k].Idx
        if kitchenList[kitchenSubIdx].Order_qty + orderList[orderIdx].Qty <= kitchenList[kitchenSubIdx].Max_capacity {
          chosenKitchen = k
          kitchenList[kitchenSubIdx].Order_qty += orderList[orderIdx].Qty
          kitchenList[kitchenSubIdx].Order_count++
          kitchenList[i].Order_count--
          kitchenList[i].Order_qty -= orderList[orderIdx].Qty
          break
        }
      }
      if chosenKitchen==-1{
        continue      
      }
      for k:=chosenKitchen;k>0;k--{
        tmp:=orderList[orderIdx].Kitchen_distance[k]
        orderList[orderIdx].Kitchen_distance[k] = orderList[orderIdx].Kitchen_distance[k-1] 
        orderList[orderIdx].Kitchen_distance[k-1] = tmp
      }
    }
  } 
  return orderList,kitchenList
}

func NormalLessMin(orderList []Order, kitchenList []Kitchen) ([]Order,[]Kitchen){
  for i:= range kitchenList{
    if kitchenList[i].Order_qty>=kitchenList[i].Min_capacity{
      continue
    }
    if kitchenList[i].Order_qty>=int(float64(kitchenList[i].Min_capacity)*0.5){
      for j:= 0; j<len(kitchenList[i].Order_distance) &&  kitchenList[i].Order_qty < kitchenList[i].Min_capacity;j++{
        orderIdx:=kitchenList[i].Order_distance[j].Idx
        currentKitchenIdx := orderList[orderIdx].Kitchen_distance[0].Idx
        if currentKitchenIdx==i || kitchenList[currentKitchenIdx].Order_qty-orderList[orderIdx].Qty<kitchenList[currentKitchenIdx].Min_capacity {
          continue
        }
        thisKitchen := -1
        for k:= 1; k<len(kitchenList); k++{
          if orderList[orderIdx].Kitchen_distance[k].Idx==i{
            thisKitchen=k
            kitchenList[currentKitchenIdx].Order_qty -= orderList[orderIdx].Qty
            kitchenList[currentKitchenIdx].Order_count--
            kitchenList[i].Order_count++
            kitchenList[i].Order_qty += orderList[orderIdx].Qty
            break 
          }
        }
        for k:=thisKitchen;k>0;k--{
          tmp:=orderList[orderIdx].Kitchen_distance[k]
          orderList[orderIdx].Kitchen_distance[k] = orderList[orderIdx].Kitchen_distance[k-1] 
          orderList[orderIdx].Kitchen_distance[k-1] = tmp
        }                
      }    
    }
  } 
  return orderList,kitchenList
}

func NaiveOptimizeClusterGreedy(a int,b int ,orderList []Order,kitchenList []Kitchen) ([]Order, []Kitchen){
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
	
	for i:= range orderList{
		if orderList[i].Kitchen_distance[0].Idx == a {
			orderA = append(orderA, OrderTmp{i,orderList[i].Latitude,orderList[i].Longitude,orderList[i].Qty})
		}
		if orderList[i].Kitchen_distance[0].Idx == b {
			orderB = append(orderB, OrderTmp{i,orderList[i].Latitude,orderList[i].Longitude,orderList[i].Qty})
		}
	}
	//fmt.Println(orderA)
	//fmt.Println(orderB)
	
  //"profit" = smaller sum of distances
	for i :=range orderA{
		for j:=range orderB {	
		  //determine wether to swap the orders make "profit"
			distOld := orderList[orderA[i].Idx].Kitchen_distance[0].Distance + orderList[orderB[j].Idx].Kitchen_distance[0].Distance
			distNew := Haversine(orderA[i].Lat,orderA[i].Lon,kitchenList[b].Loc_lat,kitchenList[b].Loc_lon)
			distNew += Haversine(orderB[j].Lat,orderB[j].Lon,kitchenList[a].Loc_lat,kitchenList[a].Loc_lon)
			dist := distOld-distNew
			req := -orderA[i].Req + orderB[j].Req
			if dist>0 {
				swapInfo = append(swapInfo , SwapInfo{orderA[i].Idx,orderB[j].Idx,dist,req,-req,false})
			}			
			
			
			//determine wether kitchen B to give the order to kitchen A make "profit"
			distOld = orderList[orderB[j].Idx].Kitchen_distance[0].Distance
			distNew = Haversine(orderB[j].Lat,orderB[j].Lon,kitchenList[a].Loc_lat,kitchenList[a].Loc_lon)
			dist = distOld-distNew
			req = orderB[j].Req
			if dist>0 {
				swapInfo = append(swapInfo , SwapInfo{-1,orderB[j].Idx,dist,req,-req,false})
			}						
		}
		//determine wether kitchen A to give the order to kitchen B make "profit"
		distOld := orderList[orderA[i].Idx].Kitchen_distance[0].Distance
		distNew := Haversine(orderA[i].Lat,orderA[i].Lon,kitchenList[b].Loc_lat,kitchenList[b].Loc_lon)
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
			if (kitchenList[kitchenA].Order_qty+reqToA >= kitchenList[kitchenA].Min_capacity) && (kitchenList[kitchenA].Order_qty+reqToA <= kitchenList[kitchenA].Max_capacity) && (kitchenList[kitchenB].Order_qty+reqToB >= kitchenList[kitchenB].Min_capacity) && (kitchenList[kitchenB].Order_qty+reqToB <= kitchenList[kitchenB].Max_capacity) {
				
				swapInfo[i].used = true
				kitchenList[kitchenA].Order_qty+=reqToA
				kitchenList[kitchenB].Order_qty+=reqToB
				
				if custA>-1{
				  chosenKitchen := -1
				  for k:= range kitchenList{
					  if orderList[custA].Kitchen_distance[k].Idx == kitchenB {
						  chosenKitchen = k
						  break
					  }				
				  }
				  for k:= chosenKitchen; k>0;k--{
					  tmp:= orderList[custA].Kitchen_distance[k]
					  orderList[custA].Kitchen_distance[k] = orderList[custA].Kitchen_distance[k-1]
					  orderList[custA].Kitchen_distance[k-1] = tmp
				  }
				}
				
				if custB>-1{
				  chosenKitchen := -1
				  for k:= range kitchenList{
					  if orderList[custB].Kitchen_distance[k].Idx == kitchenA {
						  chosenKitchen = k
						  break
					  }				
				  }
				  for k:= chosenKitchen; k>0;k--{
					  tmp:= orderList[custB].Kitchen_distance[k]
					  orderList[custB].Kitchen_distance[k] = orderList[custB].Kitchen_distance[k-1]
					  orderList[custB].Kitchen_distance[k-1] = tmp
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
			if (kitchenList[kitchenA].Order_qty+reqToA >= kitchenList[kitchenA].Min_capacity) && (kitchenList[kitchenA].Order_qty+reqToA <= kitchenList[kitchenA].Max_capacity) && (kitchenList[kitchenB].Order_qty+reqToB >= kitchenList[kitchenB].Min_capacity) && (kitchenList[kitchenB].Order_qty+reqToB <= kitchenList[kitchenB].Max_capacity) {
				
				swapInfo[i].used = true
				kitchenList[kitchenA].Order_qty+=reqToA
				kitchenList[kitchenB].Order_qty+=reqToB
				
				chosenKitchen := -1
				for k:= range kitchenList{
					if orderList[custA].Kitchen_distance[k].Idx == kitchenB {
						chosenKitchen = k
						break
					}				
				}
				for k:= chosenKitchen; k>0;k--{
					tmp:= orderList[custA].Kitchen_distance[k]
					orderList[custA].Kitchen_distance[k] = orderList[custA].Kitchen_distance[k-1]
					orderList[custA].Kitchen_distance[k-1] = tmp
				}
				
				chosenKitchen = -1
				for k:= range kitchenList{
					if orderList[custB].Kitchen_distance[k].Idx == kitchenA {
						chosenKitchen = k
						break
					}				
				}
				for k:= chosenKitchen; k>0;k--{
					tmp:= orderList[custB].Kitchen_distance[k]
					orderList[custB].Kitchen_distance[k] = orderList[custB].Kitchen_distance[k-1]
					orderList[custB].Kitchen_distance[k-1] = tmp
				}
				
				for j:= i+1;j<len(swapInfo);j++{
					if swapInfo[j].IdxA == custA || swapInfo[j].IdxB == custB {
						swapInfo[j].used=true
					}
				}
			}
		}
	}	
	return orderList,kitchenList
}

func NaiveOptimizeCluster(orderList []Order, kitchenList []Kitchen,cycle int) ([]Order, []Kitchen){
	for k:=1;k<=cycle;k++{
		for i:=range kitchenList{
			for j:=range kitchenList{
				if i<j{
					orderList,kitchenList = NaiveOptimizeClusterGreedy(i,j,orderList,kitchenList)
				}
			}
		}
	}
	for i:=range(kitchenList){
	  kitchenList[i].Order_count = 0
	  for j:=range(orderList){
	    if orderList[j].Kitchen_distance[0].Idx == i{
	      kitchenList[i].Order_count++
	    }
	  }
	}
	return orderList,kitchenList
}
