package demo

type Singleton interface {
	Increment()
	Get() int
}

var instance singletonImpl

func Instance() Singleton {
	return &instance
}

type singletonImpl struct {
	value int
}

func (s *singletonImpl) Increment() {
	s.value++
}

func (s *singletonImpl) Get() int {
	return s.value
}
