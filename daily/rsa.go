package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

var createCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create")

		pri, pub, err := genKey()

		if err != nil {
			log.Fatal(err)
		}

		err = saveGobKey(args[0]+".gob", pri)
		if err != nil {
			log.Fatal(err)
		}

		err = saveGobKey(args[0]+".pub.gob", pub)

		if err != nil {
			log.Fatal(err)
		}

		err = savePEMKey(args[0]+".pem", pri)

		if err != nil {
			log.Fatal(err)
		}

	},
}

var loadCmd = &cobra.Command{
	Use: "load",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("load")

		var priv rsa.PrivateKey
		err := loadGobKey(args[0], &priv)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Private key Primes:", priv.Primes[0], priv.Primes[1])
		fmt.Println("Private key exponet:", priv.D.String())

		pub := priv.PublicKey

		fmt.Println("Public key Modulus:", pub.N.String())
		fmt.Println("Public key exponet:", pub.E)

	},
}

func init() {
	rootCmd.AddCommand(createCmd, loadCmd)
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

func genKey() (*rsa.PrivateKey, rsa.PublicKey, error) {

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return privKey, rsa.PublicKey{}, err
	}

	fmt.Println("Private key Primes:", privKey.Primes[0], privKey.Primes[1])
	fmt.Println("Private key exponet:", privKey.D.String())

	pub := privKey.PublicKey

	fmt.Println("Public key Modulus:", pub.N.String())
	fmt.Println("Public key exponet:", pub.E)

	return privKey, pub, nil

}

func saveGobKey(filename string, key interface{}) error {

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	enc := gob.NewEncoder(f)

	return enc.Encode(key)

}

func savePEMKey(filename string, key *rsa.PrivateKey) error {

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	return pem.Encode(f, block)

}

func loadGobKey(filename string, key *rsa.PrivateKey) error {

	f, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer f.Close()

	dec := gob.NewDecoder(f)

	return dec.Decode(key)

}
