package main
import (
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
  "math"
  "path/filepath"
  "text/template"
  "net/http"
  "sort"
  "encoding/json"
  "time"
  "io/ioutil"
  "log"
  "strconv"
  "fmt"
)

func checkErr(err error){
  if err!=nil{
    panic(err.Error())
  }
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.once.Do(func() {
    t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
  })
  data := DataTemplate {kitchenList,orderList }
  t.templ.Execute(w,data)
}

func (k *Kitchen) ScanFromSQL(row *sql.Rows){
  var k1 KitchenScan
  err:=row.Scan(&k1.id,&k1.name,&k1.min_capacity,&k1.max_capacity,&k1.optimum_capacity,&k1.loc_lat,&k1.loc_lon) 
  checkErr(err)
  k.Id=int(k1.id.Int64)
  k.Name=k1.name.String
  k.Min_capacity=int(k1.min_capacity.Int64)
  k.Max_capacity=int(k1.max_capacity.Int64)
  k.Optimum_capacity=int(k1.optimum_capacity.Int64)
  k.Loc_lat=k1.loc_lat.Float64
  k.Loc_lon=k1.loc_lon.Float64
}

func (order *Order) ScanFromSQL(row *sql.Rows){
  var order1 OrderScan
  err:=row.Scan(&order1.id,&order1.qty,&order1.latitude,&order1.longitude)  
  checkErr(err)
  order.Id=int(order1.id.Int64)
  order.Qty=int(order1.qty.Int64)
  order.Latitude=order1.latitude.Float64
  order.Longitude=order1.longitude.Float64
}

func Haversine(lat1,lon1,lat2,lon2 float64) float64 {
	DY := math.Abs(lat1-lat2)/180*math.Pi
	DX := math.Abs(lon1-lon2)/180*math.Pi
	Y1 := lat1/180*math.Pi
	Y2 := lat2/180*math.Pi
	R := 6371.00000
	a := math.Sin(DY/2)*math.Sin(DY/2) + math.Cos(Y1)*math.Cos(Y2)*math.Sin(DX/2)*math.Sin(DX/2)
	c := 2*math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	return R*c
}

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

func BuildPathList (orderList []Order, kitchenList []Kitchen) []Path{  
  orderLeft := len(orderList)  
  
  var pathList []Path
  
  for orderLeft>0 && len(pathList)<1 {    
    chosenKitchen:=-1
    maxOrder:=0
    for i:=range(kitchenList){
      if kitchenList[i].Order_count >maxOrder{
        maxOrder = kitchenList[i].Order_count
        chosenKitchen = i
      }
    }
    
    newPath := BuildPath(orderList,kitchenList,chosenKitchen)
    newPath.Id = len(pathList)
    newPath.GiveName(orderList,kitchenList)
    pathList = append(pathList,newPath)
    
    orderLeft = 0
    for i:=range(kitchenList){
      orderLeft+= kitchenList[i].Order_count
    }
    fmt.Println(orderLeft)
  }
  return pathList
}

func (p *Path) GiveName(orderList []Order, kitchenList []Kitchen){
  p.Path_name = make([]string,len(p.Path_idx))
  p.Path_name[0] = kitchenList[p.Path_idx[0]].Name
  for i:=1 ; i<len(p.Path_name); i++{
    p.Path_name[i] = strconv.Itoa(orderList[p.Path_idx[i]].Id)
  }
} 


func GetGoogleDistance(lat1,lon1,lat2,lon2 float64) int{
  
  //convert lat and lon to string for google api query
  slat1:= strconv.FormatFloat(lat1, 'f', 7, 64)
  slon1:= strconv.FormatFloat(lon1, 'f', 7, 64)
  slat2:= strconv.FormatFloat(lat2, 'f', 7, 64)
  slon2:= strconv.FormatFloat(lon2, 'f', 7, 64)     
  
  //fmt.Println(slat1,slon1,slat2,slon2)
  url := "https://maps.googleapis.com/maps/api/distancematrix/json?origins="+slat1+","+slon1+"&destinations="+slat2+","+slon2+"&departure_time=1541202457&traffic_model=best_guess&key="+googleDistanceKey
  //fmt.Println(url)

  spaceClient := http.Client{
        Timeout: time.Second * 5, // Maximum of 2 secs
  }
  req, err := http.NewRequest(http.MethodGet, url, nil)
  if err != nil {
      log.Fatal(err)
  }

  req.Header.Set("User-Agent", "kulina-matching-beta")

  res, getErr := spaceClient.Do(req)
  if getErr != nil {
      log.Fatal(getErr)
  }

  body, readErr := ioutil.ReadAll(res.Body)
  if readErr != nil {
      log.Fatal(readErr)
  }

  var googleDistance GoogleDistanceResponse
  jsonErr := json.Unmarshal(body, &googleDistance)
  if jsonErr != nil {
      log.Fatal(jsonErr)
  }
  return googleDistance.Rows[0].Elements[0].Distance.Value
}

func GetNearestOrderDistance(lat,lon float64 ,path Path, kitchenIdx int, kitchen Kitchen, orderList []Order) Distance {
  var tmpList []Distance

  fmt.Println(kitchen.Name,kitchen.Order_qty,kitchen.Max_capacity)
  //get 5 haversine closest order, and from the same kitchen
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
