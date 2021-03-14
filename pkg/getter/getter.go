package getter

// Example showing readability difference between using potentially expensive closure once vs many times.

type Report interface {
	Error() error
}

type Reporter interface {
	Get() []Report
}

func FailureRatio(reports Reporter) float64 {
	if len(reports.Get()) == 0 {
		return 0
	}

	var sum float64
	for _, report := range reports.Get() {
		if report.Error() != nil {
			sum++
		}
	}
	return sum / float64(len(reports.Get()))
}

func FailureRatio2(reports Reporter) float64 {
	got := reports.Get()
	if len(got) == 0 {
		return 0
	}

	var sum float64
	for _, report := range got {
		if report.Error() != nil {
			sum++
		}
	}
	return sum / float64(len(got))
}
