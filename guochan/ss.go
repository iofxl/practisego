package main

import (
	"log"
	"net"
	"practisego/guochan/guochan"
)

// SServer is ...
type SServer struct {
	Addr     string
	Hostname string
	cfg      Config
	logger   *log.Logger
}

// Serve is Serve
func (s *SServer) Serve(l net.Listener) {

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			s.logger.Println(err)
			return
		}
		ss := s.newSession(c)
		go ss.Serve()

	}

}

// ListenAndServe is ListenAndServe
func (s *SServer) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":443"
	}

	//	l, err := net.Listen("tcp", addr)
	g := guochan.New国产器(guochan.Method(s.cfg.Method), []byte(cfg.Secret))
	l, err := guochan.Listen(g, "tcp", addr)
	if err != nil {
		return err
	}
	s.Serve(l)
	return nil

}

// ListenAndServeSS is ...
func ListenAndServeSS(addr string, cfg Config, logger *log.Logger) error {
	s := &SServer{Addr: addr, cfg: cfg, logger: logger}
	return s.ListenAndServe()
}

// SSession is Session
type SSession struct {
	s *SServer
	net.Conn
}

func (s *SServer) newSession(conn net.Conn) *SSession {
	return &SSession{
		s:    s,
		Conn: conn,
	}
}

// Serve is ...
func (ss *SSession) Serve() {
	defer ss.Close()

	addr, err := ss.ReadAddr()
	if err != nil {
		ss.s.logger.Println(err)
		return
	}

	ss.s.logger.Printf("proxy: %s\n", addr.String())

	dstconn, err := net.Dial("tcp", addr.String())
	if err != nil {
		ss.s.logger.Println(err)
		return
	}
	err = ss.SendReply(0x00)
	if err != nil {
		ss.s.logger.Println(err)
		return
	}

	Proxy(dstconn, ss)
	return

}

// ReadByte is ReadByte
func (ss *SSession) ReadByte() (byte, error) {
	b := make([]byte, 1, 1)
	_, err := ss.Read(b)
	return b[0], err
}

// ReadAddr is used for read dst addr
func (ss *SSession) ReadAddr() (Addr, error) {

	var addr Addr

	atyp, err := ss.ReadByte()
	if err != nil {
		return addr, err
	}

	switch atyp {

	case AddrTypeIPv4:
		b := make([]byte, 6)
		n, err := ss.Read(b)

		if err != nil {
			return addr, err
		}
		addr.IP = net.IP(b[:n-2])
		addr.Port = int(b[n-2])<<8 | int(b[n-1])
	case AddrTypeIPv6:
		b := make([]byte, 18)
		n, err := ss.Read(b)
		if err != nil {
			return addr, err
		}
		addr.IP = net.IP(b[:n-2])
		addr.Port = int(b[n-2])<<8 | int(b[n-1])
	case AddrTypeFQDN:
		// read url length
		urllen, err := ss.ReadByte()

		if err != nil {
			return addr, err
		}
		b := make([]byte, urllen+2)

		n, err := ss.Read(b)

		if err != nil {
			return addr, err
		}
		addr.Name = string(b[:n-2])
		addr.Port = int(b[n-2])<<8 | int(b[n-1])
	default:
		return addr, ErrATYPNotSupported
	}
	return addr, nil
}

// SendReply is SendReply
func (ss *SSession) SendReply(reply byte) error {

	_, err := ss.Write([]byte{Version5, reply, RSV, AddrTypeIPv4, 0, 0, 0, 0, 0, 0})

	return err
}
