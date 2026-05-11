package dto

type PingMessage struct {
	Hello *string `msgpack:"hello"`
	World *string `msgpack:"world"`
	Arg1  *int    `msgpack:"arg1"`
	Arg2  *string `msgpack:"arg2"`
}

func (p *PingMessage) Validate() bool {
	return p.Hello != nil && p.World != nil
}
