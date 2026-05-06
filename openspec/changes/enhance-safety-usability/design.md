## Context

ClashTUI is a terminal UI for managing mihomo proxy core. Currently, it directly modifies system state (proxy settings, config files) without safeguards. Errors can leave system in broken state requiring manual recovery via `--restore-network`. The tool needs robustness improvements to prevent network disruption and enable automatic recovery.

Current state:
- Config files overwritten without backup
- System proxy changed without state tracking
- TUN mode enabled without capability verification
- Core process may leave zombie processes
- Errors show cryptic messages without recovery steps

## Goals / Non-Goals

**Goals:**
- Prevent network disruption through validation and safeguards
- Enable automatic rollback on failure
- Improve error clarity with actionable recovery
- Guarantee cleanup on exit (no zombie processes)
- Make tool safe for daily use without fear of breaking network

**Non-Goals:**
- Not adding new proxy features or protocols
- Not changing existing user workflows or UI structure
- Not adding configuration encryption or security features
- Not modifying mihomo core behavior

## Decisions

### 1. Config Backup Strategy: Simple Timestamp Backups
**Decision**: Backup configs with timestamp suffix before overwriting.

**Rationale**: 
- Simple to implement and understand
- No external dependencies
- Easy rollback by copying backup back
- Alternatives considered:
  - Git-based versioning: Overkill, adds git dependency
  - Database-based: Too complex for config files
  - In-memory rollback: Doesn't survive crashes

**Implementation**: Before saving config.yaml, create config.yaml.backup.timestamp. Keep last 3 backups, auto-clean older ones.

### 2. Validation Approach: Pre-flight YAML Parse + Basic Structure Check
**Decision**: Parse YAML and validate required fields before starting core.

**Rationale**:
- Catch syntax errors early (before core crashes)
- Validate minimal structure (proxies, proxy-groups, rules)
- Alternatives considered:
  - Full schema validation: Too strict, rejects valid configs
  - No validation: Current approach, causes silent failures
  - Core dry-run mode: Mihomo doesn't support dry-run

**Implementation**: Parse YAML with go-yaml, check for proxies/proxy-groups/rules sections exist. Reject config if parse fails or structure invalid.

### 3. System Proxy State Tracking: File-based State File
**Decision**: Write state file tracking current proxy settings.

**Rationale**:
- Survives crashes and restarts
- Enables cleanup on next startup if previous session died
- Alternatives considered:
  - In-memory only: Lost on crash
  - systemd tracking: Only works with systemd, not portable
  - Process lock: Already have PID lock, but doesn't track proxy state

**Implementation**: Write ~/.config/clashtui/network-state.json with current proxy mode (system-proxy/tun/off). On startup, if mihomo not running but state file shows system-proxy/tun, cleanup.

### 4. Health Check Strategy: Startup Diagnostics
**Decision**: Run diagnostics on startup to detect broken state.

**Rationale**:
- Proactive detection before user starts operations
- Auto-recover from previous session's failures
- Alternatives considered:
  - No health checks: User discovers issues during operation
  - Background monitoring: Adds complexity, not needed for TUI
  - On-demand checks: Requires user action

**Implementation**: On TUI start, check: mihomo running?, system proxy matches settings?, TUN mode matches settings?. If mismatch, cleanup to safe state.

### 5. Atomic Operations: Transaction Pattern with Rollback
**Decision**: Wrap critical operations in transaction pattern (prepare → commit → rollback on failure).

**Rationale**:
- Prevent partial updates from breaking system
- Clear failure recovery path
- Alternatives considered:
  - Try-catch only: Doesn't handle all failure cases
  - Two-phase commit: Overkill for single-node operations
  - No transactions: Current approach, causes partial failures

**Implementation**: For system proxy changes: 1) Backup current state 2) Apply new settings 3) Verify success 4) If fail, restore backup. For config changes: 1) Backup config 2) Write new config 3) Validate 4) If fail, restore backup.

### 6. Zombie Process Prevention: Proper Process Wait + PID File Cleanup
**Decision**: Always Wait() after Kill(), clear PID file on successful stop.

**Rationale**:
- Guarantees zombie cleanup (POSIX requirement)
- PID file cleanup prevents stale PID reads
- Alternatives considered:
  - Let OS reap: Zombies accumulate, bad practice
  - Ignore zombies: Current approach, causes zombie buildup
  - Double-fork: Overkill for simple process management

**Implementation**: In stopInternal(), always call Process.Wait() after Kill(). Clear PID file only after successful stop. On startup, check for stale PID file and kill if process exists.

### 7. TUN Mode Enhancement: Capability Pre-check + Graceful Fallback
**Decision**: Verify capability before enabling TUN, fallback to system proxy if fail.

**Rationale**:
- Prevent silent failures when capability missing
- Graceful degradation instead of hard failure
- Alternatives considered:
  - Auto-set capability: Requires sudo, security concern
  - Fail hard: Current approach, confusing error
  - No TUN mode: Removes useful feature

**Implementation**: Before enabling TUN, check getcap. If missing, show instructions and keep system-proxy mode. If TUN enabled and core fails, detect DNS failure and disable TUN.

## Risks / Trade-offs

### [Risk] Backup files accumulate over time
**Mitigation**: Auto-clean backups older than 7 days or keep only last 3 backups.

### [Risk] State file corrupted
**Mitigation**: Validate state file on read, if invalid assume safe state (no proxy) and cleanup.

### [Risk] Validation rejects valid configs
**Mitigation**: Minimal validation (syntax + basic structure), allow configs with extra fields.

### [Risk] Health checks slow startup
**Mitigation**: Run checks asynchronously in background, don't block UI startup.

### [Risk] Transaction rollback fails
**Mitigation**: Keep backups and state files as fallback. Manual recovery via --restore-network remains available.

### [Risk] Capability check requires external command
**Mitigation**: getcap is standard on Linux, check for command existence first. If missing, skip check but warn user.

## Open Questions

None - all technical decisions resolved.