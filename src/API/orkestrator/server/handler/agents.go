package handler

import (
	"encoding/json"
	"fmt"
	"github.com/rendizi/daee/src/API/orkestrator/internal/accessible"
	"github.com/rendizi/daee/src/API/orkestrator/server"
	"net/http"
)

func Agents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		server.Error(map[string]interface{}{"message": "method is not allowed"}, w)
		return
	}

	var responseAgents []accessible.AgentJson
	for _, agent := range accessible.Agents {
		_, err := accessible.Ping(agent.Addr)
		if err == nil || err.Error() == "agent is busy" {
			responseAgents = append(responseAgents, accessible.AgentJson{
				Addr:     agent.Addr,
				IsBusy:   agent.IsBusy,
				LastPing: agent.LastPing,
			})
		}
	}

	responseJSON, err := json.Marshal(responseAgents)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
