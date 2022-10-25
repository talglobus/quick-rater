package input

type Key interface {
	IsEscape() bool
	IsBackspace() bool
	GetRating() (int, bool)
	GetBool() (bool, bool)
	GetNumeric() (int, bool)
}
