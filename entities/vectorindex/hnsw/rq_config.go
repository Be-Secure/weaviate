package hnsw

import (
	"github.com/weaviate/weaviate/entities/vectorindex/common"
)

const (
	DefaultRQEnabled = false
	DefaultDataBits  = 1
	DefaultQueryBits = 1
)

type RQConfig struct {
	Enabled   bool  `json:"enabled"`
	DataBits  int16 `json:"dataBits"`
	QueryBits int16 `json:"queryBits"`
}

func parseRQMap(in map[string]interface{}, rq *RQConfig) error {
	rqConfigValue, ok := in["rq"]
	if !ok {
		return nil
	}

	rqConfigMap, ok := rqConfigValue.(map[string]interface{})
	if !ok {
		return nil
	}

	if err := common.OptionalBoolFromMap(rqConfigMap, "enabled", func(v bool) {
		rq.Enabled = v
	}); err != nil {
		return err
	}

	if err := common.OptionalIntFromMap(rqConfigMap, "dataBits", func(v int) {
		rq.DataBits = int16(v)
	}); err != nil {
		return err
	}

	if err := common.OptionalIntFromMap(rqConfigMap, "queryBits", func(v int) {
		rq.QueryBits = int16(v)
	}); err != nil {
		return err
	}

	return nil
}
