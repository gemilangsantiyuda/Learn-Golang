package main
import (
   _ "github.com/go-sql-driver/mysql"
   "database/sql"
   "sync"
   "text/template"
)

type KitchenScan struct {
  id sql.NullInt64
  name sql.NullString
  min_capacity,max_capacity,optimum_capacity sql.NullInt64
  loc_lat,loc_lon sql.NullFloat64
}

type OrderScan struct{
  id sql.NullInt64
  qty sql.NullInt64  
  latitude,longitude sql.NullFloat64    
}

type Distance struct{
  Idx int
  Distance float64
}

type Kitchen struct {
  Id int
  Name string
  Min_capacity int
  Max_capacity int
  Optimum_capacity int
  Loc_lat float64  
  Loc_lon float64    
  Order_distance []Distance
  Order_count int
  Order_qty int
}

type Order struct{
  Id int
  Qty int
  Latitude,Longitude float64
  Kitchen_distance []Distance 
  Order_distance []Distance
}

type templateHandler struct {
  once sync.Once
  filename string
  templ *template.Template
}

type DataTemplate struct {
  KitchenList []Kitchen
  OrderList []Order
} 