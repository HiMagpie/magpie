package cfg

var (
	C *Conf
)

func init() {
	ParseArgs()
}

func ParseArgs() {
	args := GetArgs()
	C = NewConf()
	parseConfigFile(args.ConfigPath)
}

