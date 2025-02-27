package main

import (
	"encoding/json"
	"fmt"
	rabbithttp "github.com/michaelklishin/rabbit-hole/v2"
	"strings"
)

func main() {
	list_nodes()
}

func list_nodes() {
	rmqc, err := rabbithttp.NewClient("http://localhost:15672", "admin", "hFPTzOE6F.;v3VqA")
	if err != nil {
		return
	}

	str := "https://abcd/"
	lastIndex := strings.LastIndex(str, "/")
	if lastIndex != -1 && lastIndex < len(str)-1 {
		suffix := str[lastIndex+1:]
		fmt.Println(suffix)
	} else {
		fmt.Println("")
	}

	//nodes, err := rmqc.ListNodes()
	//if err != nil {
	//	fmt.Println("Failed to list nodes", err)
	//	return
	//}

	//queues, err := rmqc.ListQueues()
	//if err != nil {
	//	fmt.Println("Failed to list queues", err)
	//	return
	//}
	//
	//for _, q := range queues {
	//	fmt.Println(q.Name)
	//	fmt.Println(q.Type)
	//	fmt.Println(q.Node)
	//}

	//overview, err := rmqc.Overview()
	//if err != nil {
	//	return
	//}

	queues, err := rmqc.ListQueues()
	if err != nil {
		fmt.Println("Failed to list queues", err)
		return
	}

	dump(queues)
}

func dump(data interface{}) {
	b, _ := json.MarshalIndent(data, "", "  ")
	fmt.Print(string(b))
}
