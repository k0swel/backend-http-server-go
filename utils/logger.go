package utils

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "[k0swel web_application] ", log.LstdFlags|log.Lshortfile)
