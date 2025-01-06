package goflag

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type ControllerFlag interface {
	Controller
	Info() *ControllerInformation
}

type Controller interface {
	http.Handler
	chi.Routes
}

type ControllerInformation struct {
	Information
	Name string
}

func (info *ControllerInformation) WithName(name string) *ControllerInformation {
	info.Name = name
	return info
}

type flaggedController struct {
	Controller
	flag[ControllerInformation]
}

func FlagController(controller Controller, info ...ControllerInformation) ControllerFlag {
	var _info ControllerInformation
	if len(info) > 0 {
		_info = info[0]
	}

	return &flaggedController{
		Controller: controller,
		flag:       flag[ControllerInformation]{info: _info},
	}
}
