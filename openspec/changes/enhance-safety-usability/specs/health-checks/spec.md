## ADDED Requirements

### Requirement: Startup health diagnostics
The system SHALL run health checks on TUI startup to detect broken state.

#### Scenario: Check mihomo running status
- **WHEN** TUI starts
- **THEN** system checks if mihomo process is running via API connection

#### Scenario: Check system proxy consistency
- **WHEN** TUI starts and settings show system-proxy enabled
- **THEN** system verifies actual system proxy matches settings port

#### Scenario: Check TUN mode consistency
- **WHEN** TUI starts and settings show TUN mode enabled
- **THEN** system verifies mihomo is running with TUN enabled in config

### Requirement: Automatic recovery from mismatch
The system SHALL auto-recover when health check detects mismatch.

#### Scenario: Recover from stale proxy
- **WHEN** health check detects system proxy enabled but mihomo stopped
- **THEN** system clears system proxy and updates state file

#### Scenario: Recover from stale TUN
- **WHEN** health check detects TUN mode in settings but mihomo stopped
- **THEN** system clears TUN config and updates state file

#### Scenario: Recover from proxy mismatch
- **WHEN** health check detects system proxy port differs from settings
- **THEN** system re-applies correct proxy port from settings

### Requirement: Non-blocking health checks
The system SHALL run health checks asynchronously without blocking UI startup.

#### Scenario: UI starts immediately
- **WHEN** TUI initializes
- **THEN** UI displays within 1 second while health checks run in background

#### Scenario: Health check results shown later
- **WHEN** health checks complete
- **THEN** results displayed as log message or status update

### Requirement: Health check error reporting
The system SHALL report health check findings to user clearly.

#### Scenario: Report mismatch found
- **WHEN** health check detects mismatch
- **THEN** system shows log message: "⚠ Found stale [system-proxy/TUN], cleaned up"

#### Scenario: Report no issues
- **WHEN** health checks pass with no mismatch
- **THEN** no special message shown (normal startup)

#### Scenario: Report cleanup action
- **WHEN** health check triggers cleanup
- **THEN** system shows recovery action in log: "✓ Restored network to safe state"