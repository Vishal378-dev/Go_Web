package Handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	utils "github.com/vishal/reservation_system/Handlers/Utils"
)

func WrongPathTemplate(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("./static/*")
	fmt.Println(os.Getwd())
	if err != nil {
		fmt.Println(err)
		utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("EndPoint 1 Not Implemented"), http.StatusNotImplemented))
		return
	}
	err = tpl.ExecuteTemplate(w, "RouteNotFound.html", nil)
	if err != nil {
		fmt.Println(err)
		utils.ResponseWriter(w, http.StatusNotImplemented, utils.CommonError(fmt.Errorf("EndPoint 2 Not Implemented"), http.StatusNotImplemented))
		return
	}
	return
}
