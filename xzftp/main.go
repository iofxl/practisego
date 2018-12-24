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
	"regexp"
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
	Regexp string
	Enable bool
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

		if w.Enable {
			wg.Add(1)
			log.Printf("Begin handleWork %s\n", w.Des)
			go handleWork(w, wg)
		}
	}

	wg.Wait()

}

func handleWork(w work, wg sync.WaitGroup) {
	defer wg.Done()

	switch w.Do {
	case "upload":
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

	log.Printf("%s initFTP %s:%v Successful", w.Des, w.Host, w.Port)

	for {

		entries, ftperr, err := w.getEntrynames(c)

		for ftperr != nil || err != nil {
			if ftperr != nil {
				log.Println(w.Des, ftperr)
				w.download()
			}

			if err != nil {
				log.Println(w.Des, err)
				time.Sleep(10 * time.Second)
				entries, ftperr, err = w.getEntrynames(c)

			}
		}

		for _, e := range entries {
			ftperr, err := w.downloadEntry(c, e)
			if ftperr != nil {
				log.Println(w.Des, ftperr)
				w.download()
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

		if e.Type == ftp.EntryTypeFile && e.Size != uint64(0) && w.validateFile(e.Name) {
			entrynames = append(entrynames, e.Name)
		}
	}

	if len(entrynames) == 0 {
		return nil, nil, errors.New("No availble entries left.")
	}

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

	log.Printf("%s initFTP %s:%v Successful", w.Des, w.Host, w.Port)

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
			if ftperr != nil {
				log.Println(w.Des, ftperr)
				w.upload()
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

		if !fi.IsDir() && fi.Size() != int64(0) && w.validateFile(fi.Name()) {
			files = append(files, fi.Name())
		}
	}

	if len(files) == 0 {
		return nil, errors.New("No availble files left.")
	}

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

	if w.Remote != "/" {

		err = c.ChangeDir(w.Remote)
		if err != nil {
			return c, err
		}
	}

	return c, err
}

func (w work) validateFile(filename string) bool {
	re := regexp.MustCompile(w.Regexp)
	return re.MatchString(filename) && !strings.HasPrefix(filename, w.Prefix) && !strings.HasSuffix(filename, w.Suffix)
}

func LoadConfig(cfgFile string) (works, error) {

	var ws works

	cfg, err := ioutil.ReadFile(cfgFile)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(cfg, &ws)

	for _, w := range ws {

		if w.Enable {

			w.Do = strings.TrimSpace(strings.ToLower(w.Do))

			fi, err := os.Stat(w.Local)

			if err != nil {
				log.Printf(" %s Stat %s ", w.Des, w.Local)
				log.Fatal(err)
			}

			if !filepath.IsAbs(w.Local) {
				log.Fatal(" %s %s need  an Abs path", w.Des, w.Local)
			}

			if !fi.IsDir() {
				log.Fatal(" %s %s is not a dir", w.Des, w.Local)
			}

			if w.Do == "upload" {

				fi, err := os.Stat(w.Backup)

				if err != nil {
					log.Fatal(" %s Stat %s ", w.Des, w.Local)
					log.Fatal(err)
				}

				if !fi.IsDir() {
					log.Printf(" %s  %s is not a dir", w.Des, w.Backup)
					log.Fatal(err)
				}
			}
		}
	}

	return ws, err

}
