package Logger

type LogWriterInterface interface {
	Write(prefix string, data interface{}) error
}
