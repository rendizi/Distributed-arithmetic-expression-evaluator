package handler

import (
	"encoding/json"
	"fmt"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/internal/accessible"
	"github.com/rendizi/Distributed-arithmetic-expression-evaluator/internal/server"
	"net/http"
)

func Agents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		server.Error(map[string]interface{}{"message": "method is not allowed"}, w)
		return
	}

	//проходимся по всем агентам и проверяем их доступность
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

	//маршалим и возвращаем
	responseJSON, err := json.Marshal(responseAgents)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
