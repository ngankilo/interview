package model

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

type DataObj struct {
	Ping     int64  `json:"ping"`
	ServerID string `json:"server_id"`
}

type PingObj struct {
	Data DataObj `json:"data"`
	Type string  `json:"type"`
	Id   string  `json:"id"`
}

func (p *PingObj) Validate() error {
	if p.Data.Ping <= 0 {
		log.Warn().Int64("ping", p.Data.Ping).Msg("Invalid ping timestamp")
		return errors.New("invalid ping timestamp")
	}
	if p.Type == "" {
		log.Warn().Msg("Missing ping type")
		return errors.New("missing ping type")
	}
	return nil
}

func NewPingObj(serverID, id string) *PingObj {
	return &PingObj{
		Data: DataObj{
			Ping:     time.Now().Unix(),
			ServerID: serverID,
		},
		Type: "ping",
		Id:   id,
	}
}