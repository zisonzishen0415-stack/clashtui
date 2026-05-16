# AGENTS.md - Go CLI/TUI Project

Compact guidance for OpenCode sessions.

## Build Commands

```bash
make build       # Build (or: go build -o clashtui .)
make run         # Run (or: go run .)
make install     # Install to PATH
go test ./...    # Run tests (if any)
```

- Go 1.24.0 required
- No lint configured (use `golangci-lint run` manually if needed)

## Architecture (Read These Files First)

| File | Purpose | Lines |
|------|---------|-------|
| `main.go:22-54` | Entry, CLI flags, switch dispatch | ~30 |
| `internal/app/app.go` | Main TUI model (3 tabs) | ~200 |
| `internal/clash/core.go` | Mihomo process management | ~150 |
| `internal/clash/client.go` | REST API client (port 9090) | ~100 |

## Key Patterns

- **BubbleTea**: Message-driven TUI, see `internal/tui/nodes.go`
- **Single instance**: `/tmp/clashtui.pid` + socket IPC
- **CLI flags**: Simple switch in `main.go` (add new flags there)

## Defaults

- Proxy port: 7890
- API port: 9090
- Mihomo core: v1.18.10
- Config dir: `~/.config/clashtui/`

## Config Directory

```
~/.config/clashtui/
├── core/clash      # mihomo binary
├── config.yaml     # Current Clash config
├── settings.json   # User settings
├── proxy.sh        # Terminal proxy (auto-generated)
├── Country.mmdb    # GeoIP
└── geosite.dat     # GeoSite
```

## Runtime Dependencies

- Clipboard: `wl-clipboard` (Wayland) or `xclip`/`xsel` (X11)
- TUN mode: `sudo setcap cap_net_admin+ep ~/.config/clashtui/core/clash`

## Common Tasks

| Task | Approach |
|------|----------|
| Add CLI flag | Add case in `main.go` switch |
| Add TUI feature | New component in `internal/tui/` |
| Modify proxy logic | `internal/proxy/proxy.go` |
| Change API calls | `internal/clash/client.go` |

## See Also

- `CLAUDE.md` - Detailed architecture (180+ lines)
- `README.md` - User documentation, protocol support list