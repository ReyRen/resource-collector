package main

type recvMsg struct {
	NodeName string `json:"nodeName"`
}

type sendMsg struct {
	Utilize  string `json:"utilize"`
	MemUsed  string `json:"memUsed"`
	MemFreed string `json:"memFreed"`
	Occupied string `json:"occupied"`
}

type socketRecvMsg struct {
	Nodelist string `json:"nodelist"`
}

type socketSendMsg struct {
	NodeName string `json:"nodeName"`
	Occupied string `json:"occupied"`
}
