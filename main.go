package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "net/url"
	"strings"

	"github.com/gorilla/websocket"
)

//bidding system
//functionalities i want to support:
/* 	1) create a bid (takes an amount and also a timer should start for all the bidders so the timer should be broadcasted to everyone)
*	2) place a bid (ask)
*	3) create a room for bidding
*	4) Join a room for bidding
*
 */
//each auction will be in a specific room(roomid)
//people in a room should be able to interact with one another


var addr = flag.String("addr","localhost:8080","http service address")

var upgrader = websocket.Upgrader{} 



// func createRoom(w http.ResponseWriter, r *http.Request){

// }//should create a room and redirect to joinRoom

type Bid struct {
	Item string
	Bid uint
	Bidder string
}

type AuctionManager struct {
	Item string
	StartingPrice uint
	CurrentBid uint
	Seller string
	TakingBids bool
	Bids []Bid
}

type Room struct {
	RoomId string
	Crowd uint
	Manager AuctionManager
}

type AuctionStartReq struct {
	Item string
	Price uint
	Seller string
}

func (a *AuctionManager) createBid(bid Bid){
	a.Bids = append(a.Bids, bid)
}

func (a *AuctionManager) startAuction(item string, price uint, seller string){
	a.Item = item
	a.StartingPrice = price
	a.CurrentBid = a.StartingPrice
	a.Seller = seller
	a.TakingBids = true
}

func (r *Room) CustomerManager(c *websocket.Conn){
	var currentRoomBid Bid
	for {
		_, message, err := c.ReadMessage()
		fmt.Printf("the message received was: %s\n",message)
		if err != nil {
			log.Println("read: ",err)
			break
		}

		err = json.Unmarshal(message,&currentRoomBid)
		if err!=nil {
			fmt.Println(err)
		}

		r.Manager.createBid(currentRoomBid)
		fmt.Printf("the current bid is: %v\n",currentRoomBid)
		fmt.Printf("bids received so far: %v\n",r.Manager.Bids)

		err = c.WriteMessage(websocket.TextMessage, []byte("your message was received"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}	
}

// func (a *AuctionManager) SellerManager(){
// 	for {
// 		_, message, err := c.ReadMessage()
// 		fmt.Printf("the message received was: %s\n",message)
// 		if err != nil {
// 			log.Println("read: ",err)
// 			break
// 		}

// 		err = json.Unmarshal(message,&currentRoomBid)
// 		if err!=nil {
// 			fmt.Println(err)
// 		}
// 		fmt.Printf("the current bid is: %v\n",currentRoomBid)

// 		err = c.WriteMessage(websocket.TextMessage, []byte("your message was received"))
// 		if err != nil {
// 			log.Println("write:", err)
// 			break
// 		}
// 	}	
// }

var Rooms = make(map[string]Room)

func joinRoom(w http.ResponseWriter, r *http.Request) {
	c, err :=  upgrader.Upgrade(w,r,nil)
	if err!= nil {
		fmt.Println("Error:",err)
		return
	}
	defer c.Close()

	suffix := strings.TrimPrefix(r.URL.Path,"/room/")
	id, err := strconv.Atoi(suffix)
	if err!= nil {
		fmt.Println("invalid id")
		fmt.Println("Error:",err)
		return
	}

	_, exists := Rooms[strconv.Itoa(id)]
	if !exists {
		startReq := AuctionStartReq{}
		c.ReadJSON(&startReq)
		currManager := AuctionManager{}
		currManager.startAuction(startReq.Item,startReq.Price,startReq.Seller)

		Rooms[strconv.Itoa(id)] = Room{
			RoomId: strconv.Itoa(id),
			Crowd: 1,
			Manager: currManager,
		}
		currRoom := Rooms[strconv.Itoa(id)]
		log.Printf("someone started auction at roomId: %v\n",currRoom.RoomId)
		log.Printf("auction item: %v\n",currManager.Item)
		log.Printf("seller: %v\n",currManager.Seller)
		log.Printf("starting Pirce: %v\n",currManager.StartingPrice)
		currRoom.CustomerManager(c)
	}

	currRoom := Rooms[strconv.Itoa(id)]
	currRoom.Crowd++
	log.Printf("someone joined on roomId: %v\n",currRoom.RoomId)
	currRoom.CustomerManager(c)
	
}

func main(){
	flag.Parse()
	http.HandleFunc("/room/",joinRoom)

	
	http.ListenAndServe(*addr,nil)
}