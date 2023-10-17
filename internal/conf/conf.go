package conf

type ReveConfig struct {
	Http *HttpConf
}
type HttpConf struct {
	Method string
	Url    string
	Header map[string][]string
}

type Config struct {
	RdsConf     RdsConf
	GptKey      string
	GptKeywords string
}

type RdsConf struct {
	Addr     string
	Password string
	PoolSize int
}
