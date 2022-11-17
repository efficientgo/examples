package getter

type Report interface {
	Error() error
}

type ReportGetter interface {
	Get() []Report
}

// FailureRatio represents slightly less efficient (also less safe) use of getters.
// Read more in "Efficient Go"; Example 1-1.
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

// FailureRatio_Better represents more efficient (also safer and more readable) use of getters.
// Read more in "Efficient Go"; Example 1-2 (called `FailureRatio` in example).
func FailureRatio_Better(reports ReportGetter) float64 {
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
