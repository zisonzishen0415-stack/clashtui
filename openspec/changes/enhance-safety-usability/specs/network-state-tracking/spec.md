## ADDED Requirements

### Requirement: Network state file creation
The system SHALL maintain a network state file tracking current proxy mode.

#### Scenario: State file on system proxy enable
- **WHEN** system proxy is enabled
- **THEN** network-state.json records mode as "system-proxy" with port number

#### Scenario: State file on TUN mode enable
- **WHEN** TUN mode is enabled
- **THEN** network-state.json records mode as "tun"

#### Scenario: State file on proxy disable
- **WHEN** system proxy or TUN is disabled
- **THEN** network-state.json records mode as "off"

### Requirement: State file format and location
The system SHALL use JSON format at ~/.config/clashtui/network-state.json.

#### Scenario: State file structure
- **WHEN** state file is written
- **THEN** file contains JSON with fields: mode (string), port (int if system-proxy), timestamp (int64)

#### Scenario: State file location
- **WHEN** state file is created
- **THEN** file is located at ~/.config/clashtui/network-state.json

### Requirement: State-based cleanup on startup
The system SHALL detect mismatch between state file and actual system state, then cleanup.

#### Scenario: Cleanup stale system proxy
- **WHEN** TUI starts and state file shows "system-proxy" but mihomo not running
- **THEN** system clears system proxy settings and updates state to "off"

#### Scenario: Cleanup stale TUN mode
- **WHEN** TUI starts and state file shows "tun" but mihomo not running
- **THEN** system disables TUN in config and updates state to "off"

#### Scenario: No cleanup when state matches
- **WHEN** TUI starts and state file shows "off" and mihomo not running
- **THEN** no cleanup performed

### Requirement: State file validation on read
The system SHALL validate state file integrity when reading.

#### Scenario: Valid state file read
- **WHEN** state file contains valid JSON with required fields
- **THEN** system reads state successfully

#### Scenario: Corrupted state file handling
- **WHEN** state file is corrupted or invalid JSON
- **THEN** system treats as "off" mode and performs cleanup to ensure safe state

#### Scenario: Missing state file handling
- **WHEN** state file does not exist
- **THEN** system assumes "off" mode (safe default)