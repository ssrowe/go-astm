package net

type Mode int

const (
	TCP_CLIENT            Mode = 1
	TCP_SERVER            Mode = 2
	TCP_CLIENT_AND_SERVER Mode = 3
	TCP_SFTP              Mode = 4
	TCP_FTP               Mode = 5
)

type NetworkAPI interface {
	Run()  // Blocking, run thread
	Stop() // Clean exit from blocking run
}

type Handler interface {
	Incoming([]byte)
}
