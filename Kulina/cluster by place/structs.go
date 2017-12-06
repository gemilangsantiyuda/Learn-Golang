package main
import (
   _ "github.com/go-sql-driver/mysql"
   "database/sql"
   "sync"
   "text/template"
)

type CoordinateSQLContainer struct {
  Latitude,Longitude sql.NullFloat64
}

type Coordinate struct {
  Latitude,Longitude float64
}

type KitchenCapacitySQLContainer struct {
  Min,Optimum,Max sql.NullInt64
}

type KitchenCapacity struct {
  Min,Optimum,Max int
}

type KitchenSQLContainer struct {
  Id sql.NullInt64
  Name sql.NullString
  Capacity KitchenCapacitySQLContainer
  Coord CoordinateSQLContainer
}

type OrderSQLContainer struct{
  Id sql.NullInt64  
  Qty sql.NullInt64  
  Coord CoordinateSQLContainer   
  PlaceId sql.NullString
}

type PlaceSQLContainer struct {
  Id sql.NullString
  Qty sql.NullInt64
}

type Distance struct{
  Index int 
  Distance float64
}

type Kitchen struct {
  Id int
  Index int
  Name string
  Capacity KitchenCapacity
  Coord Coordinate
  OrderDistanceList []Distance
  PlaceDistanceList []Distance
  DistinctOrderCount int
  OrderQty int
  IndexOnDistanceMatrix int
}

type Order struct{
  Id int
  Qty int
  Coord Coordinate
  KitchenDistanceList []Distance 
  OrderDistanceList []Distance
  ServingKitchen *Kitchen
  PlaceId string
}

type Place struct{
  Id string
  Qty int
  Coord Coordinate
  OrderList []*Order
  KitchenDistanceList []Distance
  PlaceDistanceList []Distance
  ServingKitchen *Kitchen
}

type templateHandler struct {
  once sync.Once
  filename string
  templ *template.Template
}

type DataForTemplate struct {
  KitchenList []Kitchen
  OrderList []Order
  GoogleAPIKey string
} 

type Path struct {
  Id int
  IndexList []int
  IdList []string
  ServingKitchen *Kitchen
  Length float64
  OrderQty int
}

type GoogleAPIDistanceResponse struct {
  Destination_addresses []string
  Origin_addresses []string
  Rows []struct{
    Elements []struct {
      Distance struct {
        Text string
        Value int
      }
      Duration struct {
        Text string
        Value int
      }
      Status string
    }    
  }
  Status string
}
