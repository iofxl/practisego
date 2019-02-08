package main

/*
func handleConn(conn net.Conn) {

	err := 国产.Negotiate(conn)

	if err != nil {
		log.Println("Negotiation error:", err)
	}

	err = 国产.HandleRequest(conn)
	if err != nil {
		log.Println("HandleRequest error:", err)
		return
	}

	b := make([]byte, 2048)

	// read addr here
	n, err := conn.Read(b)

	if err != nil {
		log.Printf("GetAddr error:", err)
		return
	}

	dst, err := 国产.ParseAddr(b[:n])

	if err != nil {
		log.Println("ParseAddr error:", err)
		return
	}

	// dial server here
	g := 国产.New国产器(cfg.M, []byte(cfg.Secret))
	srvconn, err := 国产.Dial(g, "tcp", cfg.Server)

	if err != nil {

		国产.SendReply(conn, 0x01)
		log.Println("DialServer error:", err)
		return
	}

	if _, err := srvconn.Write(b[:n]); err == nil {
		log.Printf("Connect: %s\n", dst.String())
	} else {
		log.Printf("Write dst %s error: %v\n", dst.String(), err)
	}

	reply, _ := 国产.ReadReply(srvconn)
	国产.SendReply(conn, reply)
	guochan.Proxy(conn, srvconn)
}
*/
