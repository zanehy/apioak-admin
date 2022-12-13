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

func PluginTypeList(c *gin.Context) {
	pluginAllTypes := utils.PluginAllTypes()

	utils.Ok(c, pluginAllTypes)
}

func PluginUpdate(c *gin.Context) {
	pluginResId := strings.TrimSpace(c.Param("id"))

	var validatorPluginUpdate = validators.ValidatorPluginUpdate{}
	if msg, err := packages.ParseRequestParams(c, &validatorPluginUpdate); err != nil {
		utils.Error(c, msg)
		return
	}
	validators.GetPluginUpdateAttributesDefault(&validatorPluginUpdate)

	pluginModel := models.Plugins{}
	pluginInfos, pluginInfosErr := pluginModel.PluginInfosByResIds([]string{pluginResId})
	if pluginInfosErr != nil {
		utils.Error(c, pluginInfosErr.Error())
		return
	}
	if len(pluginInfos) == 0 {
		utils.Error(c, enums.CodeMessages(enums.PluginNull))
		return
	}

	updateErr := services.PluginUpdate(pluginResId, &validatorPluginUpdate)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func PluginDelete(c *gin.Context) {
	pluginResId := strings.TrimSpace(c.Param("id"))

	pluginModel := models.Plugins{}
	pluginInfos, pluginInfosErr := pluginModel.PluginInfosByResIds([]string{pluginResId})
	if pluginInfosErr != nil {
		utils.Error(c, pluginInfosErr.Error())
		return
	}

	if len(pluginInfos) == 0 {
		utils.Error(c, enums.CodeMessages(enums.PluginNull))
		return
	}

	routePluginModel := models.RoutePlugins{}
	routePluginInfos, routePluginInfosErr := routePluginModel.RoutePluginInfosByPluginResIds([]string{pluginResId})
	if routePluginInfosErr != nil {
		utils.Error(c, routePluginInfosErr.Error())
		return
	}
	if len(routePluginInfos) != 0 {
		utils.Error(c, enums.CodeMessages(enums.PluginRouteExist))
		return
	}

	deleteErr := pluginModel.PluginDelete(pluginResId)
	if deleteErr != nil {
		utils.Error(c, deleteErr.Error())
		return
	}

	utils.Ok(c)
}

func PluginList(c *gin.Context) {
	var validatorPluginList = validators.PluginList{}
	if msg, err := packages.ParseRequestParams(c, &validatorPluginList); err != nil {
		utils.Error(c, msg)
		return
	}

	structPluginInfo := services.StructPluginInfo{}
	routeList, total, err := structPluginInfo.PluginListPage(&validatorPluginList)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = validatorPluginList
	result.Page = validatorPluginList.Page
	result.PageSize = validatorPluginList.PageSize
	result.Total = total
	result.Data = routeList

	utils.Ok(c, result)
}

func PluginInfo(c *gin.Context) {
	pluginResId := strings.TrimSpace(c.Param("id"))

	checkPluginExistErr := services.CheckPluginExist(pluginResId)
	if checkPluginExistErr != nil {
		utils.Error(c, checkPluginExistErr.Error())
		return
	}

	pluginInfoService := services.PluginInfoService{}
	pluginInfo, pluginInfoErr := pluginInfoService.PluginInfoByResId(pluginResId)
	if pluginInfoErr != nil {
		utils.Error(c, pluginInfoErr.Error())
		return
	}

	utils.Ok(c, pluginInfo)
}

func PluginAddList(c *gin.Context) {
	pluginAddListItem := services.PluginAddListItem{}
	pluginAddList, err := pluginAddListItem.PluginAddList()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, pluginAddList)
}
