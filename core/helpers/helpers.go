package helpers

func MapForEachWithErr[INP, OUT any](input []INP, mapper func(INP) (OUT, error)) (output []OUT, err error) {
	for _, elem := range input {
		mapped, err := mapper(elem)
		if err != nil {
			return []OUT{}, err
		}
		output = append(output, mapped)
	}
	return
}

func MapForEach[INP, OUT any](input []INP, mapper func(INP) OUT) (output []OUT) {
	for _, elem := range input {
		output = append(output, mapper(elem))
	}
	return
}
