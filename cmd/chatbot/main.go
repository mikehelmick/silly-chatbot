package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/common-nighthawk/go-figure"
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
		var response chat.Card
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

	// Handle individual commands.
	if incoming.Type == "MESSAGE" {
		if incoming.Message.SlashCommand != nil {
			id := incoming.Message.SlashCommand.CommandId

			log.Printf("here")
			var textResponse *TextResponse
			user := fmt.Sprintf("<%s> says:\n", incoming.Message.Sender.Name)
			log.Printf("user: %q", user)

			if id == 1 {
				textResponse = &TextResponse{
					Text: "```\nhold my beer...\n         . .\n       .. . *.\n- -_ _-__-0oOo\n _-_ -__ -||||)\n    ______||||______\n~~~~~~~~~~`\"\"'\n```",
				}

			} else if id == 2 {
				textResponse = &TextResponse{
					Text: "```\n                 //\n                //\n               //\n              //\n      _______||\n ,-'''       ||`-.\n(            ||   )\n|`-..._______,..-'|\n|            ||   |\n|     _______||   |\n|,-'''_ _  ~ ||`-.|\n|  ~ / `-.\\ ,-'\\ ~|\n|`-...___/___,..-'|\n|    `-./-'_ \\/_| |\n| -'  ~~     || -.|\n(   ~      ~   ~~ )\n`-..._______,..-'```",
				}
			} else if id == 3 {
				asciiMsg := strings.Builder{}
				parts := strings.Split(strings.TrimSpace(incoming.Message.ArgumentText), " ")
				for _, p := range parts {
					asciiMsg.WriteString("\n")
					myFigure := figure.NewFigure(p, "", true)
					asciiMsg.WriteString(myFigure.String())
				}

				textResponse = &TextResponse{
					Text: fmt.Sprintf("```%s\n```", asciiMsg.String()),
				}
			}

			if textResponse != nil {
				writeResponse(w, fmt.Sprintf("%s%s", user, textResponse))
				return
			}
		}
	}

	{
		var response chat.Card
		response.Header = &chat.CardHeader{
			Title: "Hello! Nice to meet you!",
		}
		response.Sections = []*chat.Section{
			{
				Widgets: []*chat.WidgetMarkup{
					{
						TextParagraph: &chat.TextParagraph{
							Text: `Here is what I can help with<br><b>/drink</b> - display a drink<br><b>/holdmybeer</b> - tell the room to hold your beer<br><b>/ascii</b> - ascii print a message`,
						},
					},
				},
			},
		}
		writeResponse(w, makeCardResponse(&response))
	}
}
