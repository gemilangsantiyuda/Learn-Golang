package main

import(
  "fmt"
  "encoding/csv"
  "os"
  "strconv"
  "bufio"
  "math"
)


var kitchenList []Kitchen
var customerList []Customer
var riderList []Rider
var message []string


func main(){

//Read file csv of kitchen and store it to Kitchens array
	fileKitchen,_ := os.Open("TSP_kulina_dapur")
	r:=csv.NewReader(bufio.NewReader(fileKitchen))
	for i:=0;;i++{
		record,err := r.Read()
		if i==0{
			continue		
		}
		if err!= nil{
			break		
		}
		var k Kitchen;
		k.Name 		= record[0]
		k.Lon,_ 	= strconv.ParseFloat(record[1],64)
		k.Lat,_ 	= strconv.ParseFloat(record[2],64)
		k.MinCap,_	= strconv.Atoi(record[3])
		k.MaxCap,_	= strconv.Atoi(record[5])
		kitchenList = append(kitchenList,k)
	}
//Read file csv of customers and store it to Customers array
	fileCustomer,_ := os.Open("TSP_kulina_kust")
	r=csv.NewReader(bufio.NewReader(fileCustomer))
	for i:=0;;i++{
		record,err := r.Read()
		if i==0{
			continue		
		}
		if err!= nil{
			break		
		}
		var c Customer;
		c.Name		= record[0]
		c.Lon,_		= strconv.ParseFloat(record[1],64)
		c.Lat,_		= strconv.ParseFloat(record[2],64)
		c.Req,_		= strconv.Atoi(record[3])
		customerList = append(customerList,c)
	}		
//build Distance for Kitchens and Customers
	for i:=range(kitchenList){
		kitchenList[i].CalculateDistance(customerList)
	}
	for i:=range(customerList){
		customerList[i].CalculateDistance(customerList)	
	}
  //Cluster the Customers into their most suitable Kitchens
	customerList,kitchenList:=	NaiveClustering(customerList,kitchenList)
  //swap customers among each pair of kitchens	
 	//customerList,kitchenList= NaiveOptimizeCluster(customerList,kitchenList)

//Making Path
  totalTime:=0
  totalLength:=0.000
  for true {
    var rider Rider
    customerList,kitchenList = rider.MakePath(customerList,kitchenList)
    if rider.PathLength == 0 {
      break  
    }
    rider.GiveNamePath(customerList,kitchenList)
    fmt.Println("Rider ",len(riderList)+1,"dengan Path :\n",rider.NamePath,"\nMembawa ",rider.CurrentCap," makanan\ndengan Jarak Tempuh: ",rider.PathLength)
    fmt.Println("-----------------------------------------------------------------")    
    totalTime+= int(math.Ceil(rider.PathLength/20))
    totalLength+=rider.PathLength
    riderList = append(riderList,rider)
  }
  fmt.Println("Total jarak tempuh :",totalLength,"\nTotal waktu yang dibayar : ",totalTime)
  fmt.Println("\nOptimasi Dilakukan :")
  for i:= range riderList{
    if len(riderList[i].IdxPath)<=21{
      riderList[i].OptimizePath(customerList,kitchenList)
    }
    riderList[i].GiveNamePath(customerList,kitchenList)
    fmt.Println("Rider ",i,"dengan Path :\n",riderList[i].NamePath,"\nMembawa ",riderList[i].CurrentCap," makanan\ndengan Jarak Tempuh: ",riderList[i].PathLength)
    fmt.Println("-----------------------------------------------------------------")    
  }
  totalTime=0
  totalLength=0.000
  for i:=range riderList{
  	totalTime+= int(math.Ceil(riderList[i].PathLength/20))
  	totalLength+= riderList[i].PathLength
  }
  fmt.Println("Total jarak tempuh :",totalLength,"\nTotal waktu yang dibayar : ",totalTime)
}
