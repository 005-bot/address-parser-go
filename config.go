package address

// Config controls parser initialization. DBPath optionally points to an
// external SQLite streets database; when empty, the embedded copy is used.
type Config struct {
	DBPath string `koanf:"db_path"`
}
