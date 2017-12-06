package main

import (
  "sort"
  //"fmt"
)

func GreedyClusteringForPlace(){  
  for i:= range PlaceList{
    place:=&PlaceList[i]
    place.ServingKitchen=&KitchenList[place.KitchenDistanceList[0].Index]
    kitchen:=place.ServingKitchen
    kitchen.OrderQty+=place.Qty
    kitchen.DistinctOrderCount+=len(place.OrderList)
    for j:= range place.OrderList{
      order:=place.OrderList[j]
      order.ServingKitchen=kitchen
    }
  } 
  ResolveMaximumCapacityViolationForPlace()  
  ResolveMinimumCapacityViolationForPlace()
} 


func ResolveMaximumCapacityViolationForPlace(){
  //For kitchens violating the maximum capacity :
  //Starting from the furthes place to the nearest
  //Give the place which is served by the kitchen
  //To the alternative kitchen nearest to that place and is able to serve that place
  //Do this until the kitchen does not violate the maximum capacity
  for i:= range KitchenList{
    kitchenOrigin:= &KitchenList[i]
    for j:= len(kitchenOrigin.PlaceDistanceList)-1; j>=0 &&  kitchenOrigin.OrderQty > kitchenOrigin.Capacity.Max;j--{
      placeToLetGo:=&PlaceList[kitchenOrigin.PlaceDistanceList[j].Index]
      if placeToLetGo.ServingKitchen!=kitchenOrigin{
        continue
      }      
      for k:= 1; k<len(placeToLetGo.KitchenDistanceList); k++{
        kitchenDestination:= &KitchenList[placeToLetGo.KitchenDistanceList[k].Index]
        if kitchenDestination.CanServePlace(placeToLetGo) {
          kitchenOrigin.GivePlaceToKitchen(placeToLetGo,kitchenDestination)
          break
        }
      }      
    }
  } 
}

func ResolveMinimumCapacityViolationForPlace(){  
  for i:= range KitchenList{
    thisKitchen:=&KitchenList[i]    
    if thisKitchen.OrderQty>=thisKitchen.Capacity.Min{
      continue
    }    
    //For kitchens violating the minimum capacity and has more than 50% of the minimum capacity:
    //Starting from the nearest place to the furthest
    //Get the place from the original kitchen serving the place
    //Only when that kitchen is able to let go the forementioned place without violating its minimum capacity
    //Do this until the kitchen does not violate the minimum capacity
    if thisKitchen.OrderQty>=int(float64(thisKitchen.Capacity.Min)*0.5){
      kitchenDestination:=thisKitchen
      for j:=0; j<len(kitchenDestination.PlaceDistanceList) &&  kitchenDestination.OrderQty < kitchenDestination.Capacity.Min;j++{
        placeToTakeOver :=&PlaceList[kitchenDestination.PlaceDistanceList[j].Index]
        kitchenOrigin   :=placeToTakeOver.ServingKitchen
        if (kitchenOrigin!=kitchenDestination) && kitchenDestination.CanLetGoPlace(placeToTakeOver){
          kitchenOrigin.GivePlaceToKitchen(placeToTakeOver,kitchenDestination)
        }
      }    
    } else {                 
    //if it has less than 50% capacity, distribute the orders it serves to other kitchens
    //nearest to that order, and can serve the order without violating maximum capacity
      kitchenOrigin:=thisKitchen
      for j:=0;j<len(kitchenOrigin.PlaceDistanceList) && kitchenOrigin.OrderQty>0;j++{
        placeToLetGo := &PlaceList[kitchenOrigin.PlaceDistanceList[j].Index]        
        if placeToLetGo.ServingKitchen!=kitchenOrigin{
          continue
        }
        for k:= 1; k<len(KitchenList); k++{
          kitchenDestination:= &KitchenList[placeToLetGo.KitchenDistanceList[k].Index]
          if kitchenDestination.CanServePlace(placeToLetGo) {
            kitchenOrigin.GivePlaceToKitchen(placeToLetGo,kitchenDestination)          
            break
          }
        }      
      }
    }
  } 
}

func GreedyMaximumOrderMatchingForPlace(kitchenA,kitchenB *Kitchen){
	var placeListA,placeListB []*Place
	type SwapInfo struct{
		PlaceA,PlaceB *Place
		DistanceProfit float64
		IsUsed bool
	}
	var swapInfo []SwapInfo	
	for i:= range PlaceList{
	  place:=&PlaceList[i]	  	  
	  servingKitchen:=place.ServingKitchen
		if servingKitchen == kitchenA {
			placeListA = append(placeListA, place)
		} else if servingKitchen == kitchenB {
			placeListB = append(placeListB, place)		  
		}		
	}
	
	
	for i :=range placeListA{
	  placeA:=placeListA[i]
	  //debug fmt.Println(orderA.Id,KitchenList[orderA.GetServerKitchenIndex()].Name)	  
		for j:=range placeListB {	
		  placeB:=placeListB[j]
		  //determine wether swapping the place make "distanceProfit" or not		  
			oldDistanceSum := GetHaversineDistance(placeA.Coord,kitchenA.Coord) + GetHaversineDistance(placeB.Coord,kitchenB.Coord)
			newDistanceSum := GetHaversineDistance(placeA.Coord,kitchenB.Coord) + GetHaversineDistance(placeB.Coord,kitchenA.Coord)
			distanceProfit := oldDistanceSum-newDistanceSum
			if distanceProfit>0 {
				swapInfo = append(swapInfo , SwapInfo{placeA,placeB,distanceProfit,false})
			}									
		}
	}
	
  //determine wether kitchen A giving its place to kitchen B make "distanceProfit"
  for i:= range placeListA{
    placeA:=placeListA[i]
    oldDistanceSum := GetHaversineDistance(placeA.Coord,kitchenA.Coord)
    newDistanceSum := GetHaversineDistance(placeA.Coord,kitchenB.Coord)
    distanceProfit := oldDistanceSum-newDistanceSum    
    if distanceProfit>0 {
	    swapInfo = append(swapInfo , SwapInfo{placeA,nil,distanceProfit,false})
    }			
  }	
  
  //determine wether kitchen B giving its place to kitchen A make "distanceProfit"
  for i:= range placeListB{
    placeB:= placeListB[i]    
    oldDistanceSum := GetHaversineDistance(placeB.Coord,kitchenB.Coord)
    newDistanceSum := GetHaversineDistance(placeB.Coord,kitchenA.Coord)
    distanceProfit := oldDistanceSum-newDistanceSum
    if distanceProfit>0 {
	    swapInfo = append(swapInfo , SwapInfo{nil,placeB,distanceProfit,false})
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
	  placeA:=swapInfo[i].PlaceA
	  placeB:=swapInfo[i].PlaceB
	  if placeB==nil{
  	  //if kitchenA is giving away order
	    if kitchenA.CanLetGoPlace(placeA) && kitchenB.CanServePlace(placeA) {
	        kitchenA.GivePlaceToKitchen(placeA,kitchenB)
	        swapInfo[i].IsUsed = true
	    }
	  } else if placeA==nil{
  	  //if kitchenB is giving away order
	    if kitchenB.CanLetGoPlace(placeB) && kitchenA.CanServePlace(placeB) {
	        kitchenB.GivePlaceToKitchen(placeB,kitchenA)
	        swapInfo[i].IsUsed = true
	    }
	  } else {
	    //if kitchens are swapping orders  
	    if PlacesCanSwapKitchen(placeA,placeB) {
	      kitchenA.GivePlaceToKitchen(placeA,kitchenB)
	      kitchenB.GivePlaceToKitchen(placeB,kitchenA)	      
        swapInfo[i].IsUsed = true
	    }
	  }
	  
	  //if swap happens
		if swapInfo[i].IsUsed {
		  //flag every swap info that contains orderA or orderB
		  for j:=i+1;j<len(swapInfo);j++{
		    if (placeA!=nil && swapInfo[j].PlaceA!=nil) && (swapInfo[j].PlaceA.Id == placeA.Id) {
		      swapInfo[j].IsUsed = true
		    }
		    if (placeB!=nil && swapInfo[j].PlaceB!=nil) && (swapInfo[j].PlaceB.Id == placeB.Id) {
		      swapInfo[j].IsUsed = true
		    }
		  }
		}
	}
}

func GreedyClusterOptimizationForPlace(optimizationRepetition int){
  //repeat for as many as optimizationRepetition
  //for each pair of kitchen, let them swap orders
	for k:=1;k<=optimizationRepetition;k++{	  
		for i:=range KitchenList{
			for j:=range KitchenList{
				if i<j{
          GreedyMaximumOrderMatchingForPlace(&KitchenList[i],&KitchenList[j])					
				}				
			}
		}
	}
}
