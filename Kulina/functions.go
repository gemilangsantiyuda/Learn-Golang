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
