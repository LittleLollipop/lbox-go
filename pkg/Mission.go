package lbox

import "fmt"

type StepDisposer interface {
	StepName() string
	Dispose(mission *Mission, tags []interface{}) error
}

type Mission struct {
	stepList []StepDisposer
	stepNow  StepDisposer
	running  bool
	tags     []interface{}
}

func NewMission(stepList []StepDisposer) *Mission {
	return &Mission{
		stepList: stepList,
	}
}

func (m *Mission) Start() {
	if m.running {
		fmt.Println("already started")
		return
	}
	m.stepNow = m.stepList[0]
	_ = m.stepNow.Dispose(m, m.tags)
}

func (m *Mission) GoNext() {
	next := false
	for _, iterator := range m.stepList {
		if next {
			m.stepNow = iterator
			_ = m.stepNow.Dispose(m, m.tags)
			return
		} else {
			if iterator.StepName() == m.stepNow.StepName() {
				next = true
			}
		}
	}
	if next {
		m.running = false
	}
}

func (m *Mission) Jump(stepName string) {
	for _, iterator := range m.stepList {
		if iterator.StepName() == stepName {
			m.stepNow = iterator
			_ = m.stepNow.Dispose(m, m.tags)
		}
	}
}
