package handler

import "net/http"

func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) {

	//get login + pass from req

	//validate for some points

	//try to userService.Login() - should return user_id

	//if success - generate tokens

	//set tokens in the redis map

	//set tokens into cookies
}
