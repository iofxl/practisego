package main

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

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
	Prefix string
	Suffix string
}

type works []work

func main() {

	var cfgFile string

	var ws works

	flag.StringVar(&cfgFile, "c", "Xzftp.yaml", "config file")
	flag.Parse()

	ws, err := LoadConfig(cfgFile)

	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, w := range ws {
		wg.Add(1)
		log.Printf("Begin handleWork %s\n", w.Des)
		go handleWork(w, wg)
	}

	wg.Wait()

}

func handleWork(w work, wg sync.WaitGroup) {
	defer wg.Done()

	fi, err := os.Stat(w.Local)

	if !filepath.IsAbs(w.Local) || err != nil || !fi.IsDir() {
		log.Printf(" %s Local dir error or %s is not a dir or is not an Abs path", w.Des, w.Local)
		log.Fatal(err)
	}

	do := strings.ToLower(w.Do)

	switch do {
	case "upload":
		fi, err := os.Stat(w.Backup)

		if err != nil || !fi.IsDir() {
			log.Printf(" %s Backup dir error or %s is not a dir", w.Des, w.Backup)
			log.Fatal(err)
		}

		w.upload()
	case "download":
		w.download()
	}
}

func (w work) download() {

	c, ftperr := w.initFTP()
	for ftperr != nil {
		log.Println(w.Des, ftperr)
		time.Sleep(1 * time.Minute)
		c, ftperr = w.initFTP()
	}

	for {

		entries, ftperr, err := w.getEntrynames(c)

		for ftperr != nil || err != nil {
			if ftperr != nil {
				log.Println(w.Des, ftperr)
				time.Sleep(1 * time.Minute)
				c, ftperr = w.initFTP()
			}
			if err != nil {
				log.Println(w.Des, err)
				time.Sleep(10 * time.Second)
				entries, ftperr, err = w.getEntrynames(c)

			}

		}

		for _, e := range entries {
			ftperr, err := w.downloadEntry(c, e)
			for ftperr != nil {
				log.Println(w.Des, ftperr)
				time.Sleep(1 * time.Minute)
				c, ftperr = w.initFTP()
			}
			if err != nil {
				log.Println(w.Des, err)
				time.Sleep(5 * time.Second)
				continue
			}
			log.Printf("%s DownloadEntry %s Successful", w.Des, e)

		}
	}

}

func (w work) getEntrynames(c *ftp.ServerConn) ([]string, error, error) {
	// Returns information of a file or directory if specified, else information of the current working directory is returned.
	entries, ftperr := c.List("")
	if ftperr != nil {
		return nil, ftperr, nil
	}

	if len(entries) == 0 {
		return nil, nil, errors.New("No entries left.")
	}

	var entrynames []string
	for _, e := range entries {

		if !(e.Type == ftp.EntryTypeFolder) && !strings.HasPrefix(e.Name, w.Prefix) && !strings.HasSuffix(e.Name, w.Suffix) {
			entrynames = append(entrynames, e.Name)
		}
	}

	if len(entrynames) == 0 {
		return nil, nil, errors.New("Only tmp entry or dir left.")
	}

	log.Printf("%s Find entries %v", w.Des, entrynames)
	return entrynames, nil, nil

}

func (w work) downloadEntry(c *ftp.ServerConn, entryname string) (error, error) {
	res, ftperr := c.Retr(entryname)
	if ftperr != nil {
		return ftperr, nil
	}

	tmpfname := w.Prefix + entryname + w.Suffix

	f, err := os.Create(filepath.Join(w.Local, tmpfname))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, res)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}
	err = os.Rename(filepath.Join(w.Local, tmpfname), filepath.Join(w.Local, entryname))
	if err != nil {
		return nil, err
	}

	ftperr = res.Close()
	if ftperr != nil {
		return ftperr, nil
	}
	ftperr = c.Delete(entryname)
	return ftperr, nil
}

func (w work) upload() {

	c, ftperr := w.initFTP()
	for ftperr != nil {
		log.Println(w.Des, ftperr)
		time.Sleep(1 * time.Minute)
		c, ftperr = w.initFTP()
	}

	for {

		files, err := w.getFilenames()

		for err != nil {
			log.Println(w.Des, err)
			time.Sleep(10 * time.Second)
			files, err = w.getFilenames()
			c.NoOp()
		}

		for _, file := range files {
			ftperr, err := w.uploadFile(c, filepath.Join(w.Local, file))
			for ftperr != nil {
				log.Println(w.Des, ftperr)
				time.Sleep(1 * time.Minute)
				c, ftperr = w.initFTP()
			}
			if err != nil {
				log.Println(w.Des, err)
				time.Sleep(5 * time.Second)
				continue
			}
			log.Printf("%s uploadFile %s Successful", w.Des, file)
		}
	}

}

func (w work) uploadFile(c *ftp.ServerConn, file string) (error, error) {

	basename := filepath.Base(file)

	rfname := w.Prefix + basename + w.Suffix

	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	ftperr := c.Stor(rfname, f)

	if ftperr != nil {
		return ftperr, nil
	}

	ftperr = c.Rename(rfname, basename)

	if ftperr != nil {
		return ftperr, nil
	}
	err = f.Close()

	if err != nil {
		return nil, err
	}

	err = os.Rename(file, filepath.Join(w.Backup, filepath.Base(file)))

	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	return nil, nil
}

func (w work) getFilenames() ([]string, error) {

	fis, err := ioutil.ReadDir(w.Local)

	if err != nil {
		return nil, err
	}

	if len(fis) == 0 {
		return nil, errors.New("No files left.")
	}

	var files []string

	for _, fi := range fis {

		if !fi.IsDir() && !strings.HasPrefix(fi.Name(), w.Prefix) && !strings.HasSuffix(fi.Name(), w.Suffix) {
			files = append(files, fi.Name())
		}
	}

	if len(files) == 0 {
		return nil, errors.New("Only tmp file or dir left.")
	}

	log.Printf("%s Find files %v", w.Des, files)
	return files, nil
}

func (w work) initFTP() (*ftp.ServerConn, error) {

	addr := net.JoinHostPort(w.Host, strconv.Itoa(w.Port))

	c, err := ftp.DialTimeout(addr, 5*time.Second)

	if err != nil {
		return c, err
	}

	err = c.Login(w.User, w.Pass)

	if err != nil {
		return c, err
	}

	log.Printf("%s Login FTP %s:%v Successful", w.Des, w.Host, w.Port)

	if w.Remote != "/" {

		err = c.ChangeDir(w.Remote)
		if err != nil {
			return c, err
		}
		log.Printf("%s FTP ChangeDir to %s Successful", w.Des, w.Remote)
	}

	return c, err
}

func LoadConfig(cfgFile string) (works, error) {

	var w works

	cfg, err := ioutil.ReadFile(cfgFile)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(cfg, &w)

	return w, err

}
