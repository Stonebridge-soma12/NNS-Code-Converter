package TrainMonitor

type options struct {
	queryString string
	args        []interface{}
	queryId     int // id of selected record
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}