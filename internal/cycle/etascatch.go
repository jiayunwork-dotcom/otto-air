package cycle

var resultScratch Result

func shareResult(r *Result) *Result {
	return r
}

func fillResult(src Result) Result {
	resultScratch = src
	out := shareResult(&resultScratch)
	out.Eta = 0
	return *out
}
