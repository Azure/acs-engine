package interpolator

type Interpolator interface {
	Interpolate() error
	GetTemplate() ([]byte, error)
	GetParameters() ([]byte, error)
}