package input

type boolKey struct {
	value bool
}

func (_ boolKey) IsEscape() bool {
	return false
}

func (_ boolKey) IsBackspace() bool {
	return false
}

func (_ boolKey) GetRating() (int, bool) {
	return 0, false
}

func (k boolKey) GetBool() (bool, bool) {
	return k.value, true
}

func (k boolKey) GetNumeric() (int, bool) {
	if k.value {
		return 5, true
	} else {
		return 0, true
	}
}
