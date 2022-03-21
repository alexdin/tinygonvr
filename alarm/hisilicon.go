package alarm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Server struct {
	Port           string
	WaitGroup      *sync.WaitGroup
	MessageHandler func(data []byte)
}

func (server *Server) handleTcpConnection(conn net.Conn) {
	defer conn.Close()
	//fmt.Printf("HISI: DEVICE CONNECTED: %s\n", conn.RemoteAddr().String())
	var buf bytes.Buffer

	_, err := io.Copy(&buf, conn)
	if err != nil {
		fmt.Printf("HISI: TCP READ ERROR: %s\n", err)
		return
	}
	bufString := buf.String()
	resultString := bufString[strings.IndexByte(bufString, '{'):]
	//fmt.Printf("HISI: DEVICE ALERT: %s\n", resultString)
	var dataMap map[string]interface{}

	if err := json.Unmarshal([]byte(resultString), &dataMap); err != nil {
		fmt.Printf("HISI: JSON PARSE ERROR: %s\n", err)
		return
	}
	if dataMap["Address"] != nil {
		hexAddrStr := fmt.Sprintf("%v", dataMap["Address"])
		dataMap["ipAddr"] = hexIpToCIDR(hexAddrStr)
	}

	jsonBytes, err := json.Marshal(dataMap)
	if err != nil {
		fmt.Printf("HISI: JSON STRINGIFY ERROR: %s\n", err)
		return
	}

	if dataMap["SerialID"] == nil {
		fmt.Println("HISI: UNKNOWN DEVICE SERIAL ID")
		fmt.Println(dataMap)
		return
	}

	server.MessageHandler(jsonBytes)
}

func (server *Server) Start() {
	if server.Port == "" {
		server.Port = "15002"
	}

	go func() {
		defer server.WaitGroup.Done()
		server.WaitGroup.Add(1)

		// run
		tcpListener, err := net.Listen("tcp", ":"+server.Port)
		if err != nil {
			panic(err)
		}
		fmt.Println("Listen port:", server.Port)
		defer tcpListener.Close()

		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				panic(err)
			}
			go server.handleTcpConnection(conn)
		}
	}()
}

func hexIpToCIDR(hexAddr string) string {
	hexAddrStr := fmt.Sprintf("%v", hexAddr)[2:]
	ipAddrHexParts := strings.Split(hexAddrStr, "")

	var decParts []string
	lastPart := ""
	for ind, part := range ipAddrHexParts {
		if ind%2 == 0 {
			lastPart = part
		} else {
			decParts = append(decParts, lastPart+part)
		}
	}
	var strParts []string
	for _, part := range decParts {
		dec, _ := strconv.ParseInt(part, 16, 64)
		// PREPEND RESULT TO SLICE
		strParts = append(strParts, "")
		copy(strParts[1:], strParts)
		strParts[0] = strconv.Itoa(int(dec))
	}
	ipAddr := fmt.Sprintf("%s", strings.Join(strParts[:], "."))
	return ipAddr
}
