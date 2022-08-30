package main

import (
	"encoding/json"
	"strings"

	"github.com/TorchedSammy/Econode"
	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.New("$ ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	var server *network

	for {
		line, err := rl.Readline()
		if err != nil { // ctrl d
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		
		args := splitQuote(line)
		switch args[0] {
			case "connect":
				address := args[1]
				if address == "" {
					println("missing server address")
					continue
				}

				var err error
				server, err = connectToNetwork(address)
				if err != nil {
					panic(err)
				}
				server.listenIncoming(func(sr econode.ServerResponse) {
					switch sr.Method {
						case "welcome":
							var welcome econode.WelcomePayload
							jsondata, _ := json.Marshal(sr.Data)
							err := json.Unmarshal(jsondata, &welcome)
							if err != nil {
								panic(err)
							}
							
							message(welcome.MOTD)
					}
				})
		}
	}
}

func splitQuote(str string) []string {
	split := []string{}
	sb := &strings.Builder{}
	quoted := false

	for _, r := range str {
		if r == '"' {
			quoted = !quoted
			sb.WriteRune(r)
		} else if !quoted && r == ' ' {
			split = append(split, sb.String())
			sb.Reset()
		} else {
			sb.WriteRune(r)
		}
	}

	if sb.Len() > 0 {
		split = append(split, sb.String())
	}

	return split
}

func message(str string) {
	// 1 up then to the right, and yes we want the newline
	println("\u001b[1A\u001b[9999C")
	println(str)
}
