## ADDED Requirements

### Requirement: Transaction pattern for critical operations
The system SHALL wrap critical operations in prepare-commit-rollback pattern.

#### Scenario: System proxy change transaction
- **WHEN** system proxy setting changes
- **THEN** system follows steps: 1) backup current state 2) apply new setting 3) verify success 4) rollback if fail

#### Scenario: Config modification transaction
- **WHEN** config.yaml is modified
- **THEN** system follows steps: 1) backup config 2) write new config 3) validate new config 4) rollback if fail

#### Scenario: TUN mode change transaction
- **WHEN** TUN mode toggled
- **THEN** system follows steps: 1) backup config 2) modify config for TUN 3) restart core 4) verify success 5) rollback config if fail

### Requirement: Transaction rollback on any failure
The system SHALL rollback transaction on any failure during commit phase.

#### Scenario: Rollback on verification failure
- **WHEN** verification step fails during transaction
- **THEN** system restores backup from prepare phase

#### Scenario: Rollback on core start failure
- **WHEN** core fails to start during config transaction
- **THEN** system restores config backup and shows error

#### Scenario: Rollback on validation failure
- **WHEN** config validation fails during transaction
- **THEN** system restores previous config from backup

### Requirement: Transaction success guarantees
The system SHALL only commit transaction when all steps succeed.

#### Scenario: Commit only on full success
- **WHEN** all transaction steps complete successfully
- **THEN** system finalizes transaction and cleans up backup (keeps for rollback window)

#### Scenario: Partial success treated as failure
- **WHEN** some steps succeed but verification fails
- **THEN** system treats as failure and rolls back

#### Scenario: Transaction atomic guarantee
- **WHEN** transaction completes (success or rollback)
- **THEN** system state is consistent (either fully changed or unchanged, no partial state)

### Requirement: Transaction logging
The system SHALL log transaction progress for debugging.

#### Scenario: Log transaction steps
- **WHEN** transaction executes
- **THEN** each step logged: prepare, commit attempt, success/rollback

#### Scenario: Log rollback reason
- **WHEN** rollback occurs
- **THEN** system logs which step failed and why

#### Scenario: Log transaction outcome
- **WHEN** transaction completes
- **THEN** system logs final outcome: "✓ [operation] completed" or "⚠ [operation] failed, rolled back"