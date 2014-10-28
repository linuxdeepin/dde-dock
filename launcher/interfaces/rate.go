package interfaces

type RateConfigFileInterface interface {
	Free()
	SetUint64(string, string, uint64)
	GetUint64(string, string) (uint64, error)
	ToData() (uint64, string, error)
}
