package dto

type SumMessage struct {
	Arg1   *int `msgpack:"arg1"`
	Arg2   *int `msgpack:"arg2"`
	Result *int `msgpack:"result"`
}

type SumResponse struct {
	Arg1   int `msgpack:"arg1"`
	Arg2   int `msgpack:"arg2"`
	Result int `msgpack:"result"`
}

func (s *SumMessage) Validate() bool {
	return s.Arg1 != nil && s.Arg2 != nil && s.Result != nil
}
