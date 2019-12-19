package service_trigger

func Tr(in string) string {
	return in
}

var _ = Tr("\"%s\" did not pass the system security verification, and cannot run now")
