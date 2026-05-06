package health

import (
	"fmt"
	"strings"

	"clashtui/internal/clash"
	"clashtui/internal/config"
	"clashtui/internal/proxy"
	"clashtui/internal/settings"
	"clashtui/internal/state"
)

type HealthCheckResult struct {
	Issue       string
	ActionTaken string
}

func CheckCoreRunning(apiPort int) bool {
	client := clash.NewClient(apiPort)
	return client.IsConnected()
}

func CheckSystemProxyState(expectedPort int) (bool, int) {
	currentPort := proxy.GetSystemProxyPort()
	return currentPort == expectedPort, currentPort
}

func CheckTUNState() bool {
	if !config.Exists() {
		return false
	}

	data, err := config.LoadConfigNoValidation()
	if err != nil {
		return false
	}

	content := string(data)
	return strings.Contains(content, "tun:") && strings.Contains(content, "enable: true")
}

func RunHealthChecks(s settings.Settings) []HealthCheckResult {
	results := []HealthCheckResult{}
	coreRunning := CheckCoreRunning(s.APIPort)

	if s.SystemProxy && !coreRunning {
		proxy.UnsetSystemProxy()
		state.SaveState(state.NetworkState{Mode: state.ModeOff})
		results = append(results, HealthCheckResult{
			Issue:       "stale system proxy",
			ActionTaken: "cleared system proxy and updated state to off",
		})
	}

	if s.TUNMode && !coreRunning {
		proxy.UnsetSystemProxy()
		data, _ := config.LoadConfigNoValidation()
		newData := clash.ProcessConfigForTUN(data, false)
		config.SaveConfig(newData)
		state.SaveState(state.NetworkState{Mode: state.ModeOff})
		results = append(results, HealthCheckResult{
			Issue:       "stale TUN mode",
			ActionTaken: "disabled TUN in config and cleared system proxy",
		})
	}

	if s.SystemProxy && coreRunning {
		matches, actualPort := CheckSystemProxyState(s.ProxyPort)
		if !matches && actualPort > 0 {
			proxy.SetSystemProxy(s.ProxyPort)
			results = append(results, HealthCheckResult{
				Issue:       fmt.Sprintf("system proxy port mismatch (expected %d, got %d)", s.ProxyPort, actualPort),
				ActionTaken: fmt.Sprintf("re-applied correct proxy port %d", s.ProxyPort),
			})
		}
	}

	return results
}