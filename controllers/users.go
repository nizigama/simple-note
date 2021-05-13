package controllers

import "net/http"

type UsersController struct{}

func NewUsersController() UsersController {
	return UsersController{}
}

func (uc UsersController) Login(w http.ResponseWriter, req *http.Request) {

}
