package dock

type positionType int32

const (
	positionTop positionType = iota
	positionRight
	positionBottom
	positionLeft
)

func (p positionType) String() string {
	switch p {
	case positionTop:
		return "Top"
	case positionRight:
		return "Right"
	case positionBottom:
		return "Bottom"
	case positionLeft:
		return "Left"
	default:
		return "Unknown"
	}
}
