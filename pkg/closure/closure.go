package closure

import "errors"

// Example showing readability difference between using potentially expensive closure once vs many times.

type Report interface {
	Error() error
}

func New(getReports func() []Report) Reporter { return Reporter{getReports: getReports} }

type Reporter struct {
	getReports func() []Report
}

func (r *Reporter) FailureRatio() (float64, error) {
	if len(r.getReports()) == 0 {
		return 0, errors.New("operation still pending")
	}

	var sum float64
	for _, report := range r.getReports() {
		if report.Error() != nil {
			sum++
		}
	}
	return sum / float64(len(r.getReports())), nil
}
