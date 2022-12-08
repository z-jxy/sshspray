package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"sshspray/spray"

	"github.com/spf13/cobra"
)

var _hosts string
var _users string
var _passwords string
var _delay int
var _cmd string

var _sshport int

var cmdSpray = &cobra.Command{
	Use:   "spray",
	Short: "spray",
	Long:  `Test SSH access with a set of credentials`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var _x string
		var _y int

		if strings.Contains(_hosts, ":") {
			target := strings.Split(_hosts, ":")
			_x = target[0]
			_y, _ = strconv.Atoi(target[1])
		} else {
			_x = _hosts
			_y = _sshport
		}

		host := _x
		port := _y
		user := _users
		pass := _passwords

		_input := spray.SprayModule{
			Hosts:     host,
			Users:     user,
			Passwords: pass,
		}
		fmt.Printf("Spray Settings: \n\nHost: %s \nPort: %d \nUser: %s \nPass: %s \n", host, port, user, pass)

		if len(_cmd) > 0 {
			fmt.Printf("Execute command: '%s' \n", _cmd)
		} else {
			fmt.Println("Execute command: False")
		}

		val := strconv.Itoa(_delay) + "ms"

		delay, _ := time.ParseDuration(val)
		fmt.Printf("Delay: %v\n", delay)

		spray.InitSshTask(_input, delay, _cmd, _sshport)
	},
}
