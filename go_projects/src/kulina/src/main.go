	package main

import(
	"fmt"
	"encoding/csv"
	"os"
	"bufio"
	"strconv"
  "src/kulina"
)


var KitchenList []kulina.Kitchen
var CustomerList []kulina.Customer


func main(){
//Read file csv of kitchen and store it to Kitchens array
	FileKitchen,_ := os.Open("TSP_kulina_dapur")
	r:=csv.NewReader(bufio.NewReader(FileKitchen))
	for i:=0;;i++{
		record,err := r.Read()
		if i==0{
			continue		
		}
		if err!= nil{
			break		
		}
		var k kulina.Kitchen;
		k.Name 		= record[0]
		k.Lon,_ 	= strconv.ParseFloat(record[1],64)
		k.Lat,_ 	= strconv.ParseFloat(record[2],64)
		k.MinCap,_	= strconv.Atoi(record[3])
		k.MaxCap,_	= strconv.Atoi(record[5])
		KitchenList = append(KitchenList,k)
	}
//Read file csv of customers and store it to Customers array
	FileCustomer,_ := os.Open("TSP_kulina_kust")
	r=csv.NewReader(bufio.NewReader(FileCustomer))
	for i:=0;;i++{
		record,err := r.Read()
		if i==0{
			continue		
		}
		if err!= nil{
			break		
		}
		var c kulina.Customer;
		c.Name		= record[0]
		c.Lon,_		= strconv.ParseFloat(record[1],64)
		c.Lat,_		= strconv.ParseFloat(record[2],64)
		c.Req,_		= strconv.Atoi(record[3])
		CustomerList = append(CustomerList,c)
	}		
//build Distance for Kitchens and Customers
	for i:=range(KitchenList){
		KitchenList[i].CalculateDistance(CustomerList)
	}
	for i:=range(CustomerList){
		CustomerList[i].CalculateDistance(CustomerList)	
	}
//Cluster the Customers into their most suitable Kitchens
	a,b:=	kulina.NaiveClustering(CustomerList,KitchenList)
	fmt.Print(" ",len(a)," ",len(b))
}
