package handler

import (
	"testing"

	"agenthub/backend/internal/model"
)

func TestResolveMessageRouteDirectMention(t *testing.T) {
	sessions := []model.Session{
		{SessionID: "orch", AgentType: "orchestrator", AgentName: "manager"},
		{SessionID: "front", AgentType: "codex", AgentName: "frontend"},
	}

	route, err := resolveMessageRoute(RunTaskReq{
		SessionID: "orch",
		AgentType: "orchestrator",
		Message:   "@frontend fix the button",
	}, sessions)
	if err != nil {
		t.Fatalf("resolveMessageRoute() error = %v", err)
	}
	if route.Mode != routeModeDirect {
		t.Fatalf("route.Mode = %q, want %q", route.Mode, routeModeDirect)
	}
	if route.SessionID != "front" || route.AgentType != "codex" {
		t.Fatalf("route target = %s/%s, want front/codex", route.SessionID, route.AgentType)
	}
	if route.AgentMessage != "fix the button" {
		t.Fatalf("route.AgentMessage = %q", route.AgentMessage)
	}
	if route.DisplayMessage != "@frontend fix the button" {
		t.Fatalf("route.DisplayMessage = %q", route.DisplayMessage)
	}
}

func TestResolveMessageRouteGroupNoMentionUsesOrchestrator(t *testing.T) {
	sessions := []model.Session{
		{SessionID: "orch", AgentType: "orchestrator", AgentName: "manager"},
		{SessionID: "worker", AgentType: "claude-code", AgentName: "worker"},
	}

	route, err := resolveMessageRoute(RunTaskReq{
		SessionID: "orch",
		Message:   "please plan this",
	}, sessions)
	if err != nil {
		t.Fatalf("resolveMessageRoute() error = %v", err)
	}
	if route.Mode != routeModeOrchestrator || route.SessionID != "orch" {
		t.Fatalf("route = %s/%s, want orchestrator/orch", route.Mode, route.SessionID)
	}
}

func TestResolveMessageRouteSingleNoMentionUnchanged(t *testing.T) {
	sessions := []model.Session{
		{SessionID: "solo", AgentType: "codex", AgentName: "Codex"},
	}

	route, err := resolveMessageRoute(RunTaskReq{
		SessionID: "solo",
		Message:   "hello",
	}, sessions)
	if err != nil {
		t.Fatalf("resolveMessageRoute() error = %v", err)
	}
	if route.Mode != routeModeUnchanged || route.SessionID != "solo" || route.AgentType != "codex" {
		t.Fatalf("route = %s/%s/%s, want unchanged/solo/codex", route.Mode, route.SessionID, route.AgentType)
	}
}

func TestResolveMessageRouteSkipUserMessageUnchanged(t *testing.T) {
	sessions := []model.Session{
		{SessionID: "orch", AgentType: "orchestrator", AgentName: "manager"},
		{SessionID: "worker", AgentType: "claude-code", AgentName: "worker"},
	}

	route, err := resolveMessageRoute(RunTaskReq{
		SessionID:       "worker",
		AgentType:       "claude-code",
		Message:         "@manager do not reroute this",
		SkipUserMessage: true,
	}, sessions)
	if err != nil {
		t.Fatalf("resolveMessageRoute() error = %v", err)
	}
	if route.Mode != routeModeUnchanged || route.SessionID != "worker" {
		t.Fatalf("route = %s/%s, want unchanged/worker", route.Mode, route.SessionID)
	}
}

func TestBuildRouteAgentsMakesDuplicateTypesAddressable(t *testing.T) {
	agents := buildRouteAgents([]model.Session{
		{SessionID: "a", AgentType: "claude-code"},
		{SessionID: "b", AgentType: "claude-code"},
	})

	if len(agents) != 2 {
		t.Fatalf("len(agents) = %d, want 2", len(agents))
	}
	if agents[0].RouteID != "claude-code" || agents[1].RouteID != "claude-code-2" {
		t.Fatalf("route ids = %q/%q", agents[0].RouteID, agents[1].RouteID)
	}

	route, err := resolveMessageRoute(RunTaskReq{
		SessionID: "a",
		Message:   "@claude-code-2 check this",
	}, []model.Session{
		{SessionID: "a", AgentType: "claude-code"},
		{SessionID: "b", AgentType: "claude-code"},
	})
	if err != nil {
		t.Fatalf("resolveMessageRoute() error = %v", err)
	}
	if route.SessionID != "b" {
		t.Fatalf("route.SessionID = %q, want b", route.SessionID)
	}
}
