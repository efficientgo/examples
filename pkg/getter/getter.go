package getter

// Example showing readability difference between using potentially expensive closure once vs many times.
// You can read more in "Efficient Go" book, Chapter 1.

type Report interface {
	Error() error
}

type ReportGetter interface {
	Get() []Report
}

func FailureRatio(reports ReportGetter) float64 {
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

func FailureRatio2(reports ReportGetter) float64 {
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
