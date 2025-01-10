# Port Scanner

- **Author: Rafael Chavez Solis**
- ***Email: rafaelchavezsolis@gmail.com***

## Description

This is a command line application for port scanning. For a range of port numbers of a given ip address, this program will try to establish a TCP connection. If the connection its successfully established the port will be consider as open.

This program use command line arguments and flags to parse the parameters, also creates a pool of workers to scan ports concurrently and use mutexes for store results

## How to use

For scanning use the command line argument to set the ip address or hostname and the flags --from and --to to set the port range numbers.

<binary_filename> -f(--from)=<start_port> -t(--to)=<end_port> <ipaddress_or_host>

**Example**
```bash
 ./bin/main  -f 10 -t 500 scanme.nmap.org
```

## Tech Stack

- Go 1.23.4
