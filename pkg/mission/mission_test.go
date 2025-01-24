package mission

import (
	"sync"
	"testing"
	"time"
)

// Mock step implementation
type MockStep struct {
	name      string
	executed  bool
	mu        sync.Mutex
	onDispose func()
}

func NewMockStep(name string) *MockStep {
	return &MockStep{
		name: name,
	}
}

func (s *MockStep) StepName() string {
	return s.name
}

func (s *MockStep) Dispose(_ *Mission, _ []interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executed = true
	if s.onDispose != nil {
		s.onDispose()
	}
	return nil
}

func (s *MockStep) WasExecuted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.executed
}

func TestMission_Start(t *testing.T) {
	step1 := NewMockStep("step1")
	step2 := NewMockStep("step2")
	
	mission := NewMission([]StepDisposer{step1, step2})
	
	mission.Start()
	time.Sleep(10 * time.Millisecond) // Wait for async execution
	
	if !step1.WasExecuted() {
		t.Error("First step should be executed")
	}
	if step2.WasExecuted() {
		t.Error("Second step should not be executed")
	}
}

func TestMission_GoNext(t *testing.T) {
	step1 := NewMockStep("step1")
	step2 := NewMockStep("step2")
	step3 := NewMockStep("step3")
	
	mission := NewMission([]StepDisposer{step1, step2, step3})
	
	mission.Start()
	time.Sleep(10 * time.Millisecond)
	
	mission.GoNext()
	time.Sleep(10 * time.Millisecond)
	
	if !step2.WasExecuted() {
		t.Error("Second step should be executed")
	}
	
	mission.GoNext()
	time.Sleep(10 * time.Millisecond)
	
	if !step3.WasExecuted() {
		t.Error("Third step should be executed")
	}
}

func TestMission_Jump(t *testing.T) {
	step1 := NewMockStep("step1")
	step2 := NewMockStep("step2")
	step3 := NewMockStep("step3")
	
	mission := NewMission([]StepDisposer{step1, step2, step3})
	
	mission.Start()
	time.Sleep(10 * time.Millisecond)
	
	mission.Jump("step3")
	time.Sleep(10 * time.Millisecond)
	
	if !step3.WasExecuted() {
		t.Error("Should jump to step3")
	}
}

func TestMission_ConcurrencyControl(t *testing.T) {
	step1 := NewMockStep("step1")
	executionCount := 0
	
	step1.onDispose = func() {
		time.Sleep(10 * time.Millisecond) // Reduced wait time
		executionCount++
	}
	
	mission := NewMission([]StepDisposer{step1})
	
	mission.Start()
	time.Sleep(5 * time.Millisecond) // Wait for first execution to start
	mission.Start() // Try to start again
	
	time.Sleep(50 * time.Millisecond) // Ensure enough time for completion
	
	if executionCount != 1 {
		t.Errorf("Step1 should be executed only once, but was executed %d times", executionCount)
	}
}

func TestMission_TaskQueueValidation(t *testing.T) {
	step1 := NewMockStep("step1")
	step2 := NewMockStep("step2")
	
	mission := NewMission([]StepDisposer{step1, step2})
	
	// Use channel to ensure we catch the panic
	done := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- true
				return
			}
			done <- false
		}()
		
		mission.Start()
		time.Sleep(20 * time.Millisecond) // Increased wait time
		
		mission.Jump("step1")
		mission.Jump("step2") // This should trigger panic immediately
	}()
	
	// Wait for goroutine completion
	select {
	case didPanic := <-done:
		if !didPanic {
			t.Error("Should panic when executing multiple jumps")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Test timeout")
	}
}

func TestMission_CompletionState(t *testing.T) {
	step1 := NewMockStep("step1")
	step2 := NewMockStep("step2")
	
	mission := NewMission([]StepDisposer{step1, step2})
	
	mission.Start()
	time.Sleep(10 * time.Millisecond)
	
	mission.GoNext()
	time.Sleep(10 * time.Millisecond)
	mission.GoNext() // Try to move to non-existent next step
	time.Sleep(10 * time.Millisecond)
	
	if mission.running {
		t.Error("Mission should stop running after completing all steps")
	}
} 