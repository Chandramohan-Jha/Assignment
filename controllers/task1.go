package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/chandramohan/in_memory/db"
	"github.com/chandramohan/in_memory/model"
)

type data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type response struct {
	Status string `json:"status"`
	Data   data   `json:"data"`
}

func isFound(key string) (bool, *model.Data) {
	for _, v := range db.Datas {
		if v.Key == key && (!v.IsTimeGiven || time.Now().Before(v.ExpireTime)) {
			return true, &v
		}
	}
	return false, nil
}

func update(cmd []string, w http.ResponseWriter) {
	key := cmd[1]
	val := cmd[2]
	for i, v := range db.Datas {
		if v.Key == key && (!v.IsTimeGiven || time.Now().Before(v.ExpireTime)) {
			db.Datas[i].Value = val
			if len(cmd) > 4 {
				db.Datas[i].IsTimeGiven = true
				sec := cmd[4]
				exp, err := strconv.ParseInt(sec, 0, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					e := model.Error{Status: "fail", Message: "Please Enter input in correct format"}
					json.NewEncoder(w).Encode(e)
					return
				}
				db.Datas[i].ExpireTime = time.Now().Local().Add(time.Second * time.Duration(int(exp)))
			}
			w.WriteHeader(http.StatusCreated)
			res := response{Status: "success", Data: data{Key: v.Key, Value: val}}
			json.NewEncoder(w).Encode(res)

			return
		}
	}
}

func SetData(cmd []string, w http.ResponseWriter) {
	var newData model.Data
	key := cmd[1]
	// SET key_c 4 EX 60 NX
	// 0     1   2 3  4  5

	if cmd[len(cmd)-1] == "NX" {
		isPresent, _ := isFound(key)
		if isPresent {
			w.WriteHeader(http.StatusBadRequest)
			e := model.Error{Status: "fail", Message: "key already exist"}
			json.NewEncoder(w).Encode(e)
			return
		}
	}
	if cmd[len(cmd)-1] == "XX" {
		isPresent, _ := isFound(key)
		if !isPresent {
			w.WriteHeader(http.StatusBadRequest)
			e := model.Error{Status: "fail", Message: "key not exist"}
			json.NewEncoder(w).Encode(e)
			return
		} else {
			// key already exist, update it
			update(cmd, w)
			return
		}
	}
	if len(cmd) > 3 {
		newData.IsTimeGiven = true
		sec := cmd[4]
		exp, err := strconv.ParseInt(sec, 0, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			e := model.Error{Status: "fail", Message: "Please Enter input in correct format"}
			json.NewEncoder(w).Encode(e)
			return
		}
		newData.ExpireTime = time.Now().Local().Add(time.Second * time.Duration(int(exp)))
	}
	newData.Key = cmd[1]
	newData.Value = cmd[2]

	db.Datas = append(db.Datas, newData)

	w.WriteHeader(http.StatusCreated)
	res := response{Status: "success", Data: data{Key: newData.Key, Value: newData.Value}}
	json.NewEncoder(w).Encode(res)
}

func GetData(cmd []string, w http.ResponseWriter) {
	if len(cmd) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		e := model.Error{Status: "fail", Message: "invalid command"}
		json.NewEncoder(w).Encode(e)
		return
	}

	isPresent, v := isFound(cmd[1])

	if isPresent {
		w.WriteHeader(http.StatusOK)
		res := response{Status: "success", Data: data{Key: v.Key, Value: v.Value}}
		json.NewEncoder(w).Encode(res)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		e := model.Error{Status: "fail", Message: "key not found"}
		json.NewEncoder(w).Encode(e)
		return
	}
}
