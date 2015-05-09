package fakes

type FakeStore struct {
	DeleteCalled bool
	DeleteErr    error

	GetCalled bool
	GetFound  bool
	GetValue  string
	GetErr    error

	SaveCalled bool
	SaveErr    error
}

func (s *FakeStore) Delete(key string) error {
	s.DeleteCalled = true
	return s.DeleteErr
}

func (s *FakeStore) Get(key string) (string, bool, error) {
	s.GetCalled = true
	return s.GetValue, s.GetFound, s.GetErr
}

func (s *FakeStore) Save(key, value string) error {
	s.SaveCalled = true
	return s.SaveErr
}
