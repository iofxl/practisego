package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/jlaffaye/ftp"
	yaml "gopkg.in/yaml.v2"
)

type work struct {
	Des    string
	Do     string
	Host   string
	Port   int
	User   string
	Pass   string
	Local  string
	Remote string
	Backup string
}

var works []work
var cfgFile string

func main() {
	flag.StringVar(&cfgFile, "c", "Xzftp.yaml", "config file")
	flag.Parse()

	f, err := os.Open(cfgFile)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	dec := yaml.NewDecoder(f)

	err = dec.Decode(&works)

	if err != nil {
		log.Fatal(err)
	}

	for k, v := range works {
		fmt.Println(k, v)
	}
	var wg sync.WaitGroup

	for _, work := range works {
		wg.Add(1)
		go handlework(work, wg)

	}

	wg.Wait()
	fmt.Println("done")

}

func handlework(w work, wg sync.WaitGroup) {
	defer wg.Done()

	addr := net.JoinHostPort(w.Host, strconv.Itoa(w.Port))
	c, err := ftp.Dial(addr)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Login(w.User, w.Pass)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Login Successfull!")
	entries, err := c.NameList(w.Remote)

	if err != nil {
		log.Println(err)
		return
	}

	for k, v := range entries {
		fmt.Println(k, v)

	}

}
