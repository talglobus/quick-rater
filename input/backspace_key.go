package input

type backspaceKey struct{}

func (_ backspaceKey) IsEscape() bool {
	return false
}

func (_ backspaceKey) IsBackspace() bool {
	return true
}

func (_ backspaceKey) GetRating() (int, bool) {
	return 0, false
}

func (_ backspaceKey) GetBool() (bool, bool) {
	return false, false
}

func (_ backspaceKey) GetNumeric() (int, bool) {
	return 0, false
}
