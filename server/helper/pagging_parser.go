package helper

func PageParser(page int, limit int) (p int, l int, err error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 0
	}

	if page == 0 {
		p = 0
	} else {
		p = (page - 1) * limit
	}
	l = limit

	return
}
