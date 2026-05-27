package handler

import (
	"context"
	"encoding/json"
	"strconv"

	"agenthub/backend/internal/vo"
	"agenthub/backend/pkg/redis"

	"github.com/gin-gonic/gin"
)

type ResourceInfo struct {
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
	Unit  string  `json:"unit"`
}

func (h *AdminHandler) GetResources(c *gin.Context) {
	// disk + memory: delegate to agentend
	var disk ResourceInfo
	var memory ResourceInfo

	resp, err := h.agentClient.GetResources()
	if err == nil {
		defer resp.Body.Close()
		var result struct {
			Disk   ResourceInfo `json:"disk"`
			Memory ResourceInfo `json:"memory"`
		}
		if json.NewDecoder(resp.Body).Decode(&result) == nil {
			disk = result.Disk
			memory = result.Memory
		}
	}

	vo.OK(c, gin.H{
		"disk":   disk,
		"memory": memory,
		"redis":  getRedisUsage(),
	})
}

func getRedisUsage() ResourceInfo {
	client := redis.GetClient()
	if client == nil {
		return ResourceInfo{Used: 0, Total: 0, Unit: "MB"}
	}

	info, err := client.Info(context.Background(), "memory").Result()
	if err != nil {
		return ResourceInfo{Used: 0, Total: 0, Unit: "MB"}
	}

	usedMB := parseRedisInfoFloat(info, "used_memory") / 1e6
	maxMemoryStr := parseRedisInfoString(info, "maxmemory")
	var totalMB float64
	if maxMemoryStr != "" && maxMemoryStr != "0" {
		if val, err := strconv.ParseFloat(maxMemoryStr, 64); err == nil {
			totalMB = val / 1e6
		}
	}
	if totalMB == 0 {
		totalMB = 512
	}

	return ResourceInfo{Used: usedMB, Total: totalMB, Unit: "MB"}
}

func parseRedisInfoFloat(info, key string) float64 {
	target := key + ":"
	for i := 0; i < len(info); i++ {
		if i+len(target) <= len(info) && info[i:i+len(target)] == target {
			j := i + len(target)
			for j < len(info) && info[j] != '\r' && info[j] != '\n' {
				j++
			}
			val, _ := strconv.ParseFloat(info[i+len(target):j], 64)
			return val
		}
	}
	return 0
}

func parseRedisInfoString(info, key string) string {
	target := key + ":"
	for i := 0; i < len(info); i++ {
		if i+len(target) <= len(info) && info[i:i+len(target)] == target {
			j := i + len(target)
			for j < len(info) && info[j] != '\r' && info[j] != '\n' {
				j++
			}
			return info[i+len(target) : j]
		}
	}
	return ""
}
