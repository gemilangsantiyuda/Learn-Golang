package main
import (
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
  "path/filepath"
  "text/template"
  "net/http"
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
