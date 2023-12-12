package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	TYPES = "picture,video,sticker,url"
)

type Text struct {
	Value    string `json:"value"`
	Receiver string `json:"receiver"`
}

type Message struct {
	Receiver      string `json:"receiver"`
	MinApiVersion uint   `json:"min_api_version"`
	Sender        struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"sender"`
	TrackingData string `json:"tracking_data"`
	Type         string `json:"type"`
	Text         string `json:"text"`
	Media        string `json:"media"`
}

func Send(w http.ResponseWriter, r *http.Request) {
	var text Text
	json.NewDecoder(r.Body).Decode(&text)
	msg := &Message{
		Receiver:      text.Receiver,
		MinApiVersion: 1,
		Sender: struct {
			Name   string `json:"name"`
			Avatar string `json:"avatar"`
		}{
			Name:   "Dushyant",
			Avatar: "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460__340.png",
		},
		TrackingData: "tracking data",
		Type:         "text",
		Text:         text.Value,
	}
	client := &http.Client{}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	req, err := http.NewRequest("POST", "https://chatapi.viber.com/pa/send_message", strings.NewReader(string(payload)))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("X-Viber-Auth-Token", "4f4f3fe501a7e2d3-8625db40ea97ec98-48864325acb52a2d")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("%s | %s", text.Value, text.Receiver)
}

func Viber(w http.ResponseWriter, r *http.Request) {
	var response interface{}
	json.NewDecoder(r.Body).Decode(&response)
	switch res := response.(map[string]interface{}); res["event"] {
	case "webhook":
		log.Println("Webhook called")
	case "subscribed":
		log.Println("Subscribed")
	case "message":
		var txt string
		msg := res["message"].(map[string]interface{})
		msgType := msg["type"].(string)
		sender := res["sender"].(map[string]interface{})
		if msgType == "text" {
			txt = msg["text"].(string)
		} else if strings.Contains(TYPES, msgType) {
			txt = msg["media"].(string)
		}
		log.Printf("%s | %s", txt, sender["id"])
	case "unsubscribed":
		log.Println("Unsubscribed")
	case "delivered":
		// log.Println("transported to the subscriber | IZkeUOZm7xPCZlVnlHRvGw==")
	case "seen":
		log.Println("subscriber acknowldeged transported the data | IZkeUOZm7xPCZlVnlHRvGw==")
	default:
		log.Println(res)
	}
}
