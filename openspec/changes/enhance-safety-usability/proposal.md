## Why

Current tool lacks safety guards and robustness: system proxy changes without rollback, TUN mode can break DNS, corrupted configs cause silent failures, and error recovery requires manual intervention. Users risk network disruption with no automatic safeguards. This makes the tool feel unsafe and impractical for daily use.

## What Changes

- Add config validation before core start (reject invalid configs early)
- Backup configs before overwriting (enable rollback on failure)
- Robust system proxy management with state tracking and auto-rollback
- Enhanced TUN mode: verify capability before enabling, graceful fallback on failure
- Zombie process prevention and cleanup guarantees
- Clear error messages with actionable recovery steps
- Atomic operations for critical state changes (prevent partial updates)
- Health checks on startup (detect broken state from previous sessions)
- Safe subscription download with validation and retry logic

## Capabilities

### New Capabilities
- `config-backup`: Backup and rollback mechanism for config.yaml before modifications
- `config-validation`: Pre-flight validation of clash config before starting core
- `network-state-tracking`: Track system proxy state for reliable cleanup on exit/failure
- `health-checks`: Startup diagnostics to detect and auto-recover from broken states
- `atomic-operations`: Critical state changes wrapped in atomic transactions

### Modified Capabilities
- `core-management`: Enhanced process lifecycle with zombie prevention and cleanup guarantees
- `tun-mode`: Add capability verification, graceful fallback, and DNS protection
- `error-handling**: Clear actionable error messages instead of cryptic failures

## Impact

- internal/clash/core.go: Add validation, health checks, process cleanup improvements
- internal/config/: Add backup/rollback mechanism
- internal/proxy/: Add state tracking, atomic operations
- internal/app/app.go: Add startup health checks, error message improvements
- main.go: Enhance emergency restore, cleanupOnExit reliability
- User experience: Safer operations, clearer errors, automatic recovery