package main

import (
  "sort"
  //"fmt"
)

func GreedyClustering(){
  for i:= range OrderList{
    thisOrder:=&OrderList[i]
    kitchen:=&KitchenList[thisOrder.GetServerKitchenIndex()]
    kitchen.DistinctOrderCount++
    kitchen.OrderQty+= thisOrder.Qty
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

func GreedyMaximumOrderMatching(kitchenA,kitchenB *Kitchen){
	var orderListA,orderListB []*Order	
	type SwapInfo struct{
		OrderA,OrderB *Order
		DistanceProfit float64
		IsUsed bool
	}
	var swapInfo []SwapInfo	
	for i:= range OrderList{
	  thisOrder:=OrderList[i]	  	  
	  servingKitchen:=&KitchenList[thisOrder.KitchenDistanceList[0].Index]
		if servingKitchen.Id == kitchenA.Id {
			orderListA = append(orderListA, &thisOrder)
		} else if servingKitchen.Id == kitchenB.Id {
			orderListB = append(orderListB, &thisOrder)		  
		}		
	}
	for i :=range orderListA{
	  orderA:=orderListA[i]
	  //debug fmt.Println(orderA.Id,KitchenList[orderA.GetServerKitchenIndex()].Name)	  
		for j:=range orderListB {	
		  orderB:=orderListB[j]
		  //determine wether swapping the orders make "distanceProfit" or not		  
			oldDistanceSum := orderA.KitchenDistanceList[0].Distance + orderB.KitchenDistanceList[0].Distance
			newDistanceSum := GetHaversineDistance(orderA.Coord,kitchenB.Coord) + GetHaversineDistance(orderB.Coord,kitchenA.Coord)
			distanceProfit := oldDistanceSum-newDistanceSum
			if distanceProfit>0 /*&& (GetHaversineDistance(orderA.Coord,kitchenB.Coord)<=orderA.KitchenDistanceList[0].Distance) && (orderB.KitchenDistanceList[0].Distance <= GetHaversineDistance(orderB.Coord,kitchenA.Coord))*/{
				swapInfo = append(swapInfo , SwapInfo{orderA,orderB,distanceProfit,false})
			}									
		}
	}
  //determine wether kitchen A giving its order to kitchen B make "distanceProfit"
  for i:= range orderListA{
    orderA:=orderListA[i]
    oldDistanceSum := orderA.KitchenDistanceList[0].Distance
    newDistanceSum := GetHaversineDistance(orderA.Coord,kitchenB.Coord)
    distanceProfit := oldDistanceSum-newDistanceSum    
    if distanceProfit>0 {
	    swapInfo = append(swapInfo , SwapInfo{orderA,nil,distanceProfit,false})
    }			
  }	
  //determine wether kitchen B giving its order to kitchen A make "distanceProfit"
  for i:= range orderListB{
    orderB:= orderListB[i]    
    oldDistanceSum := orderB.KitchenDistanceList[0].Distance
    newDistanceSum := GetHaversineDistance(orderB.Coord,kitchenA.Coord)
    distanceProfit := oldDistanceSum-newDistanceSum
    if distanceProfit>0 {
	    swapInfo = append(swapInfo , SwapInfo{nil,orderB,distanceProfit,false})
    }			
  }			  
	sort.SliceStable(swapInfo,func (lhs,rhs int) bool {
	  return swapInfo[lhs].DistanceProfit>swapInfo[rhs].DistanceProfit
	})
	
	
	//choose wich swap will be done
	for i:=range swapInfo{	      
	  //fmt.Println(swapInfo[i].DistanceProfit)
	  if swapInfo[i].IsUsed{
	    continue
	  }
	  orderA:=swapInfo[i].OrderA
	  orderB:=swapInfo[i].OrderB
	  if orderB==nil{
  	  //if kitchenA is giving away order
	    if kitchenA.CanLetGo(orderA) && kitchenB.CanServe(orderA) {
	        kitchenA.GiveOrderToKitchen(orderA,kitchenB)
	        swapInfo[i].IsUsed = true
	    }
	  } else if orderA==nil{
  	  //if kitchenB is giving away order
	    if kitchenB.CanLetGo(orderB) && kitchenA.CanServe(orderB) {
	        kitchenB.GiveOrderToKitchen(orderB,kitchenA)
	        swapInfo[i].IsUsed = true
	    }
	  } else {
	    //if kitchens are swapping orders  
	    if OrdersCanSwapKitchen(orderA,orderB) {
	      kitchenA.GiveOrderToKitchen(orderA,kitchenB)
	      kitchenB.GiveOrderToKitchen(orderB,kitchenA)	      
        swapInfo[i].IsUsed = true
	    }
	  }
	  
	  //if swap happens
		if swapInfo[i].IsUsed {
		  //flag every swap info that contains orderA or orderB
		  for j:=i+1;j<len(swapInfo);j++{
		    if (orderA!=nil && swapInfo[j].OrderA!=nil) && (swapInfo[j].OrderA.Id == orderA.Id) {
		      swapInfo[j].IsUsed = true
		    }
		    if (orderB!=nil && swapInfo[j].OrderB!=nil) && (swapInfo[j].OrderB.Id == orderB.Id) {
		      swapInfo[j].IsUsed = true
		    }
		  }
		}
	}
}

func GreedyClusterOptimization(optimizationRepetition int){
  //repeat for as many as optimizationRepetition
  //for each pair of kitchen, let them swap orders
	for k:=1;k<=optimizationRepetition;k++{	  
		for i:=range KitchenList{
			for j:=range KitchenList{
				if i<j{
          GreedyMaximumOrderMatching(&KitchenList[i],&KitchenList[j])					
				}				
			}
		}
	}
}
