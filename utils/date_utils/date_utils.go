package date_utils

import "time"

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDbLayout = "2006-01-02 15:04:05"
	)

func GetNowString() string {
	now := time.Now().UTC()
	return now.Format(apiDateLayout)
}

func GetNowDbFormat() string {
	now := time.Now().UTC()
	return now.Format(apiDbLayout)
}