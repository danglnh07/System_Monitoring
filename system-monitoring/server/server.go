package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

/*---Variable and type declaration---*/
var (
	// Web socket upgrader
	wsUpgrader = websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		//For simplicity, we will want to return the Check Origin true for all clients connect to the server
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	//Broadcast channel for the server to send message to
	Broadcast = make(chan []byte, 1024)
)

type Server struct {
	sync.Mutex                  //Embedding mutex to avoid race condition
	mux        http.ServeMux    //The server multiplxer
	clients    map[*Client]bool //Map used to keep track of all clients currently connecting to the server
	done       chan struct{}    //Done channel, used for graceful shutdown (not implemented yet)
}

func NewServer() *Server {
	return &Server{
		mux:     *http.NewServeMux(),
		clients: make(map[*Client]bool),
		done:    make(chan struct{}),
	}
}

/*---Handle websocket---*/

func (server *Server) Serve_WebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade web socket request\nError: %v\n", err)
		return
	}

	//If a connection upgrade success, add new publisher
	client := server.AddClient(conn)

	//Start independent goroutine
	go client.SendMessages()
}

func (server *Server) AddClient(conn *websocket.Conn) *Client {
	//Lock the server struct to avoid race condition
	server.Lock()
	defer server.Unlock()

	//Create new client and registered it to the server map
	client := NewClient(conn, server)
	server.clients[client] = true

	fmt.Println("New client join the server")

	return client
}

func (server *Server) RemoveClient(client *Client) {
	//Lock the server struct to avoid race condition
	server.Lock()
	defer server.Unlock()

	//Check if the current client has been registered
	if _, hasRegistered := server.clients[client]; hasRegistered {
		//Close the channel message
		close(client.msgs)
		//Close the connection
		client.conn.Close()
		//Remove publisher from the list of publisher
		delete(server.clients, client)

		fmt.Println("Client lost connection")

		return
	}

	//If this client has not been registered, then print out the message
	fmt.Println("Cannot found this client!")
}

/*---Handle data from service sent to---*/
func HandleHardware(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//If this is the POST method, then this should be the sys-check sending data to server
		rawData, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Fail to read request body")
			fmt.Println(err)
			http.Error(w, "Fail to read request body", http.StatusInternalServerError)
			return
		}

		Broadcast <- rawData
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Data uploaded successfully"))
	}
}

/*---Config server---*/
func (server *Server) Start() {
	/*---Serve the static files---*/

	//Serve the index.html file
	fs := http.FileServer(http.Dir("./web"))
	server.mux.Handle("/", fs)

	//Serve the static resources
	fs = http.FileServer(http.Dir("./web/static"))
	server.mux.Handle("/static/", http.StripPrefix("/static", fs))

	//Handler for upgrading from HTTP to Web Socket
	server.mux.HandleFunc("/ws", server.Serve_WebSocket)

	//Handler for normal HTTP endpoint
	server.mux.HandleFunc("/hardware", HandleHardware)

	//Listen and serve
	fmt.Println("Start server at http://localhost:8080 ...")
	err := http.ListenAndServe(":8080", &server.mux)
	if err != nil {
		fmt.Printf("Failed to start server\nError: %v\n", err)
		os.Exit(1)
	}
}
