## ADDED Requirements

### Requirement: Clear actionable error messages
The system SHALL show errors with actionable recovery steps.

#### Scenario: Config validation error
- **WHEN** config validation fails
- **THEN** error shows specific issue (e.g., "missing proxies section") and suggests: "Check subscription or import valid nodes"

#### Scenario: Core start error
- **WHEN** core fails to start
- **THEN** error shows stderr output and suggests: "Run 'clashtui --restore-network' if network broken"

#### Scenario: Subscription download error
- **WHEN** subscription download fails
- **THEN** error shows status code or network error and suggests: "Check URL or try refresh again"

#### Scenario: Capability missing error
- **WHEN** TUN mode fails due to missing capability
- **THEN** error shows exact command to fix: "sudo setcap cap_net_admin+ep ~/.config/clashtui/core/clash"

### Requirement: Error context preservation
The system SHALL preserve error context for debugging.

#### Scenario: Log error details
- **WHEN** error occurs
- **THEN** system logs full error details to logs tab

#### Scenario: Show stderr on core failure
- **WHEN** core crashes on start
- **THEN** system reads and shows clash.err file content in error message

#### Scenario: Keep error visible
- **WHEN** error displayed
- **THEN** error remains visible in status line until user takes action or clears

### Requirement: Error recovery suggestions
The system SHALL suggest specific recovery steps for each error type.

#### Scenario: Network broken recovery
- **WHEN** network-related error occurs
- **THEN** system suggests: "Run 'clashtui --restore-network' to fix network"

#### Scenario: Config error recovery
- **WHEN** config validation error occurs
- **THEN** system suggests: "Delete subscription and re-add, or import nodes manually"

#### Scenario: Core error recovery
- **WHEN** core process error occurs
- **THEN** system suggests: "Stop core (x) and try restart, or check logs"

#### Scenario: Generic fallback recovery
- **WHEN** unknown error occurs
- **THEN** system suggests: "Run 'clashtui --restore-network' as last resort"

### Requirement: Error classification
The system SHALL classify errors by severity and impact.

#### Scenario: Critical network errors
- **WHEN** error breaks network connectivity
- **THEN** system shows with ⚠ icon and red status, auto-triggers cleanup

#### Scenario: Non-critical operation errors
- **WHEN** error affects operation but not network (e.g., node switch fails)
- **THEN** system shows with ⚠ icon but no auto-cleanup

#### Scenario: Warning level issues
- **WHEN** issue is recoverable without user action (e.g., auto-rollback succeeded)
- **THEN** system shows with ⚡ icon in logs only