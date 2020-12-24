package main

type recvMsg struct {
	NodeName	string	`json:"nodeName"`
}

type sendMsg struct {
	Utilize		string	`json:"utilize"`
	MemUsed		string	`json:"memUsed"`
	MemFreed	string	`json:"memFreed"`
	Occupied	string	`json:"occupied"`
}
