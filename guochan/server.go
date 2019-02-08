package main

/*
	g := guochan.New国产器(cfg.M, []byte(cfg.Secret))

	log.Printf("Listen port: %v\n", cfg.Port)
	l, err := guochan.Listen(g, "tcp4", ":"+strconv.Itoa(cfg.Port))

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Println(err)
			break
		}

		go handleConnS(conn)

	}

*/
/*
func handleConnS(conn net.Conn) {

	b := make([]byte, 2048)

	n, err := conn.Read(b)

	if err != nil {
		log.Printf("GetAddr error", err)
		return
	}

	dst, err := 国产.ParseAddr(b[:n])

	if err != nil {
		log.Println(err)
		return
	}

	dstconn, err := net.Dial("tcp", dst.String())

	if err != nil {
		log.Println(err)
		国产.SendReply(conn, 0x03)
		return
	}

	国产.SendReply(conn, 0x00)

	var wg sync.WaitGroup
	f := func(dst, src net.Conn) {
		wg.Add(1)
		defer wg.Done()
		io.Copy(dst, src)
		dst.Close()
	}
	go f(conn, dstconn)
	go f(dstconn, conn)
	log.Printf("proxy: %s <-> %s <-> %s(%s)\n", conn.RemoteAddr(), conn.LocalAddr(), dst.String(), dstconn.RemoteAddr())
	wg.Wait()

}
*/
