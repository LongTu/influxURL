package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	influxDbTag          = "/influxdbUrl"
	portNum              = "18080"
	defaultStartTime     = "1970-01-01 00:00:00.000"
	timeLengthIndex  int = 19
	timeStampSuffix      = ".000"
)

//Note: parameter name in this struct needs to start with upper case letter
//PodID and Metric are required fields, the rest are optional
type jsonStruct struct {
	PodID     string
	TimeStart string
	TimeEnd   string
	Limit     int
	Metric    string
}

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func main() {
	http.HandleFunc(influxDbTag, influxDBHandler)
	log.Fatal(http.ListenAndServe(":"+portNum, nil))
}

func influxDBHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t jsonStruct
	err := decoder.Decode(&t)
	if err != nil {
		return
	}

	/*
	log.Println(t.PodID)
	log.Println(t.TimeStart)
	log.Println(t.TimeEnd)
	log.Println(t.Limit)
	log.Println(t.Metric)
	*/
	podID := t.PodID
	startTime := t.TimeStart
	endTime := t.TimeEnd
	limitNum := strconv.Itoa(t.Limit)
	metrics := t.Metric

	if len(metrics) <= 0 {
		//log.Println("no metrics")
		return
	}

	if len(startTime) <= 0 {
		startTime = defaultStartTime
	}

	sql := "SELECT * FROM \"" + metrics + "\" WHERE time > '" + startTime + "'"
	if len(endTime) <= 0 {
		sql += " AND time < now()"
	}
	if len(podID) > 0 {
		sql += " AND PodID='" + podID + "'"
	}

	if len(limitNum) > 0 {
		sql += " LIMIT " + limitNum
	}

	//log.Println(sql)
	res, err := readInfluxDb(sql, metrics)
	if err != nil {
		log.Println(err)
		return
	}
	//log.Println(string(res[:]))
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(res[:])

	return
}

func readInfluxDb(command string, metrics string) (res []byte, err error) {
	credential, err := getCredentials()

	if err != nil {
		return res, err
	}

	//test username and pw
	if len(credential[0]) <= 0 {
		return res, errors.New("no username")
	}

	if len(credential[1]) <= 0 {
		return res, errors.New("no password")
	}

	if len(credential[2]) <= 0 {
		return res, errors.New("no dbURL")
	}

	if len(credential[3]) <= 0 {
		return res, errors.New("no dbName")
	}
	linkURL := credential[2] + "?u=" + credential[0] + "&p=" + credential[1] + "&pretty=true"
	dbNameParam := "db=" + credential[3]
	//command = "SELECT * FROM \"cpu/usage_ns_cumulative\" WHERE time > '1970-01-01 00:00:00' AND time <now() LIMIT 1"
	SQLParam := "q=" + command
	c, err := exec.Command("curl", "-G", linkURL, "--data-urlencode", dbNameParam, "--data-urlencode", SQLParam).Output()
	if err != nil {
		return res, err
	}
	res = c
	return res, nil
}

func decypher(command string) (s string, err error) {
	keyText := "astaxie12798akljzmknm.ahkjkljl;k"

	// Create the aes encryption algorithm
	c, err := aes.NewCipher([]byte(keyText))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(keyText), err)
		os.Exit(-1)
	}

	ciphertext2, _ := hex.DecodeString(command)

	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(ciphertext2))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext2)
	//fmt.Printf("%x=>%s\n", ciphertext2, plaintextCopy)
	s = string(plaintextCopy[:])
	return s, nil
}

func getCredentials() (credential []string, err error) {
	credential = make([]string, 4)
	absPath, _ := filepath.Abs("influxdbUrl/credential.config")
	file, err := os.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	username := ""
	password := ""
	dbURL := ""
	dbName := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "u=") {
			username = line[2:]
			continue
		}
		if strings.Contains(line, "p=") {
			password = line[2:]
			continue
		}
		if strings.Contains(line, "l=") {
			dbURL = line[2:]
			continue
		}
		if strings.Contains(line, "d=") {
			dbName = line[2:]
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//test username and pw
	if len(username) <= 0 {
		return credential, errors.New("no username")
	}

	realUsername, _ := decypher(username)

	if len(password) <= 0 {
		return credential, errors.New("no password")
	}
	realPassword, _ := decypher(password)

	if len(dbURL) <= 0 {
		return credential, errors.New("no dbURL")
	}
	realURL, _ := decypher(dbURL)

	if len(dbName) <= 0 {
		return credential, errors.New("no dbName")
	}
	realDBName, _ := decypher(dbName)

	credential[0] = realUsername
	credential[1] = realPassword
	credential[2] = realURL
	credential[3] = realDBName
	return credential, nil
}
