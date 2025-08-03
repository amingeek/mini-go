package qtodo

import (
	"errors"
	"time"
)

type Task interface {
	DoAction()
	GetAlarmTime() time.Time
	GetAction() func()
	GetName() string
	GetDescription() string
}

type basicTask struct {
	name        string
	description string
	alarmTime   time.Time
	action      func()
}

func (t *basicTask) DoAction() {
	if t.action != nil {
		t.action()
	}
}

func (t *basicTask) GetAlarmTime() time.Time {
	return t.alarmTime
}

func (t *basicTask) GetAction() func() {
	return t.action
}

func (t *basicTask) GetName() string {
	return t.name
}

func (t *basicTask) GetDescription() string {
	return t.description
}

func NewTask(action func(), alarmTime time.Time, name string, description string) (Task, error) {
	if name == "" {
		return nil, errors.New("task name cannot be empty")
	}
	if action == nil {
		return nil, errors.New("task action cannot be nil")
	}
	if alarmTime.IsZero() {
		return nil, errors.New("alarm time cannot be zero")
	}
	if alarmTime.Before(time.Now()) {
		return nil, errors.New("alarm time cannot be in the past")
	}
	if description == "" {
		return nil, errors.New("task description cannot be empty")
	}

	return &basicTask{
		name:        name,
		description: description,
		alarmTime:   alarmTime,
		action:      action,
	}, nil
}
