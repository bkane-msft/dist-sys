package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// -- Set

type Set struct {
	set map[int]struct{}
}

func NewSet() Set {
	return Set{
		set: make(map[int]struct{}),
	}
}

func (s *Set) Add(v int) {
	s.set[v] = struct{}{}
}

func (s *Set) Contains(v int) bool {
	_, exists := s.set[v]
	return exists
}

func (s *Set) ToSlice() []int {
	var ret []int
	for k := range s.set {
		ret = append(ret, k)
	}
	return ret
}

func main() {
	n := maelstrom.NewNode()

	// TODO: mutex here?
	var messages = NewSet()

	topology := make(map[string][]string)

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		message := int(body["message"].(float64))

		if !messages.Contains(message) {
			messages.Add(message)

			neighbors := n.NodeIDs()
			for _, neighbor := range neighbors {
				n.Send(neighbor, msg)
			}
		}

		response := map[string]string{
			"type": "broadcast_ok",
		}
		return n.Reply(msg, response)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		response := make(map[string]any)
		response["type"] = "read_ok"
		response["messages"] = messages.ToSlice()

		return n.Reply(msg, response)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {

		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		topology_resp := body["topology"].(map[string]any)
		for node, v := range topology_resp {
			neigbhors := v.([]any)
			for _, v := range neigbhors {
				neighbor := v.(string)

				topology[node] = append(topology[node], neighbor)
			}
		}

		response := map[string]string{
			"type": "topology_ok",
		}
		return n.Reply(msg, response)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

}
