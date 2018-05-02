package config

var GlobCfg = Config{}

type Config struct {
	PORT         int64    `toml:"port"`
	ALLOW_ORIGIN []string `toml:"allow_origin"`
	WEB_URL      string   `toml:"web_url"`
	TG_ENABLE    bool     `toml:"tg_enable"`
	RIRI_ADDR    string   `toml:"riri_addr"`
	RIRI_KEY     string   `toml:"riri_key"`
	DB_NAME      string   `toml:"db_name"`
	DB_USER      string   `toml:"db_user"`
	DB_PASS      string   `toml:"db_pass"`
	DB_CHARSET   string   `toml:"db_charset"`
	DB_COLLATION string   `toml:"db_collation"`
	PYTHON_CMD   string   `toml:"python_cmd"`
	TEMP_PATH    string   `toml:"temp_path"`
	ND_FOLDER    string   `toml:"netdisk_folder"`
	ND_SHARELINK string   `toml:"netdisk_sharelink"`
	ND_SHAREPASS string   `toml:"netdisk_sharepass"`
}

func ParseDSN(config Config) string {
	return config.DB_USER + ":" + config.DB_PASS + "@/" + config.DB_NAME + "?charset=" + config.DB_CHARSET + "&collation=" + config.DB_COLLATION + "&parseTime=true&loc=Local"
}
