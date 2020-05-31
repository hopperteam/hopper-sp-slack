package config

import (
	"fmt"
	"os"
	"sp-slack/logger"
)

var Port string 
var Signature string

func Init() error {
    var err error
    Port = os.Getenv("PORT")
    Signature = os.Getenv("SIGNATURE")

    logger.Infof("using %s and %s", Port, Signature)
    if Port == "" || Signature == "" {
        return fmt.Errorf("config incomplete")
    }

    return err
}

