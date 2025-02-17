package config

import (
	"errors"
	"fmt"
	"sync"
)

type handler[T any] struct {
	mu      sync.Mutex
	devices []Device[T]
}

// Devices returns the handlers devices
func (cp *handler[T]) Devices() []Device[T] {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	return cp.devices
}

// Add adds device and config
func (cp *handler[T]) Add(dev Device[T]) error {
	conf := dev.Config()

	if conf.Name == "" {
		return errors.New("missing name")
	}

	if _, err := cp.ByName(conf.Name); err == nil {
		return fmt.Errorf("duplicate name: %s already defined and must be unique", conf.Name)
	}

	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.devices = append(cp.devices, dev)

	return nil
}

// Delete deletes device
func (cp *handler[T]) Delete(name string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for i, dev := range cp.devices {
		if name == dev.Config().Name {
			cp.devices = append(cp.devices[:i], cp.devices[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("not found: %s", name)
}

// ByName provides device by name
func (cp *handler[T]) ByName(name string) (Device[T], error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for _, dev := range cp.devices {
		if name == dev.Config().Name {
			return dev, nil
		}
	}

	return nil, fmt.Errorf("not found: %s", name)
}
