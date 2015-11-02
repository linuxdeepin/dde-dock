package mock

type PinYin struct {
	data  map[string][]string
	valid bool
}

func (self *PinYin) Search(key string) ([]string, error) {
	return self.data[key], nil
}

func (self *PinYin) IsValid() bool {
	return self.valid
}

func NewPinYin(data map[string][]string, valid bool) *PinYin {
	return &PinYin{
		data:  data,
		valid: valid,
	}
}
