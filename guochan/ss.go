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
}

// Serve is Serve
func (s *SServer) Serve(l net.Listener) error {

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		ss := s.newSession(c)

		go func() {
			err := ss.Serve()
			if err != nil {
				log.Println(err)
			}
		}()

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
	return s.Serve(l)

}

// ListenAndServeSS is ...
func ListenAndServeSS(addr string, cfg Config) error {
	s := &SServer{Addr: addr, cfg: cfg}
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
func (ss *SSession) Serve() error {
	defer ss.Close()
	// do some init work
	//
	addr, err := ss.ReadAddr()
	if err != nil {
		return err
	}

	log.Printf("proxy: %s\n", addr.String())

	dstconn, err := net.Dial("tcp", addr.String())
	if err != nil {
		return err
	}
	err = ss.SendReply(0x00)
	if err != nil {
		return err
	}

	Proxy(dstconn, ss)
	return nil

}

// ReadByte is ReadByte
func (ss *SSession) ReadByte() (byte, error) {
	b := make([]byte, 1, 1)
	_, err := ss.Read(b)
	return b[0], err
}

// Negotiate is Negotiate
func (ss *SSession) Negotiate() error {

	// read version
	ver, err := ss.ReadByte()
	if err != nil {
		return err
	}

	if ver != Version5 {
		return ErrVersion
	}

	// read nmethods
	nmethods, err := ss.ReadByte()

	if err != nil {
		return err
	}

	methods := make([]byte, int(nmethods))

	_, err = ss.Read(methods)

	if err != nil {
		return err
	}

	_, err = ss.Write([]byte{Version5, AuthMethodNotRequired})

	return err

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
