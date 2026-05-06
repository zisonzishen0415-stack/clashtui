## ADDED Requirements

### Requirement: Capability pre-check before TUN enable
The system SHALL verify capability before allowing TUN mode.

#### Scenario: Check capability present
- **WHEN** user attempts to enable TUN mode
- **THEN** system checks if mihomo binary has cap_net_admin capability via getcap

#### Scenario: Capability missing handling
- **WHEN** capability check shows cap_net_admin missing
- **THEN** system prevents TUN enable and shows instructions: "sudo setcap cap_net_admin+ep <binary-path>"

#### Scenario: Capability present allows TUN
- **WHEN** capability check shows cap_net_admin present
- **THEN** system allows TUN mode enable

### Requirement: Graceful fallback on TUN failure
The system SHALL fallback to system-proxy mode if TUN fails.

#### Scenario: Fallback on core start failure
- **WHEN** core fails to start with TUN mode
- **THEN** system disables TUN in config, reverts to system-proxy mode

#### Scenario: Fallback on DNS failure
- **WHEN** TUN mode enabled but DNS fails (cannot resolve domains)
- **THEN** system detects failure, disables TUN, shows warning

#### Scenario: Keep system-proxy on TUN failure
- **WHEN** TUN mode fails to activate
- **THEN** system-proxy remains enabled (not disabled during TUN attempt)

### Requirement: TUN mode state consistency
The system SHALL maintain consistent TUN state between config and running core.

#### Scenario: TUN config matches settings
- **WHEN** TUN mode enabled in settings
- **THEN** config.yaml includes TUN section, core running with TUN

#### Scenario: TUN disabled on mismatch
- **WHEN** settings show TUN disabled but config has TUN section
- **THEN** system removes TUN section from config on startup

#### Scenario: TUN disabled on core stop
- **WHEN** core stops while TUN mode enabled
- **THEN** system updates state file to reflect network mode

### Requirement: TUN mode error messages
The system SHALL show clear errors for TUN-specific failures.

#### Scenario: Capability missing error
- **WHEN** TUN enable attempted without capability
- **THEN** error shows: "⚠ TUN needs sudo setcap. Run: sudo setcap cap_net_admin+ep <path>" with full path

#### Scenario: TUN start failure error
- **WHEN** core fails to start with TUN
- **THEN** error shows: "⚠ TUN mode failed, reverted to system-proxy" with recovery info

#### Scenario: DNS failure error
- **WHEN** TUN causes DNS failure
- **THEN** error shows: "⚠ TUN DNS failure detected, switched to system-proxy. If DNS still broken: sudo systemctl restart systemd-resolved"