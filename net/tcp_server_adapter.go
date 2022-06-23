package net

type SENDPROXY int

const (
	NO_PROXY_PROTOCOL SENDPROXY = 0
	SENDPROXY_V1      SENDPROXY = 1
	SENDPROXY_V2      SENDPROXY = 2
)

type ProtocolFSMState struct {
}

type TransmissionProtocol int

const (
	PlainDataProtocol TransmissionProtocol = 0
)

type TCPServer struct {
	port          int
	proxy         SENDPROXY
	maxConnection int
	timing        ServerTimings
	protocol      TransmissionProtocol
}

type ServerTimings struct {
	connectionTimeout_ms int
	sessionTimeout_ms    int
}

func CreateTCPServer(port int, proxy SENDPROXY,
	maxConnection int,
	protocol TransmissionProtocol, timing ServerTimings) NetworkAPI {
	return &TCPServer{
		port:          port,
		proxy:         proxy,
		maxConnection: maxConnection,
		timing:        timing,
		protocol:      protocol,
	}
}

func (server *TCPServer) Run() {

}

func (server *TCPServer) Stop() {

}
