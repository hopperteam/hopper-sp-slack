package state

import (
	"sp-slack/db"
)

const (
	appStrKey = "appStr"
)

func GetAppStr() string {
	return db.SelectState(appStrKey)
}

func SetAppStr(value string) {
	db.UpsertState(appStrKey, value)
}