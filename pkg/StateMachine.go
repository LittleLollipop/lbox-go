package lbox

import (
	"fmt"
	"sync"
)

type StateMachineInterface interface {
	InitData() error
	CheckStateChange(stateNow *State, newState *State) (bool, error)
}

type State interface {
	Name() string
	StateOut() bool
	StateIn() bool
}

type StateMap struct {
	States map[string]*State
}

type StateMachine struct {
	smi      StateMachineInterface
	stateMap StateMap
	initing  bool
	running  bool
	stateNow *State
	mu       sync.Mutex
}

func NewStateMachine(smi StateMachineInterface, stateMap StateMap) *StateMachine {
	return &StateMachine{
		smi:      smi,
		stateMap: stateMap,
	}
}

func (sm *StateMachine) Start(firstState string) {
	if sm.initing || sm.running {
		fmt.Println("already start")
		return
	}
	if sm.stateMap.States[firstState] == nil {
		panic(fmt.Errorf("state not found: %s", firstState))
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.initing = true

	sm.smi.InitData()
	sm.stateNow = sm.stateMap.States[firstState]
	(*sm.stateNow).StateIn()

	sm.running = true
	sm.initing = false
}

func (sm *StateMachine) ChangeState(stateName string) {
	if sm.stateMap.States[stateName] == nil {
		panic(fmt.Errorf("state not found: %s", stateName))
	}
	sm.doChangeState(stateName)
}

func (sm *StateMachine) doChangeState(stateName string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	newState := sm.stateMap.States[stateName]

	stateChanged, err := sm.smi.CheckStateChange(sm.stateNow, newState)
	if err != nil {
		panic(fmt.Errorf("error checking state change: %v", err))
	}

	if stateChanged {
		outStatus := (*sm.stateNow).StateOut()

		if !outStatus {
			errMsg := fmt.Sprintf("stateOut failed: %s, next state: %s", (*sm.stateNow).Name(), stateName)
			panic(errMsg)
		}

		sm.stateNow = newState
		inStatus := (*sm.stateNow).StateIn()

		if !inStatus {
			errMsg := fmt.Sprintf("stateIn failed: %s", (*sm.stateNow).Name())
			panic(errMsg)
		}
	}
}
