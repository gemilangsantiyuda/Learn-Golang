package main
import (
  "math/rand"
  "math"
  "fmt"
  "sort"
  "sync"
  "time"
)
type Process struct{
  ID string;
  ArrivalTime int;
  BurstTime int;
}

type Individu struct{
  Chromosom []int;
  Fitness float64;
}

//for concurrency purpose :
var wg sync.WaitGroup
var m sync.Mutex

func (ind *Individu) calculateFitness(processList []Process, currentTime int) {
  //calculate Average Turn Around Time = average of sum of endtime-arrivaltime for eachprocess
  defer wg.Done()
  ATAT := float64(0)
	for i:= range(ind.Chromosom){
	  currentTime+=processList[ind.Chromosom[i]].BurstTime
	  arrivalTime:=processList[ind.Chromosom[i]].ArrivalTime
		ATAT += float64(currentTime - arrivalTime)
	}
	ATAT/= float64(len(ind.Chromosom))
	ind.Fitness = 1/ATAT
}

func initChromosom(processList []Process,currentTime int,nChromosom,length int) []Individu{
  var initialChromosom []Individu  
  //usedRandom for random permutation control
  var usedRandom []int
  for i:=0;i<length;i++{
    usedRandom = append(usedRandom,-1)
  }
  
  //initiating random individus
  for i:=0;i<nChromosom;i++{    
    //randomize chromosom
    var temp []int    
    for j:=0;j<length;j++{
      randomIdx := rand.Intn(length);
      for usedRandom[randomIdx]==i {
        randomIdx= rand.Intn(length)
      } 
      usedRandom[randomIdx]=i
      temp = append(temp,randomIdx)
    }    
    //get the randomize Chromosom fitness value
    initialChromosom = append(initialChromosom,Individu{temp,0})
    wg.Add(1)
    go initialChromosom[i].calculateFitness(processList,currentTime)
  }  
  wg.Wait()
  return initialChromosom
}


func selectParents(chromosomList []Individu, nParent int, PC float64 ) []Individu{
  //roulettewheel the parents
  var candidateParentList []Individu
    
  var portion []float64 //percentage of fitness for roulette wheel
  sort.Slice(chromosomList,func(lhs,rhs int) bool {
    return chromosomList[lhs].Fitness>chromosomList[rhs].Fitness 
  })
  
  totalFitness:=float64(0)
  for i:= range(chromosomList){
    totalFitness+=chromosomList[i].Fitness
  }
  
  //generating portion (width of each individu on the wheel) for roulette wheel
  for i:= range(chromosomList){
    portion = append(portion,chromosomList[i].Fitness/totalFitness)
    if i > 0{
      portion[i]+=portion[i-1]  
    }
  }
  
  //generating parent
  for i:=0;i<nParent;i++{
    //generating 2 random numbers, rw for roullette wheel ... pc for comparison with PC
    rw := float64(rand.Intn(100001))/100000
    pc := float64(rand.Intn(100001))/100000
    if pc<PC{
      //check where rw falls in roulette wheel
      for j:= range(portion){
        if rw<=portion[j]{
          candidateParentList = append(candidateParentList,chromosomList[j])
          break
        }
      }
    }    
  }  
      
  return candidateParentList
}

func orderCrossover(parentA,parentB Individu)(Individu,Individu){
  length:= len(parentA.Chromosom)
  //dummy copy parent to offspring
  var offspringA,offspringB []int
  for i:=0;i<length;i++{
    offspringA = append(offspringA,0)
    offspringB = append(offspringB,0)    
  }  
  
  idxA:= rand.Intn(len(parentA.Chromosom))
  idxB:= rand.Intn(len(parentA.Chromosom))  
  
  for idxA==idxB{
    idxB= rand.Intn(len(parentA.Chromosom))
  }
  if idxA>idxB {
    tmp := idxA
    idxA = idxB
    idxB = tmp
  }
  
  var existA,existB []bool
  for i:=0;i<length;i++{
    existA,existB = append(existA,false),append(existB,false)
  }

  for i:=idxA;i<=idxB;i++{
    offspringA[i] = parentA.Chromosom[i]
    existA[offspringA[i]] = true
    offspringB[i] = parentB.Chromosom[i]        
    existB[offspringB[i]] = true  
  }
  //time.Sleep(3*time.Microsecond)
  //complete the offspringA
  currentIdx:=0
  for i:=0;i<length;i++{
    if currentIdx==idxA{
      currentIdx = idxB+1
    }
    if (!existA[parentB.Chromosom[i]]) {
      offspringA[currentIdx] = parentB.Chromosom[i]
      currentIdx++
    }
  }
  
  //complete the offspringB
  currentIdx=0
  for i:=0;i<length;i++{
    if currentIdx==idxA{
      currentIdx = idxB+1
    }
    if(!existB[parentA.Chromosom[i]]) {
      offspringB[currentIdx] = parentA.Chromosom[i]
      currentIdx++
    }
  }
  
  return Individu{offspringA,0},Individu{offspringB,0}
}
  
func mutate(offspringList []Individu, PM float64) []Individu{
  //determine number of mutation and where it will happen
  //fmt.Println(len(offspringList))
  length := len(offspringList[0].Chromosom)
  totalLength:= len(offspringList)*length
  nMutation := int(PM*float64(totalLength))
  var idxMutation []int
  exist:=make( map[int]bool)
  for i:=0;i<nMutation;i++{
    randomIdx:= rand.Intn(totalLength)
    for _,ok := exist[randomIdx];ok;{
      randomIdx= rand.Intn(totalLength)
      _,ok = exist[randomIdx]
    }
    exist[randomIdx]= true

    idxMutation = append(idxMutation, randomIdx)
  }

  for i:=0;i<len(idxMutation);i++{  
    n:= (idxMutation[i]/length)

    m1:= idxMutation[i]%length
    m2:= rand.Intn(length)
    for m1==m2{
      m2= rand.Intn(length)
    }        
    //fmt.Println(offspringList[n])
    //swap the chromosom of offspring[n] at gen m1 and m2
    tmp:= offspringList[n].Chromosom[m1]
    offspringList[n].Chromosom[m1] = offspringList[n].Chromosom[m2]
    offspringList[n].Chromosom[m2] = tmp
  }  

  return offspringList  
}

func crossover(processList []Process,currentTime int, parentList []Individu, PM float64) []Individu{
  var offspringList []Individu      
  var parentA,parentB Individu  
  for i:= range(parentList){
      parentA = parentList[i]
      if i==len(parentList)-1{
        parentB = parentList[0]        
  
      } else {
        parentB = parentList[i+1]
      }
      wg.Add(1)
      go func(){
        offspringA,offspringB := orderCrossover(parentA,parentB)
        m.Lock()
          offspringList = append(offspringList,offspringA)
          offspringList = append(offspringList,offspringB)               
        m.Unlock()      
        defer wg.Done()
        }()
  }
  wg.Wait()
  if len(offspringList)==0{
    return nil
  }  
  offspringList = mutate(offspringList,PM)
  return offspringList
}

func scheduling(processList []Process, currentTime int,nChromosom,nParent int, PC,PM,PT float64) []Process{
  rand.Seed( time.Now().UTC().UnixNano())


  //initialization of chromosoms (first generation)
  length:= len(processList)
  chromosomList:= initChromosom(processList,currentTime,nChromosom,length)
  
  nGeneration:=0    
  for {
    nGeneration++
  
    //parents selection using roulettewheel   
    parentList:=selectParents(chromosomList,nParent,PC)
    
    //offspring production with ordercrossover
    offspringList:= crossover(processList,currentTime,parentList,PM)
                
    //generation selection with elitism
    for i:=0;i<len(chromosomList);i++{
      offspringList = append(offspringList,chromosomList[i])
    }
    for i:=0;i<len(offspringList);i++{
      wg.Add(1)
      go offspringList[i].calculateFitness(processList,currentTime)    
    } 
    wg.Wait()
    sort.Slice(offspringList,func(lhs,rhs int) bool {
      return offspringList[lhs].Fitness>offspringList[rhs].Fitness 
    })    
    
    chromosomList = nil
    for i:=0;i<nChromosom;i++{
      chromosomList = append(chromosomList,offspringList[i])
    }
    //if difference<eps than they are considered the same
    eps:= 0.000001
    bestFitness:=chromosomList[0].Fitness
    
    
    //fmt.Println(len(chromosomList)," ",1/bestFitness)
    freqSame:=0
    for i:=0;i<len(chromosomList);i++{
      if math.Abs(chromosomList[i].Fitness-bestFitness)<eps{
        freqSame++
      }
    }
    if float64(freqSame)/float64(nChromosom) >= PT{
      break
    }
  }    
  
  //re-order the process into the calculated schedule!
  var scheduledProcess []Process
  for i:=0;i<len(chromosomList[0].Chromosom);i++{
    scheduledProcess = append(scheduledProcess,processList[chromosomList[0].Chromosom[i]])
  }
  fmt.Println(nChromosom,",",nParent,",",PC,",",PM,",",PT,",",chromosomList[0].Fitness,",",nGeneration)  
  return scheduledProcess
}
