package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	chat "google.golang.org/api/chat/v1"
)

func main() {
	log.Printf("starting chatbot")
	http.HandleFunc("/", ChatServer)
	http.ListenAndServe(":8080", nil)
}

func writeResponse(w http.ResponseWriter, output interface{}) {
	bytes, err := json.Marshal(output)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

type Event struct {
	Type    string       `json:"type"`
	Message chat.Message `json:"message"`
}

type TextResponse struct {
	Text string `json:"text"`
}

type CardResponse struct {
	Cards []*chat.Card `json:"cards"`
}

func makeCardResponse(cards ...*chat.Card) *CardResponse {
	return &CardResponse{
		Cards: cards,
	}
}

func ChatServer(w http.ResponseWriter, r *http.Request) {
	var response chat.Card

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.Write([]byte("bad request"))
		return
	}

	//log.Printf("BODY: %+v", strings.ReplaceAll(string(bodyBytes), "\n", " "))

	var incoming Event
	err = json.Unmarshal(bodyBytes, &incoming)
	if err != nil {
		response.Header = &chat.CardHeader{
			Title: "Invalid Request",
		}
		response.Sections = []*chat.Section{
			{
				Widgets: []*chat.WidgetMarkup{
					{
						TextParagraph: &chat.TextParagraph{
							Text: fmt.Sprintf("invalid request: %v", err),
						},
					},
				},
			},
		}
		writeResponse(w, makeCardResponse(&response))
		return
	}

	log.Printf("INCOMING: %+v", incoming)

	// Handle individual commands.
	if incoming.Type == "MESSAGE" {
		log.Printf("COMMAND: %+v", incoming.Message.SlashCommand)
		if incoming.Message.SlashCommand != nil {
			id := incoming.Message.SlashCommand.CommandId

			var textResponse *TextResponse

			if id == 1 {
				log.Printf("COMMAND 1")
				textResponse = &TextResponse{
					Text: "```\nhold my beer...\n         . .\n       .. . *.\n- -_ _-__-0oOo\n _-_ -__ -||||)\n    ______||||______\n~~~~~~~~~~`\"\"'\n```",
				}

			} else if id == 2 {
				log.Printf("COMMAND 2")
				textResponse = &TextResponse{
					Text: "```\n                 //\n                //\n               //\n              //\n      _______||\n ,-'''       ||`-.\n(            ||   )\n|`-..._______,..-'|\n|            ||   |\n|     _______||   |\n|,-'''_ _  ~ ||`-.|\n|  ~ / `-.\\ ,-'\\ ~|\n|`-...___/___,..-'|\n|    `-./-'_ \\/_| |\n| -'  ~~     || -.|\n(   ~      ~   ~~ )\n`-..._______,..-'```",
				}
			}

			if textResponse != nil {
				log.Printf("RESPONSE: %+v", textResponse)
				writeResponse(w, textResponse)
				return
			}
		}
	}

	w.Write([]byte(""))
}
