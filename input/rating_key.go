package input

type ratingKey struct {
	value int
}

func (_ ratingKey) IsEscape() bool {
	return false
}

func (_ ratingKey) IsBackspace() bool {
	return false
}

func (k ratingKey) GetRating() (int, bool) {
	return k.value, true
}

func (_ ratingKey) GetBool() (bool, bool) {
	return false, false
}

func (k ratingKey) GetNumeric() (int, bool) {
	return k.value, true
}
