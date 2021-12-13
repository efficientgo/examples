package dup

import "sort"

type Element struct {
	Something int
}

func needsChange(Element) bool      { return false }
func needsToBeRemoved(Element) bool { return false }
func isTypeA(Element) bool          { return false }
func isTypeB(Element) bool          { return false }
func isX(Element) bool              { return false }

var (
	typeA, typeB, typeUnknown int
	thingsToChange            []Element
	thingsToRemove            []Element
)

func ProcessInput(input []Element) {
	for _, i := range input {
		if needsChange(i) {
			thingsToChange = append(thingsToChange, i)
		}
		if needsToBeRemoved(i) {
			thingsToRemove = append(thingsToRemove, i)
		}
	}

	for _, i := range input {
		if isTypeA(i) {
			typeA++
			continue
		}
		if isTypeB(i) {
			typeB++
			continue
		}
		typeUnknown++
	}

	var x *Element
	for _, i := range input {
		if isX(i) {
			x = &i
			break
		}
	}

	sort.Slice(input, func(i, j int) bool {
		return input[i].Something < input[j].Something
	})

	// ...use processed information
}
