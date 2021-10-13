package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"encoding/json"
	"errors"
	"strings"
)

func CheckClusterNodeNull(id string) error {
	clusterNodesModel := models.ClusterNodes{}
	clusterNodeInfo := clusterNodesModel.ClusterNodeInfoById(id)
	if clusterNodeInfo.ID != id {
		return errors.New(enums.CodeMessages(enums.ClusterNodeNull))
	}

	return nil
}

func CheckClusterNodeEnableChange(id string, enable int) error {
	clusterNodesModel := models.ClusterNodes{}
	clusterNodeInfo := clusterNodesModel.ClusterNodeInfoById(id)
	if clusterNodeInfo.IsEnable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func CheckClusterNodeEnableOn(id string) error {
	clusterNodesModel := models.ClusterNodes{}
	clusterNodeInfo := clusterNodesModel.ClusterNodeInfoById(id)
	if clusterNodeInfo.IsEnable == utils.EnableOn {
		return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
	}

	return nil
}

func ClusterNodeSwitchEnable(id string, enable int) error {
	clusterNodesModel := models.ClusterNodes{}
	updateErr := clusterNodesModel.ClusterNodeSwitchEnable(id, enable)
	if updateErr != nil {
		return updateErr
	}

	// @todo 触发远程发送数据 开启/停止 网关服务，会保持与远程服务的通信

	return nil
}

type ClusterNodeListInfo struct {
	ID         string `json:"id"`
	NodeIP     string `json:"node_ip"`
	NodeStatus int    `json:"node_status"`
	IsEnable   int    `json:"is_enable"`
}

func (c *ClusterNodeListInfo) ClusterNodeListPage(param *validators.ClusterNodeList) ([]ClusterNodeListInfo, int, error) {
	clusterNodesModel := models.ClusterNodes{}
	clusterNodeListInfos, total, clusterNodeListInfosErr := clusterNodesModel.ClusterNodeListPage(param)

	clusterNodeList := make([]ClusterNodeListInfo, 0)
	if len(clusterNodeListInfos) != 0 {
		for _, clusterNodeListInfo := range clusterNodeListInfos {
			clusterNodeInfo := ClusterNodeListInfo{}
			clusterNodeInfo.ID = clusterNodeListInfo.ID
			clusterNodeInfo.NodeIP = clusterNodeListInfo.NodeIP
			clusterNodeInfo.NodeStatus = clusterNodeListInfo.NodeStatus
			clusterNodeInfo.IsEnable = clusterNodeListInfo.IsEnable

			clusterNodeList = append(clusterNodeList, clusterNodeInfo)
		}
	}

	return clusterNodeList, total, clusterNodeListInfosErr
}

func ClusterNodeDelete(id string) error {
	clusterNodesModel := models.ClusterNodes{}
	deleteErr := clusterNodesModel.ClusterNodeDelete(id)
	if deleteErr != nil {
		return deleteErr
	}

	// @todo 触发远程发送数据 开启/停止 网关服务，停止与远程服务的通信

	return nil
}

func ClusterNodeWatchAdd(watchValue string) {
	defer func() {
		if r := recover(); r != nil {
			// @todo 记录日志信息（详细报错信息）
		}
	}()

	type clusterNodeInfo struct {
		NodeIP   string `json:"node_ip"`
		IsEnable int    `json:"is_enable"`
	}

	clusterNode := clusterNodeInfo{}
	err := json.Unmarshal([]byte(watchValue), &clusterNode)
	if err != nil {
		// @todo 记录日志错误信息
	}

	ipTypeName, ipTypeNameErr := utils.DiscernIP(clusterNode.NodeIP)
	if ipTypeNameErr != nil {
		// @todo 记录日志错误信息
	}

	ipType, ipTypeErr := utils.IPNameToType(ipTypeName)
	if ipTypeErr != nil {
		// @todo 记录日志错误信息
	}

	clusterNodesModel := models.ClusterNodes{
		NodeIP:     clusterNode.NodeIP,
		IPType:     ipType,
		NodeStatus: utils.ClusterNodeStatusHealth,
		IsEnable:   clusterNode.IsEnable,
	}

	clusterNodesInfo := clusterNodesModel.ClusterNodeInfoByIp(clusterNode.NodeIP)

	if strings.TrimSpace(clusterNodesInfo.NodeIP) == strings.TrimSpace(clusterNode.NodeIP) {

		updateErr := clusterNodesModel.ClusterNodeUpdate(clusterNodesInfo.ID, &clusterNodesModel)
		if updateErr != nil {
			// @todo 记录日志错误信息
		}

	} else {
		addErr := clusterNodesModel.ClusterNodeAdd(&clusterNodesModel)
		if addErr != nil {
			// @todo 记录日志错误信息
		}
	}

	return
}
