package server

type (
	ServerConfigs struct {
		Addr       string
		LeaderAddr string
		IsLeader   bool
	}

	Server struct {
		ServerConfigs
	}
)
