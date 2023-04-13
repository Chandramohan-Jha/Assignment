package routers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/chandramohan/in_memory/controllers"
	"github.com/chandramohan/in_memory/model"
)

type incomingData struct {
	Command string `json:"command"`
}

func Route(w http.ResponseWriter, r *http.Request) {
	var data incomingData
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &data)
	cmd := strings.Split(data.Command, " ")
	switch cmd[0] {
	case "SET":
		controllers.SetData(cmd, w)
		// fmt.Println("SET Command")
	case "GET":
		controllers.GetData(cmd, w)
		// fmt.Println("GET Command")
	case "QPUSH":
		controllers.QPush(cmd, w)
		// fmt.Println("QPUSH Command")
	case "QPOP":
		controllers.QPop(cmd, w)
		// fmt.Println("QPOP Command")
	case "BQPOP":
		controllers.BQPop(cmd, w)
		// fmt.Println("BQPOP Command")
	default:
		w.WriteHeader(http.StatusBadRequest)
		e := model.Error{Status: "fail", Message: "invalid command"}
		json.NewEncoder(w).Encode(e)
		return
	}

}
