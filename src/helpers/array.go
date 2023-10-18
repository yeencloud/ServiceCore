package helpers

func ArrayOrNil[T any](a []T) []T {
	if len(a) > 0 {
		return a
	}
	return nil
}
