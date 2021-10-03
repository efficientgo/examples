package dup

func DeduplicateNaively(input []int) (output []int) {
	for _, i := range input {
		dup := false
		for _, o := range output {
			if i == o {
				dup = true
			}
		}
		if !dup {
			output = append(output, i)
		}
	}
	return output
}

func DeduplicateFaster(input []int) (output []int) {
	unique := make(map[int]struct{})
	for _, i := range input {
		if _, ok := unique[i]; ok {
			continue
		}
		unique[i] = struct{}{}
		output = append(output, i)
	}

	return output
}

func DeduplicateLessAllocs(input []int) []int {
	o := 0
	unique := make(map[int]struct{}, len(input))
	for _, i := range input {
		if _, ok := unique[i]; ok {
			continue
		}
		unique[i] = struct{}{}
		input[o] = i
		o++
	}

	return input[:o]
}

func DeduplicateLessAllocs2(input []int) []int {
	unique := make(map[int]struct{}, len(input))
	for _, i := range input {
		// Quite funny trick where we always loop through all (only write), then only read - good for cache locality.
		unique[i] = struct{}{}
	}

	o := 0
	for u := range unique {
		input[o] = u
		o++
	}
	return input[:o]
}

func DeduplicateDynamic(input []int) (output []int) {
	const switchThreshold = 10

	var unique map[int]struct{}
	for _, i := range input {
		if len(output) < switchThreshold {
			// Fast path for almost all-duplication cases.
			dup := false
			for _, o := range output {
				if i == o {
					dup = true
				}
			}
			if dup {
				continue
			}
		} else {
			// More efficient path if we have more unique elements.
			if _, ok := unique[i]; ok {
				continue
			}
			unique[i] = struct{}{}
		}

		output = append(output, i)
		if len(output) == switchThreshold {
			unique = make(map[int]struct{})
			for _, o := range output {
				unique[o] = struct{}{}
			}
		}
	}
	return output
}

func DeduplicateDynamicLessAllocs(input []int) []int {
	const switchThreshold = 10

	o := 0
	var unique map[int]struct{}
	for _, i := range input {
		if o < switchThreshold {
			// Fast path for almost all-duplication cases.
			dup := false
			for j := 0; j < o; j++ {
				if i == input[j] {
					dup = true
				}
			}
			if dup {
				continue
			}
		} else {
			// More efficient path if we have more unique elements.
			if _, ok := unique[i]; ok {
				continue
			}
			unique[i] = struct{}{}
		}

		input[o] = i
		o++
		if o == switchThreshold {
			unique = make(map[int]struct{}, len(input))
			for j := 0; j < o; j++ {
				unique[input[j]] = struct{}{}
			}
		}
	}
	return input[:o]
}

func DeduplicateDynamicLessAllocs2(input []int) []int {
	const switchThreshold = 10

	o := 0
	var unique map[int]struct{}
	for _, i := range input {
		if o < switchThreshold {
			// Fast path for almost all-duplication cases.
			dup := false
			for j := 0; j < o; j++ {
				if i == input[j] {
					dup = true
				}
			}
			if dup {
				continue
			}
			input[o] = i
			o++
			if o == switchThreshold {
				unique = make(map[int]struct{}, len(input))
				for j := 0; j < o; j++ {
					unique[input[j]] = struct{}{}
				}
			}
			continue

		}
		// More efficient path if we have more unique elements.
		unique[i] = struct{}{}
	}

	if unique != nil {
		o = 0
		for u := range unique {
			input[o] = u
			o++
		}
	}
	return input[:o]
}
