package main
import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"  
  "sort"
  "net/http"
  "log"
)

var kitchenList []Kitchen  
var orderList []Order
func main(){  
  //get handle or connection for kulina database
  db,err := sql.Open("mysql","gemilang@tcp(127.0.0.1:3306)/kulina?parseTime=true")  
  defer db.Close()  
  checkErr(err)  
  err = db.Ping()
  checkErr(err)
    
  // get kitchens' capacities from database and put into kitchen_list
  chosenDate:="2017-11-20"
  rows,err:= db.Query("select k.id, k.name, kc.min_capacity, kc.max_capacity, kc.optimum_capacity, k.loc_lat, k.loc_lon from kitchen_capacities kc join kitchens k 	on k.id = kc.kitchen_id where date = \""+chosenDate+"\";") 
  checkErr(err)  
  for rows.Next(){
    var kitchen Kitchen
    kitchen.ScanFromSQL(rows)
    kitchenList = append(kitchenList,kitchen)
  }

  //get orders from database and put into kitchen_lists
  rows,err= db.Query("select o.id, od.qty, uda.latitude, uda.longitude from orders_delivery od join orders o on o.id = od.order_id join user_delivery_addresses uda on uda.user_id = o.user_id where od.delivery_date = \""+chosenDate+"\" and o.status = 1 and o.start_date <=\""+chosenDate+"\"	and o.end_date >= \""+chosenDate+"\";") 
  checkErr(err)
  for rows.Next(){
    var order Order
    order.ScanFromSQL(rows)
    orderList = append(orderList,order)
  }


  //building haversine distance list of kitchens to orders, and orders to kitchen, and orders to orders
  for i:=range(kitchenList){
    for j:= range(orderList){
      dist:=Haversine(kitchenList[i].Loc_lat,kitchenList[i].Loc_lon,orderList[j].Latitude,orderList[j].Longitude)       
      kitchenList[i].Order_distance   = append(kitchenList[i].Order_distance,Distance{j,dist})
      orderList[j].Kitchen_distance   = append(orderList[j].Kitchen_distance,Distance{i,dist})
    }    
  }            
  for i:=0;i<len(orderList);i++{
    for j:=i+1;j<len(orderList);j++{   
       dist:=Haversine(orderList[i].Latitude,orderList[i].Longitude,orderList[j].Latitude,orderList[j].Longitude)      
       orderList[i].Order_distance   = append(orderList[i].Order_distance,Distance{j,dist})
    }
  }
  for i:=range(kitchenList){
    sort.SliceStable(kitchenList[i].Order_distance,func(lhs,rhs int) bool {
      return kitchenList[i].Order_distance[lhs].Distance<kitchenList[i].Order_distance[rhs].Distance
    })
  }
  for i:=range(orderList){
    sort.SliceStable(orderList[i].Kitchen_distance,func(lhs,rhs int) bool {
      return orderList[i].Kitchen_distance[lhs].Distance<orderList[i].Kitchen_distance[rhs].Distance
    })
    sort.SliceStable(orderList[i].Order_distance,func(lhs,rhs int) bool {
      return orderList[i].Order_distance[lhs].Distance<orderList[i].Order_distance[rhs].Distance
    })
  }
  
  //clustering!
  orderList,kitchenList = NaiveClustering(orderList,kitchenList)
  for i:=range(kitchenList){
    fmt.Println(kitchenList[i].Id,kitchenList[i].Min_capacity,kitchenList[i].Max_capacity,kitchenList[i].Order_qty)
  }  
  
  //do the maximum_matching
  
  orderList,kitchenList = NaiveOptimizeCluster(orderList,kitchenList,10)
  fmt.Println("-----------------------------")
  for i:=range(kitchenList){
    fmt.Println(kitchenList[i].Id,kitchenList[i].Min_capacity,kitchenList[i].Max_capacity,kitchenList[i].Order_qty)
  }  
     
  http.Handle("/", &templateHandler{filename: "GoogleMap.html"})
  // start the web server
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }  
}
