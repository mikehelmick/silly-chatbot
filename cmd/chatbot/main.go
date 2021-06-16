package main

import (
	"encoding/json"
	"fmt"
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

func ChatServer(w http.ResponseWriter, r *http.Request) {
	response := make([]chat.Card, 1)
	respCard := response[0]

	var incoming chat.Message
	err := json.NewDecoder(r.Body).Decode(&incoming)
	if err != nil {
		respCard.Header = &chat.CardHeader{
			Title: "Invalid Request",
		}
		respCard.Sections = []*chat.Section{
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
		writeResponse(w, &response)
		return
	}

	// Handle individual commands.
	if incoming.SlashCommand != nil {
		id := incoming.SlashCommand.CommandId

		if id == 1 {
			respCard.Sections = []*chat.Section{
				{
					Widgets: []*chat.WidgetMarkup{
						{
							TextParagraph: &chat.TextParagraph{
								Text: "<code>\nhold my beer...\n         . .\n       .. . *.\n- -_ _-__-0oOo\n _-_ -__ -||||)\n    ______||||______\n~~~~~~~~~~`\"\"'</code>",
							},
						},
					},
				},
			}

		} else if id == 2 {
			respCard.Sections = []*chat.Section{
				{
					Widgets: []*chat.WidgetMarkup{
						{
							TextParagraph: &chat.TextParagraph{
								Text: "<code>\n                 //\n                //\n               //\n              //\n      _______||\n ,-'''       ||`-.\n(            ||   )\n|`-..._______,..-'|\n|            ||   |\n|     _______||   |\n|,-'''_ _  ~ ||`-.|\n|  ~ / `-.\\ ,-'\\ ~|\n|`-...___/___,..-'|\n|    `-./-'_ \\/_| |\n| -'  ~~     || -.|\n(   ~      ~   ~~ )\n`-..._______,..-'</code>",
							},
						},
					},
				},
			}
		}

		if len(respCard.Sections) > 0 {
			writeResponse(w, &response)
			return
		}
	}

	w.Write([]byte(""))
}
