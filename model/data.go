package model

import "time"

type Data struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	ExpireTime  time.Time
	IsTimeGiven bool
}
