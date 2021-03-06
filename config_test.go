package main

import (
	"encoding/hex"
	"fmt"
	"github.com/classzz/classzz/chaincfg"
	"github.com/classzz/classzz/czzec"
	"github.com/classzz/czzutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

var (
	rpcuserRegexp = regexp.MustCompile("(?m)^rpcuser=.+$")
	rpcpassRegexp = regexp.MustCompile("(?m)^rpcpass=.+$")
)

func TestExcessiveBlockSizeUserAgentComment(t *testing.T) {
	// Wipe test args.
	os.Args = []string{"classzz"}

	cfg, _, err := loadConfig()
	if err != nil {
		t.Fatal("Failed to load configuration")
	}

	if len(cfg.UserAgentComments) != 1 {
		t.Fatal("Expected EB UserAgentComment")
	}

	uac := cfg.UserAgentComments[0]
	uacExpected := "EB32.0"
	if uac != uacExpected {
		t.Fatalf("Expected UserAgentComments to contain %s but got %s", uacExpected, uac)
	}

	// Custom excessive block size.
	os.Args = []string{"classzz", "--excessiveblocksize=64000000"}

	cfg, _, err = loadConfig()
	if err != nil {
		t.Fatal("Failed to load configuration")
	}

	if len(cfg.UserAgentComments) != 1 {
		t.Fatal("Expected EB UserAgentComment")
	}

	uac = cfg.UserAgentComments[0]
	uacExpected = "EB64.0"
	if uac != uacExpected {
		t.Fatalf("Expected UserAgentComments to contain %s but got %s", uacExpected, uac)
	}
}

func TestCreateDefaultConfigFile(t *testing.T) {
	// Setup a temporary directory
	tmpDir, err := ioutil.TempDir("", "classzz")
	if err != nil {
		t.Fatalf("Failed creating a temporary directory: %v", err)
	}
	testpath := filepath.Join(tmpDir, "test.conf")

	// Clean-up
	defer func() {
		os.Remove(testpath)
		os.Remove(tmpDir)
	}()

	err = createDefaultConfigFile(testpath)

	if err != nil {
		t.Fatalf("Failed to create a default config file: %v", err)
	}

	content, err := ioutil.ReadFile(testpath)
	if err != nil {
		t.Fatalf("Failed to read generated default config file: %v", err)
	}

	if !rpcuserRegexp.Match(content) {
		t.Error("Could not find rpcuser in generated default config file.")
	}

	if !rpcpassRegexp.Match(content) {
		t.Error("Could not find rpcpass in generated default config file.")
	}
}

func TestGenesisAdderss(t *testing.T) {

	key, _ := czzec.NewPrivateKey(czzec.S256())
	wif, _ := czzutil.NewWIF(key, &chaincfg.MainNetParams, true)

	fmt.Println("wif:", wif.String())
	fmt.Println("priv:", hex.EncodeToString(key.Serialize()))
	pk := (*czzec.PublicKey)(&key.PublicKey).SerializeCompressed()
	fmt.Println("pub:", hex.EncodeToString(pk))
	address, err1 := czzutil.NewAddressPubKeyHash(czzutil.Hash160(pk), &chaincfg.MainNetParams)

	if err1 != nil {
		t.Errorf("failed to make address for: %v", err1)
	}
	fmt.Println(address.String())
}

func TestConvertAddr(t *testing.T) {

	keyBy, _ := hex.DecodeString("59c971df5e89a7f5816eff47575378498804224ca1e468a4d8d2afea12d9dd03")
	key, _ := czzec.PrivKeyFromBytes(czzec.S256(), keyBy)
	wif, _ := czzutil.NewWIF(key, &chaincfg.MainNetParams, true)

	fmt.Println("wif:", wif.String())
	fmt.Println("priv:", hex.EncodeToString(key.Serialize()))
	pk := (*czzec.PublicKey)(&key.PublicKey).SerializeCompressed()
	fmt.Println("pub:", hex.EncodeToString(pk))
	address, err1 := czzutil.NewAddressPubKeyHash(czzutil.Hash160(pk), &chaincfg.MainNetParams)

	if err1 != nil {
		t.Errorf("failed to make address for: %v", err1)
	}
	fmt.Println(address.String())

}

func TestGenesisRegTestAdderss(t *testing.T) {

	key, _ := czzec.NewPrivateKey(czzec.S256())
	wif, _ := czzutil.NewWIF(key, &chaincfg.RegressionNetParams, true)

	fmt.Println("wif:", wif.String())
	fmt.Println("priv:", hex.EncodeToString(key.Serialize()))
	pk := (*czzec.PublicKey)(&key.PublicKey).SerializeCompressed()
	fmt.Println("pub:", hex.EncodeToString(pk))
	address, err1 := czzutil.NewAddressPubKeyHash(czzutil.Hash160(pk), &chaincfg.RegressionNetParams)

	if err1 != nil {
		t.Errorf("failed to make address for: %v", err1)
	}
	fmt.Println(address.String())
}

func TestGenesisSimNetAdderss(t *testing.T) {

	key, _ := czzec.NewPrivateKey(czzec.S256())
	wif, _ := czzutil.NewWIF(key, &chaincfg.SimNetParams, true)

	fmt.Println("wif:", wif.String())
	fmt.Println("priv:", hex.EncodeToString(key.Serialize()))
	pk := (*czzec.PublicKey)(&key.PublicKey).SerializeCompressed()
	fmt.Println("pub:", hex.EncodeToString(pk))
	address, err1 := czzutil.NewAddressPubKeyHash(czzutil.Hash160(pk), &chaincfg.SimNetParams)

	if err1 != nil {
		t.Errorf("failed to make address for: %v", err1)
	}
	fmt.Println(address.String())
}
