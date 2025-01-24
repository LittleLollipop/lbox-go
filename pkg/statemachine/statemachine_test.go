package statemachine

import (
    "errors"
    "fmt"
    "testing"
)

// GameState implements State interface for testing
type GameState struct {
    name string
}

func (s *GameState) GetName() string {
    return s.name
}

func (s *GameState) StateIn() error {
    fmt.Printf("Entering state: %s\n", s.name)
    return nil
}

func (s *GameState) StateOut() error {
    fmt.Printf("Exiting state: %s\n", s.name)
    return nil
}

// GameStateMachine implements StateMachineInterface for testing
type GameStateMachine struct{}

func (gsm *GameStateMachine) InitData() error {
    fmt.Println("Initializing game data...")
    return nil
}

func (gsm *GameStateMachine) CheckStateChange(stateNow, newState State) (bool, error) {
    fmt.Printf("Checking state change: %s -> %s\n", stateNow.GetName(), newState.GetName())
    return true, nil
}

// MockState implements State interface for testing
type MockState struct {
    name           string
    stateInError  error
    stateOutError error
}

func (s *MockState) GetName() string {
    return s.name
}

func (s *MockState) StateIn() error {
    return s.stateInError
}

func (s *MockState) StateOut() error {
    return s.stateOutError
}

// MockStateMachine implements StateMachineInterface for testing
type MockStateMachine struct {
    initError        error
    checkChangeError error
    allowChange      bool
}

func (m *MockStateMachine) InitData() error {
    return m.initError
}

func (m *MockStateMachine) CheckStateChange(stateNow, newState State) (bool, error) {
    return m.allowChange, m.checkChangeError
}

func TestStateMachine_Start(t *testing.T) {
    tests := []struct {
        name          string
        smi           *MockStateMachine
        firstState    string
        stateInError  error
        expectError   bool
    }{
        {
            name:        "Normal start",
            smi:         &MockStateMachine{},
            firstState:  "state1",
            expectError: false,
        },
        {
            name:        "State not found",
            smi:         &MockStateMachine{},
            firstState:  "nonexistent",
            expectError: true,
        },
        {
            name:        "Init failed",
            smi:         &MockStateMachine{initError: errors.New("init failed")},
            firstState:  "state1",
            expectError: true,
        },
        {
            name:          "StateIn failed",
            smi:           &MockStateMachine{},
            firstState:    "state1",
            stateInError:  errors.New("state in failed"),
            expectError:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            state1 := &MockState{name: "state1", stateInError: tt.stateInError}
            stateMap := StateMap{
                States: map[string]State{
                    "state1": state1,
                },
            }
            
            sm := NewStateMachine(tt.smi, stateMap)
            err := sm.Start(tt.firstState)
            
            if tt.expectError && err == nil {
                t.Error("Expected error but got nil")
            }
            if !tt.expectError && err != nil {
                t.Errorf("Expected no error but got: %v", err)
            }
        })
    }

    // Test duplicate start
    t.Run("Duplicate start", func(t *testing.T) {
        state1 := &MockState{name: "state1"}
        stateMap := StateMap{
            States: map[string]State{
                "state1": state1,
            },
        }
        sm := NewStateMachine(&MockStateMachine{}, stateMap)
        
        // First start
        if err := sm.Start("state1"); err != nil {
            t.Fatalf("First start failed: %v", err)
        }
        
        // Second start should fail
        if err := sm.Start("state1"); err == nil {
            t.Error("Expected duplicate start to fail, but succeeded")
        }
    })
}

func TestStateMachine_ChangeState(t *testing.T) {
    tests := []struct {
        name            string
        smi             *MockStateMachine
        fromState       string
        toState         string
        allowChange     bool
        checkChangeErr  error
        stateOutError   error
        stateInError    error
        expectError     bool
    }{
        {
            name:        "Normal state change",
            fromState:   "state1",
            toState:     "state2",
            allowChange: true,
            expectError: false,
        },
        {
            name:        "Target state not found",
            fromState:   "state1",
            toState:     "nonexistent",
            allowChange: true,
            expectError: true,
        },
        {
            name:           "State check failed",
            fromState:      "state1",
            toState:        "state2",
            checkChangeErr: errors.New("check failed"),
            expectError:    true,
        },
        {
            name:        "State change not allowed",
            fromState:   "state1",
            toState:     "state2",
            allowChange: false,
            expectError: true,
        },
        {
            name:          "StateOut failed",
            fromState:     "state1",
            toState:       "state2",
            allowChange:   true,
            stateOutError: errors.New("state out failed"),
            expectError:   true,
        },
        {
            name:         "StateIn failed",
            fromState:    "state1",
            toState:      "state2",
            allowChange:  true,
            stateInError: errors.New("state in failed"),
            expectError:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            state1 := &MockState{
                name:          tt.fromState,
                stateOutError: tt.stateOutError,
            }
            state2 := &MockState{
                name:         tt.toState,
                stateInError: tt.stateInError,
            }
            
            stateMap := StateMap{
                States: map[string]State{
                    tt.fromState: state1,
                    "state2":     state2,
                },
            }
            
            smi := &MockStateMachine{
                allowChange:      tt.allowChange,
                checkChangeError: tt.checkChangeErr,
            }
            
            sm := NewStateMachine(smi, stateMap)
            
            // Start state machine first
            if err := sm.Start(tt.fromState); err != nil {
                t.Fatalf("Failed to start state machine: %v", err)
            }
            
            // Test state change
            err := sm.ChangeState(tt.toState)
            
            if tt.expectError && err == nil {
                t.Error("Expected error but got nil")
            }
            if !tt.expectError && err != nil {
                t.Errorf("Expected no error but got: %v", err)
            }
        })
    }
}

func TestStateMachine_GetCurrentState(t *testing.T) {
    state1 := &MockState{name: "state1"}
    stateMap := StateMap{
        States: map[string]State{
            "state1": state1,
        },
    }
    
    sm := NewStateMachine(&MockStateMachine{}, stateMap)
    
    // Should be nil before start
    if state := sm.GetCurrentState(); state != nil {
        t.Error("Current state should be nil before start")
    }
    
    // Should be state1 after start
    if err := sm.Start("state1"); err != nil {
        t.Fatalf("Start failed: %v", err)
    }
    
    if state := sm.GetCurrentState(); state.GetName() != "state1" {
        t.Errorf("Expected state to be state1, got %s", state.GetName())
    }
}

func TestStateMachine_IsRunning(t *testing.T) {
    state1 := &MockState{name: "state1"}
    stateMap := StateMap{
        States: map[string]State{
            "state1": state1,
        },
    }
    
    sm := NewStateMachine(&MockStateMachine{}, stateMap)
    
    // Should be false before start
    if sm.IsRunning() {
        t.Error("IsRunning should return false before start")
    }
    
    // Should be true after start
    if err := sm.Start("state1"); err != nil {
        t.Fatalf("Start failed: %v", err)
    }
    
    if !sm.IsRunning() {
        t.Error("IsRunning should return true after start")
    }
}

// Test complete workflow
func TestStateMachine_CompleteWorkflow(t *testing.T) {
    // Create states with logging
    menuState := &GameState{name: "menu"}
    playState := &GameState{name: "play"}
    pauseState := &GameState{name: "pause"}
    
    stateMap := StateMap{
        States: map[string]State{
            "menu":  menuState,
            "play":  playState,
            "pause": pauseState,
        },
    }
    
    // Create state machine
    gsm := &GameStateMachine{}
    sm := NewStateMachine(gsm, stateMap)
    
    t.Log("Starting complete workflow test...")
    
    // Start state machine
    if err := sm.Start("menu"); err != nil {
        t.Fatalf("Failed to start state machine: %v", err)
    }
    
    // Verify initial state
    if state := sm.GetCurrentState(); state.GetName() != "menu" {
        t.Errorf("Expected initial state to be menu, got %s", state.GetName())
    }
    
    // Change to play state
    if err := sm.ChangeState("play"); err != nil {
        t.Errorf("Failed to change to play state: %v", err)
    }
    
    // Verify current state
    if state := sm.GetCurrentState(); state.GetName() != "play" {
        t.Errorf("Expected current state to be play, got %s", state.GetName())
    }
    
    // Change to pause state
    if err := sm.ChangeState("pause"); err != nil {
        t.Errorf("Failed to change to pause state: %v", err)
    }
    
    // Verify final state
    if state := sm.GetCurrentState(); state.GetName() != "pause" {
        t.Errorf("Expected final state to be pause, got %s", state.GetName())
    }
    
    // Verify state machine is still running
    if !sm.IsRunning() {
        t.Error("State machine should be running")
    }
} 