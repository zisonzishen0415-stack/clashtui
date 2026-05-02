# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Build (Go 1.24.0 required)
go build -o clashtui .

# Run
go run .

# Install to PATH
go install .

# Or use Makefile
make build && make run
```

No tests or lint commands configured in this project.

## CLI Commands

| Command | Description |
|---------|-------------|
| `clashtui` | Launch TUI interface |
| `clashtui --status` | Output JSON status for Waybar |
| `clashtui --daemon` | Background mode (for systemd service) |
| `clashtui --stop` | Stop mihomo, clear system proxy |
| `clashtui --toggle` | Toggle proxy on/off |
| `clashtui --restore-network` | Restore network after reboot (kills mihomo, clears proxy) |
| `clashtui --env` | Print proxy environment variables |

## Architecture Overview

ClashTUI is a terminal UI for managing Clash/mihomo proxy. Uses BubbleTea framework with message-driven architecture.

**Default ports:** Proxy 7890, API 9090, Mihomo core v1.18.10

### Package Structure

- `main.go` - Entry point; CLI flag handling + single-instance check
- `internal/app/app.go` - Main Model with 3 tabs (Nodes/Config/Logs); handles key events, orchestrates components
- `internal/clash/` - Clash integration:
  - `core.go` - Process management (download, start/stop mihomo binary, geo data); subscription parsing; DNS config replacement
  - `client.go` - REST API client for Clash external controller (127.0.0.1:9090)
  - `proxy.go` - ProxyInfo type and API methods (GetAllProxies, SwitchProxy, TestDelay)
- `internal/settings/settings.go` - User settings persistence (subscriptions, ports, toggles)
- `internal/config/config.go` - File paths (~/.config/clashtui/), config.yaml handling
- `internal/proxy/proxy.go` - System proxy via gsettings (GNOME) and kwriteconfig (KDE); creates proxy.sh for terminal
- `internal/clipboard/clipboard.go` - Clipboard read via wl-paste (Wayland) or xclip/xsel (X11)
- `internal/tui/` - BubbleTea components:
  - `nodes.go` - Proxy list with selection, delay testing, auto-test on load
  - `logs.go` - Log display (thread-safe, max 100 lines)
  - `styles.go` - Lipgloss styling definitions
- `internal/singleinstance/singleinstance.go` - PID file mechanism (/tmp/clashtui.pid)

### Key Message Types (internal/tui/nodes.go)

- `MsgProxiesLoaded` - Proxy list loaded from API
- `MsgProxySwitched` - Proxy selection changed
- `MsgDelayTested` - Single delay test result
- `MsgRefresh` - Trigger proxy reload (core started)
- `MsgLogLine` - Add log entry
- `MsgStopCore` - Stop core signal

### Data Flow

1. User imports subscription (clipboard/manual) → `clash.DownloadSubscription()` parses base64 links, builds config.yaml
2. Core starts → `clash.Core.Start()` spawns mihomo process with `-d` pointing to config dir
3. Nodes tab → `clash.Client.GetAllProxies()` fetches from API, auto-starts sequential delay testing
4. Proxy switch → `clash.Client.SwitchProxy()` calls PUT /proxies/Auto
5. System proxy → `proxy.SetSystemProxy()` sets gsettings (GNOME) or kwriteconfig (KDE)

### Startup Behavior (app.go:New())

Critical: On startup, if `s.SystemProxy` is true but mihomo is not running, stale proxy settings are cleared. This handles the case where gsettings persist across reboot but mihomo doesn't auto-start.

### DNS Configuration

Uses `redir-host` mode (not `fake-ip`) to prevent DNS hijacking issues. The `replaceDNSInConfig()` function in core.go enforces safe DNS settings when processing subscription configs.

Default DNS servers: 223.5.5.5, 119.29.29.29 (Chinese public DNS)
Fallback: 1.1.1.1, dns.google (for international resolution)

## Config Directory

```
~/.config/clashtui/
├── core/clash          # mihomo binary
├── config.yaml         # Current Clash config (auto-generated)
├── settings.json       # User settings (subscriptions, ports, toggles)
├── proxy.sh            # Terminal proxy script (auto-generated)
├── clash.pid           # Mihomo process PID (for cleanup)
├── Country.mmdb        # GeoIP database
└── geosite.dat         # GeoSite data
```

## Settings Structure (settings.json)

```json
{
  "subscriptions": [{"name": "...", "url": "...", "traffic": "...", "expiry": "..."}],
  "active_sub_idx": 0,
  "auto_start": false,
  "auto_test_delay": true,
  "auto_select_best": true,
  "system_proxy": true,
  "proxy_port": 7890,
  "api_port": 9090
}
```

## Runtime Requirements

- Go 1.24.0+
- Clipboard: `wl-clipboard` (Wayland) or `xclip`/`xsel` (X11)
- TUN mode: `sudo setcap cap_net_admin+ep ~/.config/clashtui/core/clash`

## Protocol Support

Subscription parsing in `clash/core.go:parseNodeConfig()` handles:
- Trojan, VLESS, VMess, Shadowsocks, ShadowsocksR
- Hysteria2/hy2, Hysteria
- SOCKS5, HTTP/HTTPS proxy
- WireGuard, TUIC, SSH

## Known Issue: Network Recovery After Reboot

If network is broken after reboot (gsettings proxy persists but mihomo doesn't run):
- Run `clashtui --restore-network` to kill stale mihomo and clear proxy settings
- If DNS still broken, may need: `sudo systemctl restart NetworkManager`