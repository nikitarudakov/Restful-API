package repoerr

type CalcRatingUserError struct {
	Msg string
}

func (m *CalcRatingUserError) Error() string {
	return m.Msg
}
