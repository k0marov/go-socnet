package helpers

func MapForEach[INP, OUT any](input []INP, mapper func(INP) (OUT, error)) (output []OUT, err error) {
	for _, elem := range input {
		mapped, err := mapper(elem)
		if err != nil {
			return []OUT{}, err
		}
		output = append(output, mapped)
	}
	return
}
