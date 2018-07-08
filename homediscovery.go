package main

import (
    "os"
    "log"
    "fmt"
    "net"
    "strings"
    "os/exec"
    "encoding/json"
    "io/ioutil"

    "github.com/mostlygeek/arp"
)

func main() {
    mapD := make(map[string]string)
    cmd := "ping"
    args := []string{"-c", "3", "-b", "192.168.1.255"}
    if err := exec.Command(cmd, args...).Run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    print("Completed: Subnet ping ")

    for ip, _ := range arp.Table() {
        host, _ := net.LookupAddr(ip)

        if len(host) > 0 {
            mapD[strings.Split(host[0], ".")[0]] = ip
        } else {
            mapD[ip] = ip 
        }
    }

    print("Completed: Arp table read")
    
    addrs, _ := net.InterfaceAddrs()

    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                localip := ipnet.IP.String()
                host, _ := net.LookupAddr(localip)
                mapD[strings.Split(host[0], ".")[0]] = localip
            }
        }
    }

    mapB, _ := json.MarshalIndent(mapD, "", "    ")
    fmt.Println(string(mapB))

    cachepath := "cache"
    if _, err := os.Stat(cachepath); os.IsNotExist(err) {
        os.Mkdir(cachepath, 0766)
    } 
 
    err := ioutil.WriteFile("cache/devices.cache", []byte(string(mapB)), 0666)
    if err != nil {
        log.Fatal(err)
    }

}
