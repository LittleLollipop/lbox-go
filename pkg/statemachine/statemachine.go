package statemachine

import (
    "fmt"
    "sync"
)

// State interface defines the behavior of a state
type State interface {
    GetName() string
    StateIn() error
    StateOut() error
}

// StateMachineInterface defines the required methods for state machine implementation
type StateMachineInterface interface {
    InitData() error
    CheckStateChange(stateNow, newState State) (bool, error)
}

// StateMap holds all available states
type StateMap struct {
    States map[string]State
}

// StateMachine implements a thread-safe state machine
type StateMachine struct {
    smi       StateMachineInterface
    stateMap  StateMap
    initing   bool
    running   bool
    stateNow  State
    stateLast State
    mu        sync.RWMutex // RWMutex for concurrent access
}

// NewStateMachine creates a new instance of StateMachine
func NewStateMachine(smi StateMachineInterface, stateMap StateMap) *StateMachine {
    return &StateMachine{
        smi:      smi,
        stateMap: stateMap,
        initing:  false,
        running:  false,
    }
}

// Start initializes the state machine with the first state
func (sm *StateMachine) Start(firstState string) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    // Check if already running
    if sm.initing || sm.running {
        return fmt.Errorf("state machine already started")
    }

    // Validate first state
    state, exists := sm.stateMap.States[firstState]
    if !exists {
        return fmt.Errorf("state not found: %s", firstState)
    }

    // Initialize state machine
    sm.initing = true
    if err := sm.smi.InitData(); err != nil {
        sm.initing = false
        return fmt.Errorf("init data failed: %w", err)
    }

    // Enter first state
    sm.stateNow = state
    if err := sm.stateNow.StateIn(); err != nil {
        sm.initing = false
        return fmt.Errorf("state in failed: %w", err)
    }

    // Mark as running
    sm.running = true
    sm.initing = false
    return nil
}

// ChangeState triggers a state transition
func (sm *StateMachine) ChangeState(stateName string) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    // Validate current state
    if !sm.running {
        return fmt.Errorf("state machine not running")
    }

    // Validate target state
    _, exists := sm.stateMap.States[stateName]
    if !exists {
        return fmt.Errorf("state not found: %s", stateName)
    }

    return sm.doChangeState(stateName)
}

// doChangeState performs the actual state transition
// Note: This method assumes the caller holds the lock
func (sm *StateMachine) doChangeState(stateName string) error {
    // Validate current state
    if sm.stateNow == nil {
        return fmt.Errorf("current state is undefined")
    }

    // Get new state and check if transition is allowed
    newState := sm.stateMap.States[stateName]
    canChange, err := sm.smi.CheckStateChange(sm.stateNow, newState)
    if err != nil {
        return fmt.Errorf("check state change failed: %w", err)
    }

    if !canChange {
        return fmt.Errorf("state change not allowed from %s to %s", 
            sm.stateNow.GetName(), stateName)
    }

    // Exit current state
    if err := sm.stateNow.StateOut(); err != nil {
        return fmt.Errorf("state out failed: %w", err)
    }

    // Update states
    sm.stateLast = sm.stateNow
    sm.stateNow = newState

    // Enter new state
    if err := sm.stateNow.StateIn(); err != nil {
        return fmt.Errorf("state in failed: %w", err)
    }

    sm.stateLast = nil
    return nil
}

// GetCurrentState returns the current state
func (sm *StateMachine) GetCurrentState() State {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    return sm.stateNow
}

// IsRunning returns whether the state machine is running
func (sm *StateMachine) IsRunning() bool {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    return sm.running
}