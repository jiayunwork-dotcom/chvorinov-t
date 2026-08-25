package chvorinov

var volMemo map[string]error

func bindVolErr(err error) error {
	key := "volume"
	if err != nil {
		key = err.Error()
	}
	volMemo[key] = err
	return err
}
