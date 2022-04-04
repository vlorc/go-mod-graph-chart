package gosrc

import (
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func CheckTcpConnect(host string, port string) (err error) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return
	}
	if conn != nil {
		defer conn.Close()
		return
	}
	return
}

func GetUnUsePort() uint32 {
	for i := 0; i < 10; i++ {
		if i < 1024 {
			continue
		}
		var newPort = rand.Intn(65535)
		if err := CheckTcpConnect("127.0.0.1", strconv.Itoa(newPort)); nil != err {
			return uint32(newPort)
		}
	}
	return 0
}

// browsers returns a list of commands to attempt for web visualization.
func browsers() []string {
	var cmds []string
	if userBrowser := os.Getenv("BROWSER"); userBrowser != "" {
		cmds = append(cmds, userBrowser)
	}
	switch runtime.GOOS {
	case "darwin":
		cmds = append(cmds, "/usr/bin/open")
	case "windows":
		cmds = append(cmds, "cmd /c start")
	default:
		// Commands opening browsers are prioritized over xdg-open, so browser()
		// command can be used on linux to open the .svg file generated by the -web
		// command (the .svg file includes embedded javascript so is best viewed in
		// a browser).
		cmds = append(cmds, []string{"chrome", "google-chrome", "chromium", "firefox", "sensible-browser"}...)
		if os.Getenv("DISPLAY") != "" {
			// xdg-open is only for use in a desktop environment.
			cmds = append(cmds, "xdg-open")
		}
	}
	return cmds
}

// open browser with url
func OpenBrowser(targetUrl string) (err error) {
	// Construct URL.
	u, _ := url.Parse(targetUrl)

	for _, b := range browsers() {
		args := strings.Split(b, " ")
		if len(args) == 0 {
			continue
		}
		viewer := exec.Command(args[0], append(args[1:], u.String())...)
		viewer.Stderr = os.Stderr
		if err = viewer.Start(); err == nil {
			return
		}
	}
	// No visualizer succeeded, so just print URL.
	fmt.Println(u.String())
	return
}
