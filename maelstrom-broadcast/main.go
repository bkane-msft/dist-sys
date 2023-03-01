package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	// TODO: mutex here?
	var messages []int

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		message := int(body["message"].(float64))
		messages = append(messages, message)

		response := map[string]string{
			"type": "broadcast_ok",
		}
		return n.Reply(msg, response)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		response := make(map[string]any)
		response["type"] = "read_ok"
		response["messages"] = messages

		return n.Reply(msg, response)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		response := map[string]string{
			"type": "topology_ok",
		}
		return n.Reply(msg, response)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

}
