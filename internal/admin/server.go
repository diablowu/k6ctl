package admin

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type AdminServer struct {
	addr string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		logrus.Debug("CheckOrigin host {} ,uri {}", r.Host, r.URL)
		return true
	},
}

func Staring(addr string) {

	adminServer := AdminServer{addr: addr}

	adminServer.Start()
}

func (as AdminServer) Start() {
	logrus.Debug("admin server staring at {}", as.addr)

	http.HandleFunc("/api/v1/agent", echo)

	http.HandleFunc("/api/v1/admin", admin)

	logrus.Error(http.ListenAndServe(as.addr, nil))
}

var conn *websocket.Conn

func admin(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("m")
	conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func echo(w http.ResponseWriter, r *http.Request) {
	var err error
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Print("upgrade:", err)
		return
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
	}
}
