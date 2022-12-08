# sshspray 

## SSH bruteforce tool.

Sshspray is a penetration testing tool for checking SSH access with a set of credentials.

You can supply individual credentials of files of hosts/users/passwords to test across a range of hosts.

## Flags:
#
-l&emsp;&emsp;&nbsp;&nbsp;User | *or users file*

-p&emsp;&emsp;&nbsp;Password | *or passwords file*

-t&emsp;&emsp;&nbsp;&nbsp;Host | *or hosts file*

#
### Optional flags:
-x&emsp;&emsp;&nbsp;Execute a command on the host upon connection | *default=False*

-d&emsp;&emsp;&nbsp;Delay between connection attempts | *default=10ms*
#
--port&nbsp;&nbsp;Port | *default=22*
#

### Examples:
#
```
./sshspray spray -l users.txt -p passwords.txt -t hosts.txt
```

```
./sshspray spray -l admin -P passwords.txt -t 127.0.0.1 -d 100 --port 82
```

```
./sshspray spray -l users.txt -p password123 -t hosts.txt -d 100 -x whoami
```
#