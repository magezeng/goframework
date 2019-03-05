package Config

type LoaderInterface interface {
	Load(configCls interface{}, configPaths ...string) error
}
