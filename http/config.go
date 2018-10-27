package http

type Config struct {
	PORT         int64    `toml:"port"`
	ALLOW_ORIGIN []string `toml:"allow_origin"`
	API_PREFIX   string   `toml:"api_prefix"`
	VIEW_DIR     string   `toml:"view_dir"`
}
