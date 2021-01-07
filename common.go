package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"os/exec"
	"strings"
	//	"strconv"
	"time"
)

const (
	// ip of mine
	websocketServer = "172.18.29.80:9400"

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Hour
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	// log
	Trace   *log.Logger // 记录所有日志
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 非常严重的问题

	RC_ENGINE_SERVER = "172.18.29.80"

	NODE1_V100_PORT   = "9401"
	NODE2_V100_PORT   = "9402"
	NODE3_V100_PORT   = "9403"
	NODE4_A100_PORT   = "9404"
	NODE5_A100_PORT   = "9405"
	NODE6_A100_PORT   = "9406"
	NODE7_2080TI_PORT = "9407"
	NODE8_2080TI_PORT = "9408"
)

var (
	GLOBALCHAN = make(chan int)
)

func jsonHandler(data []byte, v interface{}) {
	errJson := json.Unmarshal(data, v)
	if errJson != nil {
		Error.Printf("json err: %s\n", errJson)
	}
}

func getGpuRsInfo(c *Client) {
	switch c.rm.NodeName {
	case "node1":
		//NODE1_V100_PORT
		utilize, memUsed, memFreed, occupied, tmp := curl_metrics(RC_ENGINE_SERVER, NODE1_V100_PORT)
		c.sm.Utilize = utilize
		c.sm.MemUsed = memUsed
		c.sm.MemFreed = memFreed
		c.sm.Occupied = occupied
		c.sm.Temperature = tmp
	case "node2":
		//NODE2_V100_PORT
		utilize, memUsed, memFreed, occupied, tmp := curl_metrics(RC_ENGINE_SERVER, NODE2_V100_PORT)
		c.sm.Utilize = utilize
		c.sm.MemUsed = memUsed
		c.sm.MemFreed = memFreed
		c.sm.Occupied = occupied
		c.sm.Temperature = tmp
	case "node3":
		//NODE3_V100_PORT
		utilize, memUsed, memFreed, occupied, tmp := curl_metrics(RC_ENGINE_SERVER, NODE3_V100_PORT)
		c.sm.Utilize = utilize
		c.sm.MemUsed = memUsed
		c.sm.MemFreed = memFreed
		c.sm.Occupied = occupied
		c.sm.Temperature = tmp
	case "node4":
		//NODE4_A100_PORT
		utilize, memUsed, memFreed, occupied, tmp := curl_metrics(RC_ENGINE_SERVER, NODE4_A100_PORT)
		c.sm.Utilize = utilize
		c.sm.MemUsed = memUsed
		c.sm.MemFreed = memFreed
		c.sm.Occupied = occupied
		c.sm.Temperature = tmp
	case "node5":
		//NODE5_A100_PORT
		utilize, memUsed, memFreed, occupied, tmp := curl_metrics(RC_ENGINE_SERVER, NODE5_A100_PORT)
		c.sm.Utilize = utilize
		c.sm.MemUsed = memUsed
		c.sm.MemFreed = memFreed
		c.sm.Occupied = occupied
		c.sm.Temperature = tmp
	case "node6":
		//NODE6_A100_PORT
		utilize, memUsed, memFreed, occupied, tmp := curl_metrics(RC_ENGINE_SERVER, NODE6_A100_PORT)
		c.sm.Utilize = utilize
		c.sm.MemUsed = memUsed
		c.sm.MemFreed = memFreed
		c.sm.Occupied = occupied
		c.sm.Temperature = tmp
	case "node7":
		//NODE7_2080TI_PORT
		utilize, memUsed, memFreed, occupied, tmp := curl_metrics(RC_ENGINE_SERVER, NODE7_2080TI_PORT)
		c.sm.Utilize = utilize
		c.sm.MemUsed = memUsed
		c.sm.MemFreed = memFreed
		c.sm.Occupied = occupied
		c.sm.Temperature = tmp
	case "node8":
		//NODE8_2080TI_PORT
		utilize, memUsed, memFreed, occupied, tmp := curl_metrics(RC_ENGINE_SERVER, NODE8_2080TI_PORT)
		c.sm.Utilize = utilize
		c.sm.MemUsed = memUsed
		c.sm.MemFreed = memFreed
		c.sm.Occupied = occupied
		c.sm.Temperature = tmp
	}
}

func getGpuOccuppiedInfo(nodeName string, sendSocketMsg *socketSendMsg) {
	switch nodeName {
	case "node1":
		//NODE1_V100_PORT
		_, _, _, occupied, _ := curl_metrics(RC_ENGINE_SERVER, NODE1_V100_PORT)
		sendSocketMsg.NodeName = "node1"
		sendSocketMsg.Occupied = occupied
	case "node2":
		//NODE2_V100_PORT
		_, _, _, occupied, _ := curl_metrics(RC_ENGINE_SERVER, NODE2_V100_PORT)
		sendSocketMsg.NodeName = "node2"
		sendSocketMsg.Occupied = occupied
	case "node3":
		//NODE3_V100_PORT
		_, _, _, occupied, _ := curl_metrics(RC_ENGINE_SERVER, NODE3_V100_PORT)
		sendSocketMsg.NodeName = "node3"
		sendSocketMsg.Occupied = occupied
	case "node4":
		//NODE4_A100_PORT
		_, _, _, occupied, _ := curl_metrics(RC_ENGINE_SERVER, NODE4_A100_PORT)
		sendSocketMsg.NodeName = "node4"
		sendSocketMsg.Occupied = occupied
	case "node5":
		//NODE5_A100_PORT
		_, _, _, occupied, _ := curl_metrics(RC_ENGINE_SERVER, NODE5_A100_PORT)
		sendSocketMsg.NodeName = "node5"
		sendSocketMsg.Occupied = occupied
	case "node6":
		//NODE6_A100_PORT
		_, _, _, occupied, _ := curl_metrics(RC_ENGINE_SERVER, NODE6_A100_PORT)
		sendSocketMsg.NodeName = "node6"
		sendSocketMsg.Occupied = occupied
	case "node7":
		//NODE7_2080TI_PORT
		_, _, _, occupied, _ := curl_metrics(RC_ENGINE_SERVER, NODE7_2080TI_PORT)
		sendSocketMsg.NodeName = "node7"
		sendSocketMsg.Occupied = occupied
	case "node8":
		//NODE8_2080TI_PORT
		_, _, _, occupied, _ := curl_metrics(RC_ENGINE_SERVER, NODE8_2080TI_PORT)
		sendSocketMsg.NodeName = "node8"
		sendSocketMsg.Occupied = occupied
	}
}

func curl_metrics(ips string, port string) (string, string, string, string, string) {

	var utilize string
	var memUsed string
	var memFreed string
	var occupied string
	var temp string

	base_cmd_string := "curl http://" + ips + ":" + port + "/metrics | grep gpu | grep "

	gpu_util := "DCGM_FI_DEV_GPU_UTIL"
	fb_free := "DCGM_FI_DEV_FB_FREE"
	fp_used := "DCGM_FI_DEV_FB_USED"
	temp_used := "DCGM_FI_DEV_GPU_TEMP"

	gpu_util_res, _ := exec.Command("/bin/bash", "-c", base_cmd_string+gpu_util).Output()
	fb_free_res, _ := exec.Command("/bin/bash", "-c", base_cmd_string+fb_free).Output()
	fp_used_res, _ := exec.Command("/bin/bash", "-c", base_cmd_string+fp_used).Output()
	temp_res, _ := exec.Command("/bin/bash", "-c", base_cmd_string+temp_used).Output()

	trimStringValue(string(gpu_util_res), &utilize)
	trimStringValue(string(fb_free_res), &memFreed)
	trimStringValue(string(fp_used_res), &memUsed)
	trimStringOcp(string(gpu_util_res), &occupied)
	trimStringValue(string(temp_res), &temp)

	//Trace.Printf("gpu_util_res: %s\n", gpu_util_res)
	//Trace.Printf("fb_free_res: %s\n", fb_free_res)
	//Trace.Printf("fp_used_res: %s\n", fp_used_res)

	return utilize, memUsed, memFreed, occupied, temp
}

func trimStringValue(src string, dst *string) {
	src_slice := strings.Split(src, "\n")
	src_slice = src_slice[:len(src_slice)-1]
	for _, src_single := range src_slice {
		*dst += strings.Split(src_single, " ")[len(strings.Split(src_single, " "))-1]
		*dst += ","
	}
}
func trimStringOcp(src string, dst *string) {
	src_slice := strings.Split(src, "\n")
	src_slice = src_slice[:len(src_slice)-1]
	for _, src_single := range src_slice {
		if strings.Contains(src_single, "pod=\"\"") {
			*dst += "0"
			*dst += ","
		} else {
			if strings.Contains(src_single, "pod=") {
				*dst += "1"
				*dst += ","
			} else {
				*dst += " "
				*dst += ","
			}
		}
	}
}

func serverSocketCreate() {
	//建立socket
	netListen, err := net.Listen("tcp", "172.18.29.80:8082")
	if err != nil {
		Error.Printf("net.Listen: %s\n", err)
		return
	}
	defer netListen.Close()

	for {
		Trace.Printf("waiting for socket")
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		Trace.Printf("%s connected in\n", conn.RemoteAddr().String())
		handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	buffer := make([]byte, 4096)
	var readNum int
	for {
		tmpNum, _ := conn.Read(buffer)
		if tmpNum == 0 {
			Trace.Printf("read done!\n")
			break
		}
		readNum += tmpNum
	}
	Trace.Printf("read socket msg %s\n", string(buffer[:readNum]))
	message := bytes.TrimSpace(bytes.Replace(buffer[:readNum], newline, space, -1))
	fmt.Printf("received messages: %s\n", message)
	var recvMsg socketRecvMsg
	rms := &recvMsg
	jsonHandler(message, rms)
	respond(rms.Nodelist, conn)
}

func respond(nodeName string, conn net.Conn) {
	if nodeName == "" {
		return
	}
	var sendMsgs []socketSendMsg
	stringSlice := strings.Split(nodeName, ", ")
	for _, nodename := range stringSlice {
		var sendMsg socketSendMsg
		sms := &sendMsg
		getGpuOccuppiedInfo(nodename, sms)

		sendMsgs = append(sendMsgs, sendMsg)

	}
	smsSend, err := json.Marshal(sendMsgs)
	Trace.Printf("send msg %s\n", string(smsSend))
	if err != nil {
		Error.Printf("json.Marshal err: %s\n", err)
	}
	_, err = conn.Write(smsSend)
	_, err = conn.Write([]byte("\n"))
	if err != nil {
		Error.Printf("socket write err: %s\n", err)
	}
}
