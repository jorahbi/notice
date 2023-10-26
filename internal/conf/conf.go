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
	RdsConf RdsConf
	GPT     GPT
	Proxy   string
	Jobs    []Job
}

type Job struct {
	Name string
	Spec string
	To   []string
	Tpl  string
	Rule []string
}

type GPT struct {
	Key      string
	Keywords string
	BaseUrl  string
}

type RdsConf struct {
	Addr     string
	Password string
	PoolSize int
}
