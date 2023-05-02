package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"io/ioutil"
	"log"
	"strconv"
)

var targetfile = flag.String("targetfile","./targets.json",
	"JSON file specifying targets")
var bannerfile = flag.String("bannerfile","./banner.txt",
	"Text file containing banner")

type MenuItem struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int `json:"port"`
  Handler string `json:"handler"`
}

func ReadJSONMenu(jfile string) []MenuItem {
	jsonFile, err := os.Open(jfile)
	if err !=nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteVals, err := ioutil.ReadAll(jsonFile)
	if err !=nil {
		log.Fatal(err)
	}
	var result []MenuItem
	err = json.Unmarshal([]byte(byteVals), &result)
	if (err != nil) {
		log.Fatal(err)
	}
	return result
}

func ReadBanner(bannerFile string) string {
  byteVals, err := ioutil.ReadFile(bannerFile)
  if err !=nil {
    log.Fatal(err)
  }
  return string(byteVals)
}

func SetEnvs() {
  defaults := map[string]string {
    "HOME": "/nonexistent",
    "SHELL": "/usr/sbin/nologin",
    "TERM": "ansi",
  }
  for k, v := range defaults {
    if os.Getenv(k) == "" {
      os.Setenv(k,v)
    }
  }
}

func GetMenuReply(bannerfile *string, menu []MenuItem) []string {
  SetEnvs()
  banner := ReadBanner(*bannerfile)
	DisplayMenu(banner,menu)
	l := len(menu)
	scanner := bufio.NewScanner(os.Stdin)
	choice := 0
	for (choice < 1 || choice > l) {
		fmt.Printf("\n> ");
		var cs string
		var err error
		b := scanner.Scan()
		if b {
			cs = scanner.Text()
		} else {
			err = scanner.Err()
			if err != nil {
				log.Fatal(err)
			} else {
				os.Exit(0)
			}
		}
		choice, err = strconv.Atoi(cs)
		if err != nil || (choice < 1 || choice > l) {
			fmt.Printf("Error: enter number 1-%d.\n", l)
		}
	}
	cmdargs := ParseChoice(choice, menu)
	return cmdargs
}

func DisplayMenu(banner string, menu []MenuItem) {
	fmt.Printf(banner)
	fmt.Printf("\n")
	for idx, m := range menu {
		fmt.Printf("%d: %s\n", (1 + idx), m.Name)
	}
	fmt.Printf("\n")
}

func ParseChoice(choice int, menu []MenuItem) []string {
	idx := choice - 1
	item := menu[idx]
	cmd := item.Handler
  if cmd == "" {
    cmd = "socat"
  }
	host := item.Host
	port := strconv.Itoa(item.Port)
  r := []string{cmd}
  switch cmd {
  case "socat":
    r = append(r, "STDIO", fmt.Sprintf("TCP:%s:%s", host, port))
  case "c3270":
    r = append(r, "-model", "3279-2", host, port)
  case "tnz":
    r = append(r, fmt.Sprintf("%s:%s", host, port))
  default:
    r = append(r, host, port)
  }
	return r
}

func Connect(args []string) {
	cmd := exec.Command(args[0],args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if ( err != nil) {
		log.Fatal(err)
	}
}

func DoMenu(bfile string, jfile string) {
	result := ReadJSONMenu(jfile)
	cmdargs := GetMenuReply(bannerfile, result)
	Connect(cmdargs)
}

func main() {
	flag.Parse()
	DoMenu(*bannerfile, *targetfile)
}
