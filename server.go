package main

import (
		"fmt"
		"net"
		"net/rpc"
		"net/rpc/jsonrpc"
		"os"
		"math/rand"
		"strings"
		"strconv"
		"io/ioutil"
		"net/http"
		"encoding/json"	
)

var id int
var hmap map[int]string
var umap map[int]float32

type ClientRequest struct{
	RequestInfo string
	Budget float32
}

type ServerReply struct {
    ReplyInfo string
    Balance float32
    ID int
}

type CheckRequest struct{
	ID int
}

type CheckReply struct{
	CheckResult string
	CurrentMarketValue float32
	UninvestedAmount float32
}

type StockAccount struct{}

func (t *StockAccount) BuyStock(args *ClientRequest, reply *ServerReply) error{
	testString := args.RequestInfo
	budget := args.Budget
	replyString := ""
	var residualValue1 float32

	if strings.Contains(testString, ","){
		spString := strings.Split(testString, ",")
		for i := 0; i < len(spString); i++{									//我这里要写一个test是不是100
			innerString := strings.Split(spString[i], ":")
			stockPrice2 := getQuote(innerString[0])
			stockPrice_parse, _ := strconv.ParseFloat(stockPrice2, 64)
			quote1 := (float32(stockPrice_parse))

			innerString1 := strings.Split(innerString[1], "%")
			goParse := innerString1[0]
			percent1, _ := strconv.ParseFloat(goParse, 64)
			percentage1 := (float32(percent1))
			budgetStock := budget * (percentage1 /100)
			if i == 0{
				residualValue1 = budget - budgetStock
			}else{
				residualValue1 = residualValue1 - budgetStock
				if residualValue1 < 0{
					os.Exit(1)
				}
			}

			var result1 float32
			result1 = budgetStock / quote1
			var share1 int = int(result1)
			residualValue1 = budgetStock - (quote1 * float32(share1)) + residualValue1

			fmt.Println(share1)
			fmt.Println(residualValue1)

			stockShare1 := strconv.Itoa(share1)

			reply.ID = getID()

			if i == 0{
				replyString = innerString[0]+":"+stockShare1+":$"+stockPrice2
			}else{
				replyString = replyString+","+innerString[0]+":"+stockShare1+":$"+stockPrice2
			}
			
			reply.Balance = residualValue1
		}
		hmap[reply.ID] = replyString
		umap[reply.ID] = reply.Balance
	}else{
		spString := strings.Split(testString, ":")
		stockPrice := getQuote(spString[0])
		stockPrice1, _ := strconv.ParseFloat(stockPrice, 64)
		quote2 := (float32(stockPrice1))
		
		spString1 := strings.Split(spString[1],"%")
		readyParse := spString1[0]
		percent, _ := strconv.ParseFloat(readyParse, 64)
		percentage := (float32(percent))
		
		budget_stock := budget * (percentage / 100)
		residualValue := budget - budget_stock
		var result float32 
		result = budget_stock / quote2
		var share int = int(result) 
		residualValue = budget_stock - (quote2 * float32(share)) + residualValue

		fmt.Println(share)
		fmt.Println(residualValue)

		stockShare := strconv.Itoa(share)
		replyString = spString[0]+":"+stockShare+":"+"$"+stockPrice
		fmt.Println(replyString)

		reply.Balance = residualValue
		reply.ID = getID()

		hmap[reply.ID] = replyString
		umap[reply.ID] = reply.Balance

	}

	reply.ReplyInfo = replyString
	
	return nil
}


func getQuote(stockSymbol string) string{
	baseUrlLeft := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quote%20where%20symbol%20%3D%20%22"
	baseUrlRight := "%22&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
	url := baseUrlLeft+stockSymbol+baseUrlRight

	r, _ := http.Get(url)

	var data map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	byt := []byte(body)
	if err := json.Unmarshal(byt, &data); err!=nil{
		panic(err)
	}

	resource1 := data["query"].(map[string]interface{})
	resource2 := resource1["results"].(map[string]interface{})
	resource3 := resource2["quote"].(map[string]interface{})
	resource4 := resource3["LastTradePriceOnly"].(string)

	return resource4
}

func getID() int{
	if id == 0{
		for id == 0{
			id = rand.Intn(10000)
		}
	}else{
		id = id + 1
	}
	return id
}

func (t *StockAccount) CheckAccount(args *CheckRequest, reply *CheckReply) error{
	getContent := ""
	replyString1 := ""
	var currentMarket float32
	var total float32
	currentMarket = 0
	total = 0
	checkID := args.ID
	for key, value := range hmap{
		if key == checkID{
			getContent = value
			if strings.Contains(getContent, ","){
				spString2 := strings.Split(getContent, ",")
				for i := 0; i < len(spString2); i++{
					innerString1 := strings.Split(spString2[i], ":")
					stockSymbol1 := innerString1[0]
					stockShare2 := innerString1[1]
					intShare, _ := strconv.Atoi(stockShare2)
					purchasePrice := innerString1[2]
					purchasePrice1 := strings.Split(purchasePrice, "$")
					fpurchasePrice, _ := strconv.ParseFloat(purchasePrice1[1], 64)
				
					currentPPrice := (float32(fpurchasePrice))
					

					// get Current Price for the stock
					currentPrice := getQuote(stockSymbol1)
					fcurrentPrice, _ := strconv.ParseFloat(currentPrice, 64)
					currentQuote := (float32(fcurrentPrice))

					currentMarket = currentQuote * (float32(intShare))

					total = total + currentMarket
					
					cal := currentQuote - currentPPrice

					if i == 0{
						if cal > 0{
							replyString1 = stockSymbol1+":"+stockShare2+":+$"+currentPrice
						}else if cal < 0{
							replyString1 = stockSymbol1+":"+stockShare2+":-$"+currentPrice
						}else if cal == 0{
							replyString1 = stockSymbol1+":"+stockShare2+":$"+currentPrice
						}
					}else{
						if cal > 0{
							replyString1 = replyString1 + "," + stockSymbol1+":"+stockShare2+":+$"+currentPrice
						}else if cal < 0{
							replyString1 = replyString1 + "," + stockSymbol1+":"+stockShare2+":-$"+currentPrice
						}else if cal == 0{
							replyString1 = replyString1 + "," + stockSymbol1+":"+stockShare2+":$"+currentPrice
						}
					}
				}
				reply.CheckResult = replyString1
				reply.CurrentMarketValue = total
			}else{
				spString3 := strings.Split(getContent,":")
				stockSymbol1 := spString3[0]
				stockShare2 := spString3[1]
				intShare, _ := strconv.Atoi(stockShare2)
				purchasePrice := spString3[2]
				purchasePrice1 := strings.Split(purchasePrice, "$")
				fpurchasePrice, _ := strconv.ParseFloat(purchasePrice1[1], 64)
				currentPPrice := (float32(fpurchasePrice))

				currentPrice := getQuote(stockSymbol1)
				fcurrentPrice, _ := strconv.ParseFloat(currentPrice, 64)
				currentQuote := (float32(fcurrentPrice))

				cal := currentQuote - currentPPrice
				
				currentMarket = currentQuote * (float32(intShare))
				total = total + currentMarket

				if cal > 0{
					replyString1 = stockSymbol1+":"+stockShare2+":+$"+currentPrice
				}else if cal < 0{
					replyString1 = stockSymbol1+":"+stockShare2+":-$"+currentPrice
				}else{
					replyString1 = stockSymbol1+":"+stockShare2+":$"+currentPrice
				}
				reply.CheckResult = replyString1
				reply.CurrentMarketValue = total
			}
		}else{
			os.Exit(1)
		}
	}

	for key1, value1 := range umap{
		if(key1 == checkID){
			reply.UninvestedAmount = value1
		}else{
			os.Exit(1)
		}
	}
	return nil
}

func main() {
	id = 0
	hmap = make(map[int]string)
	umap = make(map[int]float32)
	stockAccount := *(new(StockAccount))
    rpc.Register(&stockAccount)

    tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
    checkError(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        jsonrpc.ServeConn(conn)
    }
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
