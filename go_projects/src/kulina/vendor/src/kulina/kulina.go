package kulina

import (
	"math"
	"fmt"
	"sort"
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
	Distance []ForDistance	
	Priority []ForDistance
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

func (c *Customer) CalculateDistance(CustomerList []Customer){
	for i:=0;i<len(CustomerList);i++{
		d :=CustomerList[i]
		distance := Haversine(c.Lat,c.Lon,d.Lat,d.Lon)
		c.Distance = append(c.Distance, ForDistance{i,distance})
	}
}
func (k *Kitchen) CalculateDistance(CustomerList []Customer){
	for i:=0;i<len(CustomerList);i++{
		d := CustomerList[i]
		distance := Haversine(k.Lat,k.Lon,d.Lat,d.Lon)
		k.Distance = append(k.Distance,ForDistance{i,distance})
	}
}
func NaiveClustering(CustomerList []Customer, KitchenList []Kitchen) ([]Customer,[]Kitchen){
	//Prioritize Customers to nearest Kitchens
	for i:= range CustomerList{
		for j:= range KitchenList{
			c := CustomerList[i]
			k := KitchenList[j]
			distance := Haversine(c.Lat,c.Lon,k.Lat,k.Lon)
			CustomerList[i].Priority = append(CustomerList[i].Priority,ForDistance{j,distance})
		}
	}
	fmt.Print(CustomerList[2].Priority)
	for i:= range CustomerList{
			
	}	
	return CustomerList,KitchenList
}
