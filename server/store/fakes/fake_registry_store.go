package fakes

type FakeRegistryStore struct {
	DeleteCalled bool
	DeleteErr    error

	GetCalled bool
	GetFound  bool
	GetValue  string
	GetErr    error

	SaveCalled bool
	SaveErr    error
}

func (s *FakeRegistryStore) Delete(key string) error {
	s.DeleteCalled = true
	return s.DeleteErr
}

func (s *FakeRegistryStore) Get(key string) (string, bool, error) {
	s.GetCalled = true
	return s.GetValue, s.GetFound, s.GetErr
}

func (s *FakeRegistryStore) Save(key, value string) error {
	s.SaveCalled = true
	return s.SaveErr
}
