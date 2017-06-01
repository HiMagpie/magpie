package connector

import (
	"magpie/internal/com/utils"
	"magpie/internal/com/errs"
	"magpie/internal/com/logger"
)

/**
 * cid和sessionId映射关系
 */
var relCidSessionMaps = map[string]map[string]uint64{}
var relSessionIdCidMaps = map[uint64]map[uint64]string{}

const (
	PREFIX_LEN = 2
	INDEX_INTERVAL = 10000
)

/**
 * 通过cid获取session
 */
func GetSessionByCid(cid string) (*Session, error) {
	prefix := utils.Substr(cid, 0, PREFIX_LEN)
	if _, ok := relCidSessionMaps[prefix]; !ok {
		return nil, errs.ERR_OPERATION
	}

	if _, ok := relCidSessionMaps[prefix][cid]; !ok {
		return nil, errs.ERR_OPERATION
	}

	sessionId := relCidSessionMaps[prefix][cid]
	return GetServer().GetSession(sessionId), nil
}

func GetCidBySessionId(sessionId uint64) (string, error) {
	index := sessionId / INDEX_INTERVAL
	if _, ok := relSessionIdCidMaps[index]; !ok {
		return "", errs.ERR_OPERATION
	}

	if _, ok := relSessionIdCidMaps[index][sessionId]; !ok {
		return "", errs.ERR_OPERATION
	}

	return relSessionIdCidMaps[index][sessionId], nil
}

/**
 * 添加一条cid和sessionId的映射关系
 */
func AddCidSessionIdRel(cid string, sessionId uint64) {
	logger.Debug("add cid", logger.Format(
		"cid", cid,
		"session_id ", sessionId,
	))
	prefix := utils.Substr(cid, 0, PREFIX_LEN)
	if _, ok := relCidSessionMaps[prefix]; !ok {
		relCidSessionMaps[prefix] = make(map[string]uint64)
	}
	relCidSessionMaps[prefix][cid] = sessionId

	index := sessionId / INDEX_INTERVAL
	if _, ok := relSessionIdCidMaps[index]; !ok {
		relSessionIdCidMaps[index] = make(map[uint64]string)
	}
	relSessionIdCidMaps[index][sessionId] = cid
}

/**
 * 移除Cid和SessionId的映射关系
 */
func RemoveCidRel(cid string) {
	prefix := utils.Substr(cid, 0, 2)
	if _, ok := relCidSessionMaps[prefix]; !ok {
		return
	}

	delete(relCidSessionMaps[prefix], cid)
}

/**
 * 获取所有cid
 */
func GetAllCids() []string {
	cids := make([]string, 0)
	for _, smap := range relCidSessionMaps {
		for k, _ := range smap {
			cids = append(cids, k)
		}
	}

	return cids;
}
