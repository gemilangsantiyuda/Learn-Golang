package main
import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"  
)

func GetDatabaseHandle() *sql.DB {
  db,err := sql.Open("mysql","gemilang@tcp(127.0.0.1:3306)/kulina?parseTime=true")  
  checkErr(err)  
  err = db.Ping()
  checkErr(err)
  return db
}

func GetKitchensFromDatabase(db *sql.DB,chosenDate string) []Kitchen {
  var kitchenList []Kitchen
  rows,err:= db.Query("select k.id, k.name, kc.min_capacity, kc.max_capacity, kc.optimum_capacity, k.loc_lat, k.loc_lon from kitchen_capacities kc join kitchens k 	on k.id = kc.kitchen_id where date = \""+chosenDate+"\";") 
  checkErr(err)  
  for rows.Next(){
    var kitchen Kitchen
    kitchen.ScanFromSQL(rows)
    kitchenList = append(kitchenList,kitchen)
  }
  return kitchenList
}
  
func GetOrdersFromDatabase(db *sql.DB , chosenDate string) []Order {
  var orderList []Order
  rows,err:= db.Query("select o.id, od.qty, uda.latitude, uda.longitude from orders_delivery od join orders o on o.id = od.order_id join user_delivery_addresses uda on uda.user_id = o.user_id where od.delivery_date = \""+chosenDate+"\" and o.status = 1 and o.start_date <=\""+chosenDate+"\"	and o.end_date >= \""+chosenDate+"\";") 
  checkErr(err)
  for rows.Next(){
    var order Order
    order.ScanFromSQL(rows)
    orderList = append(orderList,order)
  }
  return orderList
}




