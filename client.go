package main

import (
		"fmt"
		"os"
		"log"
		"net/rpc/jsonrpc"
)

type ClientRequest struct{
	RequestInfo string
	Budget float32
}

type ServerReply struct{
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

func main() {
	var i int
	fmt.Println("Please Choose 1 for purchase stock, 2 for check balance, 3 for exit")
	fmt.Scanf("%d", &i)

	switch i{
		case 1: fmt.Println("Please Enter Stock Symbol and Percentage. (Format: GOOG:50%,YHOO:50%)")
				var input string
				fmt.Scanln(&input)

				fmt.Println("Please Enter Your Budget: ")
				var Budget float32
				fmt.Scanln(&Budget)

				client, err := jsonrpc.Dial("tcp", ":1234")
				if err != nil{
					log.Fatal("dialing: ", err)
				}
				args := ClientRequest{input, Budget}
				var reply ServerReply
				err = client.Call("StockAccount.BuyStock", &args, &reply)
				if err != nil{
					log.Fatal("Buy Stock Error: ", err)
				}
				fmt.Println("-------------------------REPLY----------------------------")
				fmt.Println(reply.ReplyInfo)
				fmt.Println(reply.Balance)
				fmt.Println(reply.ID)
			break;
		case 2: fmt.Println("Plese Enter the Transaction ID you want to check")
				var id int
				fmt.Scanln(&id)

				fmt.Println(id)

				client, err := jsonrpc.Dial("tcp", ":1234")
				if err != nil{
					log.Fatal("dialing: ", err)
				}
				args := CheckRequest{id}
				var reply CheckReply
				err = client.Call("StockAccount.CheckAccount", &args, &reply)
				if err != nil{
					log.Fatal("Buy Stock Error: ", err)
				}
				fmt.Println("-------------------------REPLY----------------------------")
				fmt.Println(reply.CheckResult)
				fmt.Println(reply.CurrentMarketValue)
				fmt.Println(reply.UninvestedAmount)

			break;
		case 3: fmt.Println("Goodbye")
				os.Exit(0)
			break;
		default:
			fmt.Println("No Such Choice, Goodbye")
			os.Exit(1)
	}

}
