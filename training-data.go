package main

type recvMsg struct {
	Type         int              `json:"type"`
	OccupiedList *[]OccupiedLists `json:"occupiedList"`
	NodeName     string           `json:"nodeName"`
	PodName      string           `json:"podName"`
}

type sendMsg struct {
	GpuLabel    string `json:"gpuLabel"`
	Utilize     string `json:"utilize"`
	MemUsed     string `json:"memUsed"`
	MemFreed    string `json:"memFreed"`
	Occupied    string `json:"occupied"`
	Temperature string `json:"temp"`
	Uid         string `json:"uid"`
	Tid         string `json:"tid"`
}

type OccupiedLists struct {
	Uid     int    `json:"uid"`
	Tid     int    `json:"tid"`
	PodName string `json:"podName"`
}

type socketRecvMsg struct {
	Nodelist string `json:"nodelist"`
}

type socketSendMsg struct {
	NodeName string `json:"nodeName"`
	Occupied string `json:"occupied"`
}
