package qtodo

import (
	"fmt"
	"sync"
)

type Database interface {
	GetTaskList() []Task
	GetTask(name string) (Task, error)
	SaveTask(t Task) error
	DelTask(name string) error
}
type memoryDB struct {
	mu   sync.RWMutex
	data map[string]Task
}

func NewDatabase() Database {
	return &memoryDB{
		data: make(map[string]Task),
	}
}
func (db *memoryDB) GetTaskList() []Task {
	db.mu.RLock()
	defer db.mu.RUnlock()
	list := make([]Task, 0, len(db.data))
	for _, t := range db.data {
		list = append(list, t)
	}
	return list
}

func (db *memoryDB) GetTask(name string) (Task, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if t, ok := db.data[name]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("task %q not found", name)
}

func (db *memoryDB) SaveTask(t Task) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[t.GetName()] = t
	return nil
}

func (db *memoryDB) DelTask(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.data[name]; !ok {
		return fmt.Errorf("task %q not found", name)
	}
	delete(db.data, name)
	return nil
}
