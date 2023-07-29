package demo

import (
	"fmt"
)

type ModuleInterface interface {
	Start()
	IsRunning() bool
	Stop()
}

type MainModule struct {
	dep *SecondModule
}

func newMainModule(dep *SecondModule) *MainModule {
	return &MainModule{
		dep: dep,
	}
}

func (m *MainModule) Start() {
	m.dep.Start()
	fmt.Println("MainModule started")
}

func (m *MainModule) IsRunning() bool {
	return m.dep.IsRunning()
}

func (m *MainModule) Stop() {
	m.dep.Stop()
	fmt.Println("MainModule stopped")
}

type SecondModule struct {
	dep *ThirdModule
}

func newSecondModule(dep *ThirdModule) *SecondModule {
	return &SecondModule{
		dep: dep,
	}
}

func (m *SecondModule) Start() {
	m.dep.Start()
	fmt.Println("SecondModule started")
}

func (m *SecondModule) IsRunning() bool {
	return m.dep.IsRunning()
}

func (m *SecondModule) Stop() {
	m.dep.Stop()
	fmt.Println("SecondModule stopped")
}

type ThirdModule struct {
	running bool
}

func newThirdModule() *ThirdModule {
	return &ThirdModule{}
}

func (m *ThirdModule) Start() {
	fmt.Println("ThirdModule started")
	m.running = true
}

func (m *ThirdModule) IsRunning() bool {
	return m.running
}

func (m *ThirdModule) Stop() {
	fmt.Println("ThirdModule stopped")
	m.running = false
}

type ModuleHandler struct {
	main *MainModule
}

func (m *ModuleHandler) Bootstrap() {
	third := newThirdModule()
	second := newSecondModule(third)
	m.main = newMainModule(second)
	m.main.Start()
}

func (m *ModuleHandler) IsRunning() bool {
	return m.main.IsRunning()
}

func (m *ModuleHandler) Teardown() {
	m.main.Stop()
	m.main.dep.Stop()
	m.main.dep.dep.Stop()
}
