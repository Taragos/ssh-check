package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/ssh"
)

type ConnResult struct {
	Server  string
	Success bool
}

func main() {

	username := flag.String("u", "", "username to use for authentication")
	password := flag.String("p", "", "password to use for authentication")
	serverList := flag.String("s", "", "path to server list to try out")
	activateDebug := flag.Bool("d", false, "activate debug log")

	flag.Parse()

	if *username == "" || *password == "" || *serverList == "" {
		log.Error("A parameter was left empty, aborting script")
		os.Exit(1)
	}

	if *activateDebug {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("Creating ssh client config")
	config := &ssh.ClientConfig{
		User: *username,
		Auth: []ssh.AuthMethod{
			ssh.Password(*password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	log.Debug("Reading server list")

	servers, err := readServerList(*serverList)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	c := make(chan ConnResult, len(servers))

	for _, server := range servers {
		wg.Add(1)

		go func(serverName string) {
			defer wg.Done()

			log.Debug("Checking connection to: ", serverName)
			success := checkConnection(serverName, config)
			result := ConnResult{
				Server:  serverName,
				Success: success,
			}
			c <- result
			log.Debug("Done checkin connection to: %v", serverName)
		}(server)
	}

	wg.Wait()
	close(c)

	serverMap := make(map[string]bool)

	for result := range c {
		serverMap[result.Server] = result.Success
	}

	for _, server := range servers {
		result := serverMap[server]
		resultString := ""

		if result {
			resultString = "Success"
		} else {
			resultString = "Failure"
		}

		fmt.Printf("%v\t\t%v\n", resultString, server)
	}
}

func checkConnection(server string, clientConfig *ssh.ClientConfig) bool {
	client, err := ssh.Dial("tcp", server, clientConfig)

	if err != nil {
		return false
	}
	client.Close()

	return true
}

func readServerList(serverPath string) ([]string, error) {
	file, err := os.Open(serverPath)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, ":") {
			line = line + ":22"
		}

		lines = append(lines, line)
	}

	return lines, scanner.Err()
}
