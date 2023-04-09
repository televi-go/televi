package telegram

type ParamsWriter interface {
	WriteParameter(params Params) error
}
