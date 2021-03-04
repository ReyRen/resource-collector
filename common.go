package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"os/exec"
	"strconv"
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

	//GPU UUID
	GPU1_UUID = "GPU-b68354fb-5929-8aaf-664e-7ca9fba26a86"
	GPU2_UUID = "GPU-26aada9b-5512-8b28-b3a4-db794fad2882"
	GPU3_UUID = "GPU-defdf4e1-49f5-e7ca-f907-9454a08f3b60"
	GPU4_UUID = "GPU-0dd65ed1-c3f4-a598-e3a9-622de3a57944"
	GPU5_UUID = "GPU-56bc8b4f-e8b6-7610-0db6-84adf19b7da1"
	GPU6_UUID = "GPU-7e8b3724-616a-5152-cad9-114c119845bd"
	GPU7_UUID = "GPU-ef57863b-dac3-605d-07cc-2418f3cdd60f"
	GPU8_UUID = "GPU-7fb7e9de-0a2c-117b-fbd9-fb864c989032"
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

	RC_ENGINE_SERVER = "10.193.215.64:9400"
)

func jsonHandler(data []byte, v interface{}) {
	errJson := json.Unmarshal(data, v)
	if errJson != nil {
		Error.Printf("json err: %s\n", errJson)
	}
}

func getGpuRsInfo(c *Client) {
	var flag int
	var matchedMetrics string

	// used to check node
	var UUIDCHECK string
	switch c.rm.NodeName {
	case "node1":
		UUIDCHECK = GPU1_UUID
	case "node2":
		UUIDCHECK = GPU2_UUID
	case "node3":
		UUIDCHECK = GPU3_UUID
	case "node4":
		UUIDCHECK = GPU4_UUID
	case "node5":
		UUIDCHECK = GPU5_UUID
	case "node6":
		UUIDCHECK = GPU6_UUID
	case "node7":
		UUIDCHECK = GPU7_UUID
	case "node8":
		UUIDCHECK = GPU8_UUID
	}

	for flag = 0; flag < 10; flag++ { // retries 10 times
		metrics := getLoadbalanceMetrics(RC_ENGINE_SERVER)
		if strings.Contains(metrics, UUIDCHECK) {
			// match
			matchedMetrics = metrics
			break
		}
	}
	if flag == 10 {
		Trace.Println("Nothing get...from: ", c.rm.NodeName)
		c.sm.GpuLabel = ""
		c.sm.Utilize = ""
		c.sm.MemUsed = ""
		c.sm.MemFreed = ""
		c.sm.Occupied = ""
		c.sm.Temperature = ""
		return
	}

	handle_metrics(c, matchedMetrics)
}

func getLoadbalanceMetrics(ips string) string {
	var result string

	base_cmd_string := "curl http://" + ips + "/metrics | grep gpu"
	temp_res, _ := exec.Command("/bin/bash", "-c", base_cmd_string).Output()
	result = string(temp_res)

	return result
}

func getGpuOccuppiedInfo(nodeName string, sendSocketMsg *socketSendMsg) {

	var flag int
	var matchedMetrics string
	var UUIDCHECK string
	var nodeNameBack string

	switch nodeName {
	case "node1":
		UUIDCHECK = GPU1_UUID
		nodeNameBack = "node1"
	case "node2":
		UUIDCHECK = GPU2_UUID
		nodeNameBack = "node2"
	case "node3":
		UUIDCHECK = GPU3_UUID
		nodeNameBack = "node3"
	case "node4":
		UUIDCHECK = GPU4_UUID
		nodeNameBack = "node4"
	case "node5":
		UUIDCHECK = GPU5_UUID
		nodeNameBack = "node5"
	case "node6":
		UUIDCHECK = GPU6_UUID
		nodeNameBack = "node6"
	case "node7":
		UUIDCHECK = GPU7_UUID
		nodeNameBack = "node7"
	case "node8":
		UUIDCHECK = GPU8_UUID
		nodeNameBack = "node8"
	}

	for flag = 0; flag < 10; flag++ {
		metrics := getLoadbalanceMetrics(RC_ENGINE_SERVER)
		if strings.Contains(metrics, UUIDCHECK) {
			// match
			matchedMetrics = metrics
			break
		}
	}
	if flag == 10 {
		Trace.Println("Nothing get...when send to socket client...from: ", nodeNameBack)
		sendSocketMsg.Occupied = ""
		return
	}

	handle_metrics_to_socket(sendSocketMsg, matchedMetrics)
	sendSocketMsg.NodeName = nodeNameBack
}

func handle_metrics(c *Client, metrics string) {
	var gpuLabel string
	var utilize string
	var memUsed string
	var memFreed string
	var occupied string
	var temp string
	var uids string
	var tids string

	var flag int

	gpu_util := "DCGM_FI_DEV_GPU_UTIL"
	fb_free := "DCGM_FI_DEV_FB_FREE"
	fp_used := "DCGM_FI_DEV_FB_USED"
	temp_used := "DCGM_FI_DEV_GPU_TEMP"

	src_slice := strings.Split(metrics, "\n")
	for _, src_single := range src_slice {
		if c.rm.Type == 1 {
			if strings.Contains(src_single, gpu_util) {
				utilize += strings.Split(src_single, " ")[len(strings.Split(src_single, " "))-1]
				utilize += ","

				if strings.Contains(src_single, "pod=\"\"") {
					occupied += "0"
					occupied += ","
				} else {
					if strings.Contains(src_single, "pod=") {
						occupied += "1"
						occupied += ","
					} else {
						occupied += " "
						occupied += ","
					}
				}

				gpuLabel += trimQuotes(strings.Split(strings.Split(src_single, ",")[0], "=")[1])
				gpuLabel += ","

				if len(*(c.rm.OccupiedList)) > 0 {
					// Administrator vision
					flag = 0
					for _, v := range *(c.rm.OccupiedList) {
						if strings.Contains(src_single, v.PodName) {
							Trace.Println("administrator v.PodName = ", v.PodName)
							uids += strconv.Itoa(v.Uid)
							tids += strconv.Itoa(v.Tid)
							uids += ","
							tids += ","
							break
						}
						flag++
					}
					if flag == len(*(c.rm.OccupiedList)) {
						/*uids += " "
						tids += " "*/
						uids += ","
						tids += ","
					}
				}
			} else if strings.Contains(src_single, fb_free) {
				memFreed += strings.Split(src_single, " ")[len(strings.Split(src_single, " "))-1]
				memFreed += ","
			} else if strings.Contains(src_single, fp_used) {
				memUsed += strings.Split(src_single, " ")[len(strings.Split(src_single, " "))-1]
				memUsed += ","
			} else if strings.Contains(src_single, temp_used) {
				temp += strings.Split(src_single, " ")[len(strings.Split(src_single, " "))-1]
				temp += ","
			}
		} else if c.rm.Type == 2 {
			if strings.Contains(src_single, gpu_util) && strings.Contains(src_single, c.rm.PodName) {
				utilize += strings.Split(src_single, " ")[len(strings.Split(src_single, " "))-1]
				utilize += ","

				if strings.Contains(src_single, "pod=\"\"") {
					occupied += "0"
					occupied += ","
				} else {
					if strings.Contains(src_single, "pod=") {
						occupied += "1"
						occupied += ","
					} else {
						occupied += " "
						occupied += ","
					}
				}

				gpuLabel += trimQuotes(strings.Split(strings.Split(src_single, ",")[0], "=")[1])
				gpuLabel += ","
			} else if strings.Contains(src_single, fb_free) && strings.Contains(src_single, c.rm.PodName) {
				memFreed += strings.Split(src_single, " ")[len(strings.Split(src_single, " "))-1]
				memFreed += ","
			} else if strings.Contains(src_single, fp_used) && strings.Contains(src_single, c.rm.PodName) {
				memUsed += strings.Split(src_single, " ")[len(strings.Split(src_single, " "))-1]
				memUsed += ","
			} else if strings.Contains(src_single, temp_used) && strings.Contains(src_single, c.rm.PodName) {
				temp += strings.Split(src_single, " ")[len(strings.Split(src_single, " "))-1]
				temp += ","
			}
		}
	}
	c.sm.GpuLabel = gpuLabel
	c.sm.Utilize = utilize
	c.sm.MemUsed = memUsed
	c.sm.MemFreed = memFreed
	c.sm.Occupied = occupied
	c.sm.Temperature = temp
	c.sm.Uid = uids
	c.sm.Tid = tids
}

func handle_metrics_to_socket(sendSocketMsg *socketSendMsg, metrics string) {
	var occupied string

	gpu_util := "DCGM_FI_DEV_GPU_UTIL"

	src_slice := strings.Split(metrics, "\n")
	for _, src_single := range src_slice {
		if strings.Contains(src_single, gpu_util) {
			if strings.Contains(src_single, "pod=\"\"") {
				occupied += "0"
				occupied += ","
			} else {
				if strings.Contains(src_single, "pod=") {
					occupied += "1"
					occupied += ","
				} else {
					occupied += " "
					occupied += ","
				}
			}
		}
	}
	sendSocketMsg.Occupied = occupied
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
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
