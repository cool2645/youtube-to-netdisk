package carrier

type Config struct {
	DB_NAME      string `toml:"db_name"`
	DB_USER      string `toml:"db_user"`
	DB_PASS      string `toml:"db_pass"`
	DB_CHARSET   string `toml:"db_charset"`
	DB_COLLATION string `toml:"db_collation"`
	TEMP_PATH    string `toml:"temp_path"`
}

func parseDSN(config Config) string {
	return config.DB_USER + ":" + config.DB_PASS + "@/" + config.DB_NAME + "?charset=" + config.DB_CHARSET + "&collation=" + config.DB_COLLATION + "&parseTime=true&loc=Local"
}
