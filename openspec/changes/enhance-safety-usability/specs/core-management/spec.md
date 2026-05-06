## ADDED Requirements

### Requirement: Zombie process prevention
The system SHALL prevent zombie process accumulation by properly waiting after process kill.

#### Scenario: Wait after kill
- **WHEN** mihomo process is killed
- **THEN** system calls Process.Wait() immediately after Kill() to reap zombie

#### Scenario: Cleanup on stop command
- **WHEN** user stops core via stop command
- **THEN** system kills process, waits, clears PID file, no zombie left

#### Scenario: Cleanup on exit
- **WHEN** TUI exits
- **THEN** cleanupOnExit() kills core process, waits, clears PID file

### Requirement: PID file lifecycle management
The system SHALL manage PID file throughout core lifecycle.

#### Scenario: PID file creation on start
- **WHEN** core process starts successfully
- **THEN** system writes PID to clash.pid file

#### Scenario: PID file cleanup on successful stop
- **WHEN** core stops successfully
- **THEN** system deletes clash.pid file

#### Scenario: PID file validity check
- **WHEN** reading PID file
- **THEN** system verifies process exists before using PID

### Requirement: Stale process detection and cleanup
The system SHALL detect and clean up stale processes from previous sessions.

#### Scenario: Detect stale PID on startup
- **WHEN** TUI starts and clash.pid exists
- **THEN** system checks if PID process still running, kills if stale

#### Scenario: Kill stale process found
- **WHEN** stale process detected from previous session
- **THEN** system kills process, waits, clears PID file before starting new core

#### Scenario: No stale process
- **WHEN** clash.pid missing or process already gone
- **THEN** no cleanup needed, proceed normally

### Requirement: Process management mutex
The system SHALL use mutex to prevent concurrent process operations.

#### Scenario: Mutex on start
- **WHEN** core.Start() called
- **THEN** mutex locks to prevent concurrent start/stop

#### Scenario: Mutex on stop
- **WHEN** core.Stop() called
- **THEN** mutex locks to prevent concurrent operations

#### Scenario: Mutex prevents race conditions
- **WHEN** multiple threads attempt core operations simultaneously
- **THEN** mutex ensures operations execute sequentially