package assets

// AssetManager provides an interface to retrieve static assets.
type AssetManager interface {
	// Asset loads and returns the asset for the given name.
	// It returns an error if the asset could not be found or
	// could not be loaded.
	Asset(name string) ([]byte, error)

	// AssetNames returns the names of the assets.
	AssetNames() []string
}
