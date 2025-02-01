package choose

func If[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func IfLazyL[T any](condition bool, trueValue func() T, falseValue T) T {
	if condition {
		return trueValue()
	}
	return falseValue
}

func IfLazyR[T any](condition bool, trueValue T, falseValue func() T) T {
	if condition {
		return trueValue
	}
	return falseValue()
}
