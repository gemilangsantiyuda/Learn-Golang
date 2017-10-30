package main

import (
	"math"
	"sort"
  "gonum.org/v1/plot"
  "gonum.org/v1/plot/plotter"
  "gonum.org/v1/plot/vg"
  "image/color"
)

type ForDistance struct{
	Idx int
	Dist float64
}

type Kitchen struct{
	Name string	
	Lon,Lat float64 //coordinate on earth
	MinCap,MaxCap,CurrentCap,CurrentReq,CurrentPri int
	Distance []ForDistance
}
type Customer struct{
	Name string	
	Lon,Lat float64//coordinate on earth
	Req int
  Visited bool
  DistanceFromCenter float64
	Distance []ForDistance	
	Priority []ForDistance
}
type Rider struct{
  IdxPath []int
  NamePath []string
  PathLength float64
  PathTime float64
  CurrentCap int
}

func (r *Rider) GiveNamePath(customerList []Customer, kitchenList []Kitchen){
  r.NamePath = make([]string,len(r.IdxPath))
  r.NamePath[0] = kitchenList[r.IdxPath[0]].Name
  for i:=1 ; i<len(r.NamePath); i++{
    r.NamePath[i] = customerList[r.IdxPath[i]].Name
  }
} 

func (r *Rider) CalculatePathLength(customerList []Customer, kitchenList []Kitchen){
  length := 0.000000
  length += Haversine(kitchenList[r.IdxPath[0]].Lat,kitchenList[r.IdxPath[0]].Lon,customerList[r.IdxPath[1]].Lat,customerList[r.IdxPath[1]].Lon)
  for i:=1 ; i<len(r.IdxPath)-1; i++ {
    c1 := customerList[r.IdxPath[i]]
    c2 := customerList[r.IdxPath[i+1]]
    length+=Haversine(c1.Lat,c1.Lon,c2.Lat,c2.Lon)
  }
  r.PathLength = length
}

func (r *Rider) CalculatePathTime(){
  r.PathTime = r.PathLength*3.00000/60.000000
}

func (r *Rider) MakePath(customerList []Customer, kitchenList []Kitchen) ([]Customer,[]Kitchen) {
  MAX_DISTANCE := 60.000000000
  MAX_CAP := 40
  reqLeft := 0
  maxPri :=-1
  chosenKitchen :=-1
  for i:= range kitchenList {
    reqLeft+= kitchenList[i].CurrentReq
    if maxPri<kitchenList[i].CurrentPri {
      maxPri = kitchenList[i].CurrentPri
      chosenKitchen = i
    }
  }
  if reqLeft == 0 {
    return customerList,kitchenList
  }
  r.IdxPath = append(r.IdxPath, chosenKitchen)
  for i:= range kitchenList[chosenKitchen].Distance {
    customerIdx := kitchenList[chosenKitchen].Distance[i].Idx
    if !customerList[customerIdx].Visited && customerList[customerIdx].Priority[0].Idx == chosenKitchen && kitchenList[chosenKitchen].CurrentCap + customerList[customerIdx].Req <=kitchenList[chosenKitchen].MaxCap {
      r.IdxPath = append(r.IdxPath,customerIdx)
      r.PathLength += kitchenList[chosenKitchen].Distance[i].Dist
      r.CurrentCap += customerList[customerIdx].Req 
      customerList[customerIdx].Visited = true
      kitchenList[chosenKitchen].CurrentPri--
      kitchenList[chosenKitchen].CurrentReq -= customerList[customerIdx].Req
      kitchenList[chosenKitchen].CurrentCap += customerList[customerIdx].Req
      break
    }
  }
  //continue to build the path , including customer of the same cluster
  customerNow := r.IdxPath[1]
  found := true 
  for found{
    found = false
    for i:= range customerList[customerNow].Distance {
      customerIdx := customerList[customerNow].Distance[i].Idx
      dist := customerList[customerNow].Distance[i].Dist
      if !customerList[customerIdx].Visited && customerList[customerIdx].Priority[0].Idx == chosenKitchen && kitchenList[chosenKitchen].CurrentCap + customerList[customerIdx].Req <=kitchenList[chosenKitchen].MaxCap && r.PathLength + dist <=MAX_DISTANCE && r.CurrentCap + customerList[customerIdx].Req <= MAX_CAP{
        r.IdxPath = append(r.IdxPath,customerIdx)
        r.PathLength += dist
        r.CurrentCap += customerList[customerIdx].Req 
        customerList[customerIdx].Visited = true
        kitchenList[chosenKitchen].CurrentPri--
        kitchenList[chosenKitchen].CurrentReq -= customerList[customerIdx].Req
        kitchenList[chosenKitchen].CurrentCap += customerList[customerIdx].Req
        found = true
        customerNow = customerIdx
        break
      }
    }
  }
  //continue to build the path, including any other customer there is
  found = true
  for found{
    found = false
    for i:= range customerList[customerNow].Distance {
      customerIdx := customerList[customerNow].Distance[i].Idx
      dist := customerList[customerNow].Distance[i].Dist
      stolenKitchen := customerList[customerNow].Priority[0].Idx
      if !customerList[customerIdx].Visited && kitchenList[stolenKitchen].CurrentReq + kitchenList[stolenKitchen].CurrentCap - customerList[customerIdx].Req >= kitchenList[stolenKitchen].MinCap && kitchenList[chosenKitchen].CurrentCap + kitchenList[chosenKitchen].CurrentReq + customerList[customerIdx].Req <= kitchenList[chosenKitchen].MaxCap && r.PathLength + dist <=MAX_DISTANCE && r.CurrentCap + customerList[customerIdx].Req <= MAX_CAP {
        r.IdxPath = append(r.IdxPath,customerIdx)
        r.PathLength += dist
        r.CurrentCap += customerList[customerIdx].Req 
        customerList[customerIdx].Visited = true
        kitchenList[stolenKitchen].CurrentPri--
        kitchenList[stolenKitchen].CurrentReq -= customerList[customerIdx].Req
        kitchenList[chosenKitchen].CurrentCap += customerList[customerIdx].Req
        found = true
        customerNow = customerIdx
        break
      }
    }
  }
  return customerList,kitchenList
}


//all variable declared here are for optimization purpose
var vis  [3000000][22]uint8
var next [3000000][22]int
var memo [3000000][22]float64
var counter uint8
var nodes [][2]float64
var pathN int

func (r *Rider) OptimizePath(customerList []Customer,kitchenList []Kitchen) {

  nodes = nil  
  pathN = len(r.IdxPath)
  counter++
  for i:= range(r.IdxPath){
    if i==0 {
      nodes = append(nodes, [2]float64{kitchenList[r.IdxPath[0]].Lon, kitchenList[r.IdxPath[0]].Lat})
    } else {
      nodes = append(nodes, [2]float64{customerList[r.IdxPath[i]].Lon,customerList[r.IdxPath[i]].Lat})  
    }
  }
  var bestPath []int
  bestLength := 9999.0000000
  bestNext := 1
  for i:= 1; i < len(r.IdxPath);i++{
    dist := Haversine(nodes[0][1],nodes[0][0],nodes[i][1],nodes[i][0])
    tmp := dist + dp(1 + int(math.Pow(2,float64(i))) , i)
    if tmp<bestLength {
      bestNext = i
      bestLength = tmp
    }
  }
  bitm := 1
  bestPath = append(bestPath,r.IdxPath[0])
  for bitm< int(math.Pow(2,float64(pathN)))-1 {
    bestPath = append(bestPath,r.IdxPath[bestNext])
    bitm += int(math.Pow(2,float64(bestNext)))
    bestNext = next[bitm][bestNext]
  }
  r.IdxPath = bestPath
  r.CalculatePathLength(customerList,kitchenList)
}

func dp(bitmask ,now int) float64 {
  if bitmask == int(math.Pow(2,float64(pathN)))-1 {
    return 0
  }
  if vis[bitmask][now] == counter {
    return memo[bitmask][now]
  }
  vis[bitmask][now] = counter
  bestLength := 9999.000000
  bestNext := 1
  for i:=1 ; i<pathN; i++ {
    if int(math.Pow(2,float64(i))) & bitmask != 0 {
      continue
    } 
    bitm := bitmask + int(math.Pow(2,float64(i)))
    dist := Haversine(nodes[now][1],nodes[now][0],nodes[i][1],nodes[i][0])
    tmp := dp(bitm,i) + dist
    if tmp< bestLength {
      bestLength = tmp
      bestNext = i
    }
  }
  memo[bitmask][now]=bestLength
  next[bitmask][now]=bestNext
  return bestLength
} 

func PlotDistribution(customerList []Customer,kitchenList []Kitchen, fileName string){
  colors := [][]uint8{{255,255,0},{255,255,255},{255,0,0},{0,255,0},{0,0,255},{255,0,255},{0,255,255}}
  p, err := plot.New() 
  if err!= nil {
    panic(err)
  }
  for i:= range kitchenList{
  	ptsK := make(plotter.XYZs,1)
  	ptsK[0].X = kitchenList[i].Lon
  	ptsK[0].Y = kitchenList[i].Lat
		ptsK[0].Z = 9
		
  	bs,err := plotter.NewBubbles(ptsK,vg.Points(15),vg.Points(15))
    if err!= nil {
      panic(err)    
    }
    
    bs.Color = color.RGBA{R:colors[i][0],G:colors[i][1],B:colors[i][2]}
    p.Add(bs)
    
    pts := make(plotter.XYs,kitchenList[i].CurrentPri)
    ptsN := 0
    for j:= range customerList{
      if customerList[j].Priority[0].Idx == i {
        pts[ptsN].X = customerList[j].Lon
        pts[ptsN].Y = customerList[j].Lat
        ptsN++
      }
    }
    s,err := plotter.NewScatter(pts)
    if err!= nil {
      panic(err)    
    }
    s.GlyphStyle.Color = color.RGBA{R:colors[i][0],G:colors[i][1],B:colors[i][2]}
    s.GlyphStyle.Radius = vg.Points(12)
    p.Add(s)
  }
  if err:= p.Save(20*vg.Inch,20*vg.Inch,fileName+".png"); err != nil{
    panic(err)  
  }
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

func NormalMoreMax(customerList []Customer, kitchenList []Kitchen) ([]Customer,[]Kitchen){
  for i:= range kitchenList{
    for j:= len(kitchenList[i].Distance)-1; j>=0 &&  kitchenList[i].CurrentReq > kitchenList[i].MaxCap;j--{
      customerIdx:=kitchenList[i].Distance[j].Idx
      if customerList[customerIdx].Priority[0].Idx!=i{
        continue
      }
      chosenKitchen := -1
      for k:= 1; k<len(kitchenList); k++{
        kitchenSubIdx:= customerList[customerIdx].Priority[k].Idx
        if kitchenList[kitchenSubIdx].CurrentReq + customerList[customerIdx].Req <= kitchenList[kitchenSubIdx].MaxCap {
          chosenKitchen = k
          kitchenList[kitchenSubIdx].CurrentReq += customerList[customerIdx].Req
          kitchenList[kitchenSubIdx].CurrentPri ++
          kitchenList[i].CurrentPri--
          kitchenList[i].CurrentReq -= customerList[customerIdx].Req
          break
        }
      }
      if chosenKitchen==-1{
        continue      
      }
      for k:=chosenKitchen;k>0;k--{
        tmp:=customerList[customerIdx].Priority[k]
        customerList[customerIdx].Priority[k] = customerList[customerIdx].Priority[k-1] 
        customerList[customerIdx].Priority[k-1] = tmp
      }

    }
  } 
  return customerList,kitchenList
}
func (c *Customer) CalculateDistance(customerList []Customer){
	for i:=0;i<len(customerList);i++{
		d :=customerList[i]
		distance := Haversine(c.Lat,c.Lon,d.Lat,d.Lon)
		c.Distance = append(c.Distance, ForDistance{i,distance})
	}
  sort.SliceStable(c.Distance,func (lhs,rhs int) bool {return c.Distance[lhs].Dist<c.Distance[rhs].Dist})
}
func (k *Kitchen) CalculateDistance(customerList []Customer){
	for i:=0;i<len(customerList);i++{
		d := customerList[i]
		distance := Haversine(k.Lat,k.Lon,d.Lat,d.Lon)
		k.Distance = append(k.Distance,ForDistance{i,distance})
	}
  sort.SliceStable(k.Distance,func (lhs,rhs int) bool {return k.Distance[lhs].Dist<k.Distance[rhs].Dist})
}
func NaiveClustering(customerList []Customer, kitchenList []Kitchen) ([]Customer,[]Kitchen){
	//Prioritize Customers to nearest Kitchens
	for i:= range customerList{
		for j:= range kitchenList{
			c := customerList[i]
			k := kitchenList[j]
			distance := Haversine(c.Lat,c.Lon,k.Lat,k.Lon)
			customerList[i].Priority = append(customerList[i].Priority,ForDistance{j,distance})
		}
    sort.SliceStable(customerList[i].Priority,func (lhs,rhs int) bool {return customerList[i].Priority[lhs].Dist<customerList[i].Priority[rhs].Dist})
	}
  for i:= range customerList{
    kitchenIdx:= customerList[i].Priority[0].Idx
    kitchenList[kitchenIdx].CurrentPri++
    kitchenList[kitchenIdx].CurrentReq+=customerList[i].Req     
  }

  customerList,kitchenList = NormalMoreMax(customerList,kitchenList)

	return customerList,kitchenList
}

func NaiveOptimizeClusterGreedy(a int,b int ,customerList []Customer,kitchenList []Kitchen) ([]Customer, []Kitchen){
	type CustTmp struct{
		Idx int
		Lon,Lat float64
		Req int 
	}
	type SwapInfo struct{
		IdxA,IdxB int
		Dist float64
		ReqToA,ReqToB int
		used bool
	}
	var customerA,customerB []CustTmp
	var swapInfo []SwapInfo
	//var swapChosen []int
	
	for i:= range customerList{
		if customerList[i].Priority[0].Idx == a {
			customerA = append(customerA, CustTmp{i,customerList[i].Lon,customerList[i].Lat,customerList[i].Req})
		}
		if customerList[i].Priority[0].Idx == b {
			customerB = append(customerB, CustTmp{i,customerList[i].Lon,customerList[i].Lat,customerList[i].Req})
		}
	}
	//fmt.Println(customerA)
	//fmt.Println(customerB)
	
	for i :=range customerA{
		for j:=range customerB {
			distOld := customerList[customerA[i].Idx].Priority[0].Dist + customerList[customerB[j].Idx].Priority[0].Dist
			distNew := Haversine(customerA[i].Lat,customerA[i].Lon,kitchenList[b].Lat,kitchenList[b].Lon)
			distNew += Haversine(customerB[j].Lat,customerB[j].Lon,kitchenList[a].Lat,kitchenList[a].Lon)
			dist := distOld-distNew
			req := -customerA[i].Req + customerB[j].Req
			if dist>0 {
				swapInfo = append(swapInfo , SwapInfo{customerA[i].Idx,customerB[j].Idx,dist,req,-req,false})
			}
			
		}
	}
	sort.SliceStable(swapInfo,func (lhs,rhs int) bool {return swapInfo[lhs].Dist>swapInfo[rhs].Dist})
	//fmt.Println(swapInfo)
	
	for i:=range swapInfo{
		if !swapInfo[i].used {
			custA:=swapInfo[i].IdxA
			custB:=swapInfo[i].IdxB
			kitchenA:=a
			kitchenB:=b		
			reqToA:=swapInfo[i].ReqToA
			reqToB:=swapInfo[i].ReqToB
			if (kitchenList[kitchenA].CurrentReq+reqToA >= kitchenList[kitchenA].MinCap) && (kitchenList[kitchenA].CurrentReq+reqToA <= kitchenList[kitchenA].MaxCap) && (kitchenList[kitchenB].CurrentReq+reqToB >= kitchenList[kitchenB].MinCap) && (kitchenList[kitchenB].CurrentReq+reqToB <= kitchenList[kitchenB].MaxCap) {
				
				swapInfo[i].used = true
				kitchenList[kitchenA].CurrentReq+=reqToA
				kitchenList[kitchenB].CurrentReq+=reqToB
				
				chosenKitchen := -1
				for k:= range kitchenList{
					if customerList[custA].Priority[k].Idx == kitchenB {
						chosenKitchen = k
						break
					}				
				}
				for k:= chosenKitchen; k>0;k--{
					tmp:= customerList[custA].Priority[k]
					customerList[custA].Priority[k] = customerList[custA].Priority[k-1]
					customerList[custA].Priority[k-1] = tmp
				}
				
				chosenKitchen = -1
				for k:= range kitchenList{
					if customerList[custB].Priority[k].Idx == kitchenA {
						chosenKitchen = k
						break
					}				
				}
				for k:= chosenKitchen; k>0;k--{
					tmp:= customerList[custB].Priority[k]
					customerList[custB].Priority[k] = customerList[custB].Priority[k-1]
					customerList[custB].Priority[k-1] = tmp
				}
				
				for j:= i+1;j<len(swapInfo);j++{
					if swapInfo[j].IdxA == custA || swapInfo[j].IdxB == custB {
						swapInfo[j].used=true
					}
				}
			}
		}
	}	
	return customerList,kitchenList
}

func NaiveOptimizeCluster(customerList []Customer, kitchenList []Kitchen) ([]Customer, []Kitchen){
	for k:=1;k<=1;k++{
		for i:=range kitchenList{
			for j:=range kitchenList{
				if i<j{
					customerList,kitchenList = NaiveOptimizeClusterGreedy(i,j,customerList,kitchenList)
				}
			}
		}
	}
	return customerList,kitchenList
}
