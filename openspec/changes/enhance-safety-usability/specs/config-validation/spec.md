## ADDED Requirements

### Requirement: YAML syntax validation
The system SHALL validate config YAML syntax before starting core.

#### Scenario: Valid YAML passes validation
- **WHEN** config.yaml contains valid YAML syntax
- **THEN** validation passes and allows core start

#### Scenario: Invalid YAML rejected
- **WHEN** config.yaml contains invalid YAML syntax (e.g., malformed indentation)
- **THEN** validation fails and prevents core start with syntax error message

#### Scenario: Empty config rejected
- **WHEN** config.yaml is empty or missing
- **THEN** validation fails with "config missing" error

### Requirement: Config structure validation
The system SHALL validate minimal config structure (proxies, proxy-groups, rules).

#### Scenario: Config with required sections passes
- **WHEN** config.yaml contains proxies, proxy-groups, and rules sections
- **THEN** structure validation passes

#### Scenario: Config missing proxies rejected
- **WHEN** config.yaml lacks proxies section
- **THEN** validation fails with "missing proxies" error

#### Scenario: Config missing proxy-groups rejected
- **WHEN** config.yaml lacks proxy-groups section
- **THEN** validation fails with "missing proxy-groups" error

#### Scenario: Config missing rules rejected
- **WHEN** config.yaml lacks rules section
- **THEN** validation fails with "missing rules" error

### Requirement: Validation before core operations
The system SHALL validate config before all core start operations.

#### Scenario: Validation before manual core start
- **WHEN** user manually starts core
- **THEN** system validates config first, showing error if invalid

#### Scenario: Validation before subscription switch
- **WHEN** user switches subscription
- **THEN** system validates downloaded config before starting core

#### Scenario: Validation before config refresh
- **WHEN** user refreshes subscription
- **THEN** system validates updated config before core restart