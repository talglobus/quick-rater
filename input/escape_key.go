package input

type escapeKey struct{}

func (_ escapeKey) IsEscape() bool {
	return true
}

func (_ escapeKey) IsBackspace() bool {
	return false
}

func (_ escapeKey) GetRating() (int, bool) {
	return 0, false
}

func (_ escapeKey) GetBool() (bool, bool) {
	return false, false
}

func (_ escapeKey) GetNumeric() (int, bool) {
	return 0, false
}
