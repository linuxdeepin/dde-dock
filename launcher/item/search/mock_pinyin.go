package search

type MockPinYin struct {
	data  map[string][]string
	valid bool
}

func (self *MockPinYin) Search(key string) ([]string, error) {
	return self.data[key], nil
}

func (self *MockPinYin) IsValid() bool {
	return self.valid
}

func NewMockPinYin(data map[string][]string, valid bool) *MockPinYin {
	return &MockPinYin{
		data:  data,
		valid: valid,
	}
}
