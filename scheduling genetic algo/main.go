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
  inputFile, err := ioutil.ReadFile("inputSchedulingx.csv")
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
  nChromosom := []int{10,50,100}
  nParent := []int{10,20,40}//nParent*nChromosom = jumlah parent
  PC := []float64{0.6,0.7,0.8}
  PM := []float64{0.1,0.15,0.2}
  PT := []float64{0.5,0.6,0.7}  
  //get scheduledProcessList  
  //scheduling parameter scheduling(processList that wanna be scheduled ; currentTime ; number of chromosoms
  //we want in a generation; number of parents in a crossover ; PC probability of crossover ; PM probability of mutation, PT percentage for stopping criteria
  var scheduledProcessList []Process
  for i:= range(nChromosom){
    for j:= range(nParent){
      for k:= range(PC){
        for m:= range(PM){
          for n:= range(PT){
            scheduledProcessList = scheduling(processList,currentTime,nChromosom[i],nParent[j],PC[k],PM[m],PT[n])
        }
       } 
      }  
   } 
  }
  
   
  fmt.Println(scheduledProcessList)  
}
