package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Secret  []byte
	Salt    []byte
	Info    []byte
	M       Method
	Address string
	Server  string
}

func saveConfig(c *Config, filename string) error {

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	defer f.Close()

	if err != nil {
		return err
	}

	enc := gob.NewEncoder(f)
	err = enc.Encode(c)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(filename)

	s := base64.StdEncoding.EncodeToString(b)

	fmt.Println(s)

	return nil

}

func loadConfig(filename string) (*Config, error) {

	cfg := new(Config)

	f, err := os.Open(filename)
	if err != nil {
		return cfg, err
	}
	dec := gob.NewDecoder(f)

	err = dec.Decode(cfg)

	return cfg, err
}

func decodeConfig(s string) (*Config, error) {

	cfg := new(Config)

	b, err := base64.StdEncoding.DecodeString(s)

	if err != nil {
		return cfg, err
	}

	dec := gob.NewDecoder(bytes.NewBuffer(b))

	err = dec.Decode(cfg)

	return cfg, err

}
