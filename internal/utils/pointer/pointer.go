package pointer

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Get сахар ValueOrDefault для совместимости с AlekSi/pointer
func Get[A any](p *A) A {
	return ValueOrDefault(p)
}

// ValueOr возвращает значение p, если p != nil, иначе возвращает or
func ValueOr[A any](p *A, or A) A {
	if p == nil {
		return or
	}

	return *p
}

// ValueOrDefault возвращает значение p, если p != nil, иначе возвращает дефолтное значение типа
func ValueOrDefault[A any](p *A) A {
	var def A

	return ValueOr(p, def)
}

// CastString преобразует существующее текстовое значение к типу A и возвращает указатель на результат
func CastString[A ~string](v *string) *A {
	if v == nil {
		return nil
	}

	return To(A(*v))
}

// Map применяет функцию mapper к существующему значению и возвращает указатель на результат
func Map[A any, B any](v *A, mapper func(A) B) *B {
	if v == nil {
		return nil
	}

	return To(mapper(*v))
}

// To возвращает ссылку на value
func To[A any](value A) *A {
	return &value
}

// ToInt64 возвращает ссылку на value
func ToInt64[A number](value A) *int64 {
	return To(int64(value))
}

// ToString возвращает ссылку на value
func ToString[A ~string](value A) *string {
	return To(string(value))
}

// NilIfDefault возвращает nil, если value является дефолтным значением типа, иначе - ссылку на value
func NilIfDefault[A comparable](value A) *A {
	var def A

	return NilIf(value, func(v A) bool { return v == def })
}

// NilIf возвращает nil, если результат cond == true, иначе - ссылку на value
func NilIf[A any](value A, cond func(A) bool) *A {
	if cond(value) {
		return nil
	}

	return &value
}

// NilIfEmpty возвращает nil, если слайс пуст, иначе - указатель на слайс
func NilIfEmpty[A any](s []A) *[]A {
	if len(s) == 0 {
		return nil
	}

	return &s
}
