## 1. Backup and State Infrastructure

- [x] 1.1 Create internal/backup package with backup file management functions
- [x] 1.2 Implement CreateBackup(filepath string) function with timestamp naming
- [x] 1.3 Implement RestoreBackup(filepath string, backupPath string) function
- [x] 1.4 Implement CleanupOldBackups(filepath string, maxCount int) function
- [x] 1.5 Create internal/state package for network state tracking
- [x] 1.6 Implement NetworkState struct with mode, port, timestamp fields
- [x] 1.7 Implement SaveState(state NetworkState) function to write network-state.json
- [x] 1.8 Implement LoadState() function to read and validate network-state.json
- [x] 1.9 Add state file handling for corrupted/missing files (return safe default)

## 2. Config Validation System

- [x] 2.1 Create internal/validation package
- [x] 2.2 Implement ValidateYAML(configData []byte) function for syntax checking
- [x] 2.3 Implement ValidateStructure(configData []byte) function for required sections check
- [x] 2.4 Integrate validation into config.LoadConfig() flow
- [x] 2.5 Add validation before core.Start() and core.StartAndCheck()
- [x] 2.6 Update DownloadSubscription() to validate config before saving
- [x] 2.7 Add validation error messages with actionable suggestions

## 3. Health Check System

- [x] 3.1 Create internal/health package
- [x] 3.2 Implement CheckCoreRunning() function to verify mihomo via API
- [x] 3.3 Implement CheckSystemProxyState() function to read actual system proxy settings
- [x] 3.4 Implement CheckTUNState() function to verify TUN in config and running
- [x] 3.5 Implement RunHealthChecks() function to orchestrate all checks
- [x] 3.6 Add async health check execution in app.go Init()
- [x] 3.7 Implement auto-recovery logic for state mismatches
- [x] 3.8 Add health check logging and user notifications

## 4. Atomic Operations Pattern

- [x] 4.1 Create internal/transaction package with Transaction struct
- [x] 4.2 Implement Prepare() phase to backup current state
- [x] 4.3 Implement Commit() phase to apply changes
- [x] 4.4 Implement Rollback() phase to restore backup
- [x] 4.5 Add Verify() step for success confirmation
- [x] 4.6 Wrap system proxy operations in transaction pattern
- [x] 4.7 Wrap config modifications in transaction pattern
- [x] 4.8 Wrap TUN mode changes in transaction pattern
- [x] 4.9 Add transaction logging for debugging

## 5. Core Management Improvements

- [x] 5.1 Update core.Stop() to always call Process.Wait() after Kill()
- [x] 5.2 Update cleanupOnExit() to properly wait for process termination
- [x] 5.3 Ensure PID file cleanup only after successful stop
- [x] 5.4 Add stale PID detection on core.New() initialization
- [x] 5.5 Kill and cleanup stale process if found on startup
- [x] 5.6 Add mutex verification in start/stop operations
- [x] 5.7 Update stopAll() and restoreNetwork() to use improved cleanup

## 6. TUN Mode Enhancements

- [x] 6.1 Update NeedsCapability() to check getcap command exists first
- [x] 6.2 Add capability verification before TUN enable attempt
- [x] 6.3 Implement fallback to system-proxy on TUN start failure
- [x] 6.4 Add DNS failure detection during TUN operation
- [x] 6.5 Update handleTUNModeToggle() with graceful fallback logic
- [x] 6.6 Add clear capability error messages with exact setcap command
- [x] 6.7 Ensure TUN config consistency check in health checks
- [x] 6.8 Add TUN-specific error handling and recovery messages

## 7. Error Handling Improvements

- [x] 7.1 Create internal/errors package with error classification
- [x] 7.2 Implement CriticalError, OperationalError, WarningError types
- [x] 7.3 Add error context preservation (logs, stderr capture)
- [x] 7.4 Implement actionable recovery suggestions for each error type
- [x] 7.5 Update all error messages to show recovery steps
- [x] 7.6 Add error visibility persistence in UI status line
- [x] 7.7 Implement auto-cleanup trigger for critical network errors
- [x] 7.8 Update main.go error handlers with improved messages

## 8. Integration and Testing

- [x] 8.1 Integrate backup system into all config modification paths
- [x] 8.2 Integrate state tracking into all proxy mode changes
- [x] 8.3 Integrate validation into all core start paths
- [x] 8.4 Integrate transactions into critical operations
- [x] 8.5 Test backup/rollback with simulated failures
- [x] 8.6 Test state tracking with crash simulations
- [x] 8.7 Test health checks with various mismatch scenarios
- [x] 8.8 Test TUN mode with missing capability
- [x] 8.9 Verify zombie process cleanup works correctly
- [x] 8.10 Manual testing: verify error messages are clear and actionable