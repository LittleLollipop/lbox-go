package mission

import (
	"fmt"
	"sync"
	"time"
)

type StepDisposer interface {
	StepName() string
	Dispose(mission *Mission, tags []interface{}) error
}

type Mission struct {
	stepList    []StepDisposer
	stepNow     StepDisposer
	running     bool
	stepRunning bool
	tags        []interface{}
	taskList    []string
	mu          sync.Mutex
}

func NewMission(stepList []StepDisposer) *Mission {
	return &Mission{
		stepList: stepList,
		taskList: make([]string, 0),
	}
}

func (m *Mission) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		fmt.Println("already started")
		return
	}

	m.running = true
	m.stepNow = m.stepList[0]
	m.stepRunning = true
	_ = m.stepNow.Dispose(m, m.tags)
	m.stepRunning = false
	go m.tick()
}

func (m *Mission) GoNext() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.taskList = append(m.taskList, "next")
}

func (m *Mission) tick() {
	m.mu.Lock()
	
	if !m.running {
		m.mu.Unlock()
		return
	}
	
	if m.stepRunning {
		m.mu.Unlock()
		panic("no concurrency allowed")
	}

	if len(m.taskList) == 0 {
		m.mu.Unlock()
		time.Sleep(time.Millisecond)
		go m.tick()
		return
	}

	currentTask := m.taskList[0]
	m.taskList = m.taskList[1:]
	m.stepRunning = true

	if currentTask == "next" {
		next := false
		for _, iterator := range m.stepList {
			if next {
				m.stepNow = iterator
				_ = m.stepNow.Dispose(m, m.tags)
				next = false
				break
			} else {
				if iterator.StepName() == m.stepNow.StepName() {
					next = true
				}
			}
		}

		if next {
			m.running = false
			m.stepRunning = false
			m.mu.Unlock()
			return
		}
	} else {
		for _, iterator := range m.stepList {
			if iterator.StepName() == currentTask {
				m.stepNow = iterator
				_ = m.stepNow.Dispose(m, m.tags)
				break
			}
		}
	}

	m.stepRunning = false
	m.mu.Unlock()
	
	time.Sleep(time.Millisecond)
	go m.tick()
}

func (m *Mission) Jump(stepName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if there's already a non-"next" task in the queue
	for _, task := range m.taskList {
		if task != "next" {
			panic(fmt.Sprintf("error task list: cannot add jump task when another jump task exists: %v", m.taskList))
		}
	}
	
	// Validate if the step exists
	stepExists := false
	for _, iterator := range m.stepList {
		if iterator.StepName() == stepName {
			stepExists = true
			m.taskList = append(m.taskList, stepName)
			break
		}
	}
	
	if !stepExists {
		fmt.Printf("warning: step %s not found\n", stepName)
	}
}
