package main

import (
	"flag"
	"github.com/sevlyar/go-daemon"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", websocketServer, "resource-collector server address")

func main() {
	cntxt := &daemon.Context{
		PidFileName: "rc.pid",
		PidFilePerm: 0644,
		LogFileName: "rc.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Env:         nil,
		Args:        []string{"[resource-collector server daemon]"},
		Umask:       027,
	}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("resource-collector server started")

	listen_main()
}

func listen_main() {
	flag.Parse()

	go serverSocketCreate()
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		serveWs(writer, request)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		Error.Println("ListenAndServe: ", err)
	}
}

func init() {

	Trace = log.New(os.Stdout,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	/*Error = log.New(io.MultiWriter(file, os.Stderr),
	"ERROR: ",
	log.Ldate|log.Ltime|log.Lshortfile)*/
	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
