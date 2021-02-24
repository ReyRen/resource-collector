package main

type recvMsg struct {
	Type     int    `json:"type"`
	NodeName string `json:"nodeName"`
	PodName  string `json:"podName"`
}

type sendMsg struct {
	GpuLabel    string `json:"gpuLabel"`
	Utilize     string `json:"utilize"`
	MemUsed     string `json:"memUsed"`
	MemFreed    string `json:"memFreed"`
	Occupied    string `json:"occupied"`
	Temperature string `json:"temp"`
}

type socketRecvMsg struct {
	Nodelist string `json:"nodelist"`
}

type socketSendMsg struct {
	NodeName string `json:"nodeName"`
	Occupied string `json:"occupied"`
}
