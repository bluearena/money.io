package assets

// AssetManagerMock is a mocked implementation of the AssetManager.
type AssetManagerMock struct {
	OnAsset      func(name string) ([]byte, error)
	OnAssetNames func() []string
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func (am AssetManagerMock) Asset(name string) ([]byte, error) {
	return am.OnAsset(name)
}

// AssetNames returns the names of the assets.
func (am AssetManagerMock) AssetNames() []string {
	return am.OnAssetNames()
}
