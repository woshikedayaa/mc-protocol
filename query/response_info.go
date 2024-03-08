package query

import (
	"encoding/json"
)

func (r *BasicResponse) JSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"sessionID": r.SessionID(),
		"queryType": r.typ,
		"gameType":  r.gameType,
		"motd":      r._MOTD,
		"map":       r._map,
		"numPlayer": r.curPlayers,
		"maxPlayer": r.maxPlayer,
		"port":      r.port,
		"hostip":    r.ip,
	})
}

func (r *BasicResponse) GameType() string {
	return r.gameType
}

func (r *BasicResponse) MOTD() string {
	return r._MOTD
}

func (r *BasicResponse) Map() string {
	return r._map
}

func (r *BasicResponse) NumPlayer() int {
	return r.curPlayers
}

func (r *BasicResponse) MaxPlayer() int {
	return r.maxPlayer
}

func (r *BasicResponse) Port() int {
	return r.port
}

func (r *BasicResponse) HostIP() string {
	return r.ip
}

func (r *FullResponse) JSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"sessionID": r.SessionID(),
		"queryType": r.typ,
		"gameType":  r.gameType,
		"motd":      r._MOTD,
		"map":       r._map,
		"numPlayer": r.curPlayers,
		"maxPlayer": r.maxPlayer,
		"port":      r.port,
		"hostip":    r.ip,
		"player":    r.player,
		"plugin":    r.plugins,
		"version":   r.version,
		"gameId":    r.gameID,
	})
}

func (r *FullResponse) GameType() string {
	return r.gameType
}

func (r *FullResponse) MOTD() string {
	return r._MOTD
}

func (r *FullResponse) Map() string {
	return r._map
}

func (r *FullResponse) NumPlayer() int {
	return r.curPlayers
}

func (r *FullResponse) MaxPlayer() int {
	return r.maxPlayer
}

func (r *FullResponse) Port() int {
	return r.port
}

func (r *FullResponse) HostIP() string {
	return r.ip
}

func (r *FullResponse) Players() []string {
	return r.player
}

func (r *FullResponse) GameID() string {
	return r.gameID
}

func (r *FullResponse) Version() string {
	return r.version
}

func (r *FullResponse) Plugins() string {
	// todo parse
	return r.plugins
}

func (r *HandleShakeResponse) JSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"sessionID": r.SessionID(),
		"queryType": r.typ,
		"token":     r.token,
	})
}

func (r *HandleShakeResponse) Token() int32 {
	return r.token
}
