package stuff

import (
	"learning/stuff/demo"
	"testing"

	"github.com/stretchr/testify/require"
)

type DependencyInterface interface {
	Increment()
	Get() int
}

type Dependency struct {
	value int
}

func (d *Dependency) Increment() {
	d.value++
}

func (d *Dependency) Get() int {
	return d.value
}

func functionUnderTest(dep DependencyInterface) int {
	dep.Increment()
	return dep.Get()
}

type MockDependency struct {
	Dependency
	incrementCalls int
	getCalls       int
}

func (d *MockDependency) Increment() {
	d.incrementCalls++
	d.Dependency.Increment()
}

func (d *MockDependency) Get() int {
	d.getCalls++
	return d.Dependency.Get()
}

func TestWithMock(t *testing.T) {
	dep := &MockDependency{}
	result := functionUnderTest(dep)
	require.Equal(t, 1, result)
	require.Equal(t, 1, dep.incrementCalls)
	require.Equal(t, 1, dep.getCalls)
}

func TestSingleton(t *testing.T) {
	initial := demo.Instance().Get()
	demo.Instance().Increment()
	after := demo.Instance().Get()
	require.Greater(t, after, initial)
}

func TestModule(t *testing.T) {
	moduleMgr := &demo.ModuleHandler{}
	moduleMgr.Bootstrap()
	require.True(t, moduleMgr.IsRunning())
	moduleMgr.Teardown()
	require.False(t, moduleMgr.IsRunning())
}

type TestListenerOne struct {
	notifyCalls int
}

func (l *TestListenerOne) Notify(message string) {
	l.notifyCalls++
}

type TestListenerTwo struct {
	notifyCalls int
}

func (l *TestListenerTwo) Notify(message string) {
	l.notifyCalls++
}

func TestMessenger(t *testing.T) {
	messenger := &demo.Messenger{}
	listenerOne := &TestListenerOne{}
	listenerTwo := &TestListenerTwo{}
	messenger.Register(listenerOne)
	messenger.Register(listenerTwo)
	messenger.Send("test")
	require.Equal(t, 1, listenerOne.notifyCalls)
	require.Equal(t, 1, listenerTwo.notifyCalls)
}
