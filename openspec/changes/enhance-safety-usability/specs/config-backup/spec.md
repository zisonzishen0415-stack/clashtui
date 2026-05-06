## ADDED Requirements

### Requirement: Config backup before modification
The system SHALL create a backup of config.yaml before any modification operation.

#### Scenario: Backup before subscription download
- **WHEN** user downloads a new subscription
- **THEN** system creates backup file config.yaml.backup.timestamp before overwriting

#### Scenario: Backup before node import
- **WHEN** user imports nodes from clipboard
- **THEN** system creates backup file config.yaml.backup.timestamp before creating new config

#### Scenario: Backup before TUN mode change
- **WHEN** user toggles TUN mode
- **THEN** system creates backup of current config.yaml before modifying

### Requirement: Backup file naming and retention
The system SHALL use timestamp-based backup naming and retain last 3 backups.

#### Scenario: Backup file naming format
- **WHEN** system creates backup
- **THEN** backup filename follows format config.yaml.backup.YYYYMMDD-HHMMSS

#### Scenario: Backup retention limit
- **WHEN** number of backup files exceeds 3
- **THEN** system deletes oldest backups, keeping only 3 most recent

#### Scenario: Backup cleanup on old backups
- **WHEN** backup file is older than 7 days
- **THEN** system may delete the backup file during cleanup

### Requirement: Config rollback on failure
The system SHALL restore config from backup when operation fails.

#### Scenario: Rollback on core start failure
- **WHEN** core fails to start after config modification
- **THEN** system restores config.yaml from most recent backup

#### Scenario: Rollback on validation failure
- **WHEN** config validation fails after download
- **THEN** system restores previous config from backup and shows error message

#### Scenario: Rollback preserves original state
- **WHEN** rollback is performed
- **THEN** config.yaml matches content from backup file exactly