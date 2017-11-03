// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
  "strings"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/line/line-bot-sdk-go/linebot"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"
)

var bot *linebot.Client
var pictureURLs map[string]map[string][]string


func main() {
	var err error
	buildPictureURLs()
	bot, err = linebot.New("a718d0a9fc76394fed297f266382ef6d", "ue1Sf9ZZgQ83xFsoBvmBwFCzRMZRGFY6UUGk1P0aQo4nEV9N7RfnPqncTVsbbQLwLvwt9OpSx3d2NReJTwBEoFNJ40VINnPGswVOmFiJoS6CdsrXBD4IWjDcM163TGonGByajjS3nzVWO3jm0/xq/wdB04t89/1O/w1cDnyilFU=")
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func buildPictureURLs() {
  raw,err := ioutil.ReadFile("characters.json")
  if err != nil { 
    fmt.Println(err.Error())     
  }
  err = json.Unmarshal(raw,&pictureURLs)
  if err!= nil{
    fmt.Println(err.Error())
  }  
}

func generatePictureURL(chat string) string {
	var url string
	query := strings.Fields(message.Text)
  log.Println(query)        
  if len(query)!= 2 {
    return ""
  } 
  if val,ok := pictureURLs[query[1]]; ok {
    log.Print(val)
    if val1,ok1 := pictureURLs[query[1]][query[0]]; ok1 {
      log.Print(val1)
      rand.Seed(time.Now().UTC().UnixNano())
      idx := rand.Intn(len(pictureURLs[query[1]][query[0]]))
      url = pictureURLs[query[1]][query[0]][idx]
    } else {
    	return ""
    }
  } else {
  	return ""
  }     
  return url   		
}

func callbackHandler(w http.ResponseWriter, req *http.Request) {
	events, err := bot.ParseRequest(req)
  log.Println(events)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
        url:= generatePictureURL(message.Text)    
				if url == "" {
					return
				}        	
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(url,url)).Do(); err != nil {
					log.Print(err)
				}
			case *linebot								
			}
		}
	}
}
