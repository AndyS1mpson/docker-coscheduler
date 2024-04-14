package slices

// ToMap создает мапу с ключами из элементов списка.
func ToMap[T comparable](s []T) map[T]struct{} {
	unique := make(map[T]struct{})

	for _, elem := range s {
		unique[elem] = struct{}{}
	}

	return unique
}
