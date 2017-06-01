package cfg

import (
	"flag"
	"os"
	"strings"
	"errors"
	"log"
)

/**
 * cli parameters
 */
type Args struct {
	ConfigPath string
}

func GetArgs() *Args {
	args := new(Args)
	h := flag.Bool("h", false, "Get help informations")
	c := flag.String("c", "./src/magpie/conf/magpie.ini", "Magpie's config file path")
	flag.Parse()
	if *h {
		flag.Usage()
		os.Exit(0)
	}

	args.ConfigPath = *c
	err := checkArgs(args)
	if err != nil {
		log.Fatalln("[E] Invalid args: " + *c + err.Error())
	}
	return args
}

func checkArgs(args *Args) error {
	// Check whether config file exists
	if strings.Index(args.ConfigPath, "./") != 0 && strings.Index(args.ConfigPath, "/") != 0 {
		args.ConfigPath = "./" + args.ConfigPath
	}
	_, err := os.Stat(args.ConfigPath)
	if err != nil {
		return errors.New("Config file not found")
	}

	return nil
}

