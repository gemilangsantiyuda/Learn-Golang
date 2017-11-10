package main
import (
  "fmt"
  "io/ioutil"
  "encoding/csv"
  "strings"
  "strconv"
)
  
func main(){
  //read inputFile
  inputFile, err := ioutil.ReadFile("inputScheduling.csv")
  if err!= nil{
    fmt.Println(err)
    return
  }
  
  //turn inputFile into processList for inputFile that has currentTime
  csvReader:= csv.NewReader(strings.NewReader(string(inputFile)))    
  records,err := csvReader.ReadAll()
  
  var processList []Process  
  currentTime,_ := strconv.Atoi(records[0][1]) 
  for i:=2; i<len(records); i++{
    arrivalTime,_ := strconv.Atoi(records[i][1])
    burstTime,_ := strconv.Atoi(records[i][2])
    processList = append(processList, Process{records[i][0],arrivalTime,burstTime})  
  }
  
  fmt.Println(processList)
  //get scheduledProcessList  
  //scheduling parameter scheduling(processList that wanna be scheduled ; currentTime ; number of chromosoms
  //we want in a generation; number of parents in a crossover ; PC probability of crossover ; PM probability of mutation
  scheduledProcessList := scheduling(processList,currentTime,10,10,0.6,0.1)
  
  fmt.Println(scheduledProcessList)  
}
