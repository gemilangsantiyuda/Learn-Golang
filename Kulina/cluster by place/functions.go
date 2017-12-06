package main
import (
  _ "github.com/go-sql-driver/mysql"
  "path/filepath"
  "text/template"
  "net/http"
  "log"
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
  data:= DataForTemplate{KitchenList,OrderList,GOOGLE_API_KEY} 
  t.templ.Execute(w,data)
}

func RunKitchensAndOrdersView(){
  http.Handle("/", &templateHandler{filename: "GoogleMap.html"})
  // start the web server
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }   
}

func OrdersCanSwapKitchen(orderA,orderB *Order) bool {
  kitchenA := KitchenList[orderA.KitchenDistanceList[0].Index]
  kitchenB := KitchenList[orderB.KitchenDistanceList[0].Index]  
  kitchenANewQty := kitchenA.OrderQty - orderA.Qty + orderB.Qty
  kitchenBNewQty := kitchenB.OrderQty - orderB.Qty + orderA.Qty 
  kitchenAConstraint := kitchenANewQty>=kitchenA.Capacity.Min && kitchenANewQty<=kitchenA.Capacity.Max
  kitchenBConstraint := kitchenBNewQty>=kitchenB.Capacity.Min && kitchenBNewQty<=kitchenB.Capacity.Max
  return kitchenAConstraint && kitchenBConstraint
}

func PlacesCanSwapKitchen(placeA,placeB *Place) bool {
  kitchenA := placeA.ServingKitchen
  kitchenB := placeB.ServingKitchen
  kitchenANewQty := kitchenA.OrderQty - placeA.Qty + placeB.Qty
  kitchenBNewQty := kitchenB.OrderQty - placeB.Qty + placeA.Qty 
  kitchenAConstraint := kitchenANewQty>=kitchenA.Capacity.Min && kitchenANewQty<=kitchenA.Capacity.Max
  kitchenBConstraint := kitchenBNewQty>=kitchenB.Capacity.Min && kitchenBNewQty<=kitchenB.Capacity.Max
  return kitchenAConstraint && kitchenBConstraint
}


func CheckKitchenCapacityToOrderQty() bool {
  totalOrderQty := 0
  totalCapacity := 0
  for i:= range OrderList {
    foundKitchen := false
    for j:= range KitchenList {
      if OrderList[i].Qty <= KitchenList[j].Capacity.Max{
        foundKitchen=true
        break
      }
    }
    if !foundKitchen{
      return false
    }
  } 
  for i:=range OrderList {
    totalOrderQty += OrderList[i].Qty
  }
  for i:=range KitchenList {
    totalCapacity += KitchenList[i].Capacity.Max
  }
  return totalOrderQty<=totalCapacity
}
