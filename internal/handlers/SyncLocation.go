package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"radar/internal/tools"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type LocationMessage struct {
	Latitude     float64
	Longitude    float64
	SubscriberID string
}

func SyncLocation(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.CloseNow()
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var roomID = chi.URLParam(r, "roomID")
	// TODO: verify user has access to room
	var requestID = middleware.GetReqID(ctx)
	fmt.Println("requestID ", requestID, "roomID", roomID)

	subscriber := tools.RedisClient.Subscribe(ctx, roomID)
	var wg sync.WaitGroup
	wg.Add(2)
	go receiveSubscriptionLoop(subscriber, &ctx, conn, &wg)
	go readClientLocationLoop(roomID, requestID, &ctx, conn, &wg)
	wg.Wait()
}

func receiveSubscriptionLoop(subscriber *redis.PubSub, ctx *context.Context, conn *websocket.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg, err := subscriber.ReceiveMessage(*ctx)
		if err != nil {
			log.Printf("error received while reading pubsub: %v", err.Error())
		}

		conn.Write(*ctx, websocket.MessageText, []byte(msg.Payload))
	}
}

func readClientLocationLoop(roomID string, subscriberID string, ctx *context.Context, conn *websocket.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		var err error
		var location map[string]float64
		err = wsjson.Read(*ctx, conn, &location)
		if err != nil {
			log.Printf("error received while reading client location: %v", err.Error())
			break
		}
		var locationMessage = LocationMessage{Longitude: location["longitude"], Latitude: location["latitude"], SubscriberID: subscriberID}
		var payload []byte
		payload, err = json.Marshal(locationMessage)
		if err != nil {
			log.Printf("error received while marshalling client location: %v", err.Error())
			break
		}
		var result = tools.RedisClient.Publish(
			*ctx,
			roomID,
			payload,
		)
		log.Println(result)
	}
}
