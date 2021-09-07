package admin

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
	"strings"
)

func RouteAdd(c *gin.Context) {

	var validatorRouteAddUpdate = validators.ValidatorRouteAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &validatorRouteAddUpdate); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetRouteAttributesDefault(&validatorRouteAddUpdate)

	if validatorRouteAddUpdate.RoutePath == utils.DefaultRoutePath {
		utils.Error(c, enums.CodeMessages(enums.RouteDefaultPathNoPermission))
		return
	}

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(validatorRouteAddUpdate.ServiceID)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	err := services.CheckExistServiceRoutePath(validatorRouteAddUpdate.ServiceID, validatorRouteAddUpdate.RoutePath, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	createErr := services.RouteCreate(&validatorRouteAddUpdate)
	if createErr != nil {
		utils.Error(c, createErr.Error())
		return
	}

	utils.Ok(c)
}

func RouteList(c *gin.Context) {
	serviceId := strings.TrimSpace(c.Param("service_id"))

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	var validatorRouteList = validators.ValidatorRouteList{}
	if msg, err := packages.ParseRequestParams(c, &validatorRouteList); err != nil {
		utils.Error(c, msg)
		return
	}

	structRouteList := services.StructRouteList{}
	routeList, total, err := structRouteList.RouteListPage(serviceId, &validatorRouteList)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = validatorRouteList
	result.Page = validatorRouteList.Page
	result.PageSize = validatorRouteList.PageSize
	result.Total = total
	result.Data = routeList

	utils.Ok(c, result)
}
