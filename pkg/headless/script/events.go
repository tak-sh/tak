package script

type EventQueue chan Event

func NewEventQueue() EventQueue {
	return make(EventQueue, 10)
}

type Event interface {
	eventSigil()
}

var _ Event = &ChangeStepEvent{}

type ChangeStepEvent struct {
	// The step that is now running.
	Step *Step
	// The index of the step.
	Idx int
	// Total number of steps.
	Total int
}

func (c *ChangeStepEvent) eventSigil() {}
