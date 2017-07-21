package interpolator

// Interpolator is an interface that defines a clear boundary for interpolating specific directories in /parts.
// There (by design) may be more than one implementation of this interface, for different use cases in the program.
// The reason we pull this up to an interface is to encourage standardization of our code, and so we can use the Interpolators
// in other packages such as the InterpolatorWriter without having to care of the specific concrete implementation.
type Interpolator interface {

	// Interpolate will interpolate the minimal amount of values necessary into this specific directory.
	Interpolate() error

	// GetTemplate is an Interpolator interface method, and is used by the InterpolatorWriter. This method
	// returns the template []byte data or an error
	GetTemplate() ([]byte, error)

	// GetParameters is an Interpolator interface method, and is used by the InterpolatorWriter. This method
	// returns the parameters []byte data or an error
	GetParameters() ([]byte, error)
}
