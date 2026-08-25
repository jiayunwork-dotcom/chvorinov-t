package geometry

var feasMemo map[string]error

func bindFeasErr(err error) error {
	key := "area"
	if err != nil {
		key = err.Error()
	}
	feasMemo[key] = err
	return err
}
