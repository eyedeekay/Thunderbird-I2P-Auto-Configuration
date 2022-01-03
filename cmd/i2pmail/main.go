package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	port      = flag.String("port", "8080", "port to listen on")
	host      = flag.String("host", "i2pmail.org", "host to listen on")
	aliashost = flag.String("aliashost", "mail.i2p", "host to listen on")
	directory = flag.String("directory", "./www", "directory to serve")
)

var ispXML = `<?xml version="1.0" encoding="UTF-8"?>

<clientConfig version="1.1">
	<emailProvider id="mail.i2p">
		<domain>mail.i2p</domain>
		<displayName>Postman's I2P Mail</displayName>
		<displayShortName>Postman's I2P Mail</displayShortName>
		<incomingServer type="pop3">
			<hostname>127.0.0.1</hostname>
			<port>7660</port>
			<socketType>plain</socketType>
			<authentication>password-cleartext</authentication>
			<username>%EMAILADDRESS%</username>
		</incomingServer>
		<outgoingServer type="smtp">
			<hostname>127.0.0.1</hostname>
			<port>7659</port>
			<socketType>plain</socketType>
			<authentication>password-cleartext</authentication>
			<username>%EMAILADDRESS%</username>
		</outgoingServer>
	</emailProvider>
</clientConfig>
`

var ispAliasXML = `<?xml version="1.0" encoding="UTF-8"?>

<clientConfig version="1.1">
	<emailProvider id="i2pmail.org">
		<domain>i2pmail.org</domain>
		<displayName>Postman's I2P Mail</displayName>
		<displayShortName>Postman's I2P Mail</displayShortName>
		<incomingServer type="pop3">
			<hostname>127.0.0.1</hostname>
			<port>7660</port>
			<socketType>plain</socketType>
			<authentication>password-cleartext</authentication>
			<username>%EMAILADDRESS%</username>
		</incomingServer>
		<outgoingServer type="smtp">
			<hostname>127.0.0.1</hostname>
			<port>7659</port>
			<socketType>plain</socketType>
			<authentication>password-cleartext</authentication>
			<username>%EMAILADDRESS%</username>
		</outgoingServer>
	</emailProvider>
</clientConfig>`

// Contains a bunch of copypasta from the top answer here because good enough.
// https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

var ErrOS error = fmt.Errorf("unknown OS: %s and no /etc/hosts file found", runtime.GOOS)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func backupHosts() error {
	hostsFile := hostsFilePath()
	hostsFileBak := hostsFile + ".bak"
	return CopyFile(hostsFile, hostsFileBak)
}

func platformThunderbirdIspPath() string {
	switch runtime.GOOS {
	case "darwin":
		return "/Applications/Thunderbird.app/Contents/MacOS/isp"
	case "linux":
		return "/usr/share/thunderbird/isp"
	case "windows":
		return "C:/Program Files/Mozilla/Thunderbird/isp"
	}
	return "/usr/share/thunderbird/isp"
}

func thunderbirdIspPath() string {
	path := platformThunderbirdIspPath()
	if dirExists(path) {
		return path
	}
	return ""
}

func thunderbirdIspXMLFile() string {
	dir := thunderbirdIspPath()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, *host+".xml")
}

func thunderbirdIspXMLFileAlias() string {
	dir := thunderbirdIspPath()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, *aliashost+".xml")
}

func checkThunderbirdIsp() int {
	ispPath := thunderbirdIspPath()
	if ispPath == "" {
		log.Printf("No Thunderbird ISP directory found")
		return -1
	}
	ispFile := thunderbirdIspXMLFile()
	if !fileExists(ispFile) {
		log.Printf("No Thunderbird ISP file found at %s", ispFile)
		return 1
	}
	return 0
}

func hostsFilePath() string {
	switch runtime.GOOS {
	case "darwin":
		return "/etc/private/hosts"
	case "linux":
		return "/etc/hosts"
	case "windows":
		return "c:\\windows\\system32\\drivers\\etc\\hosts"
	default:
		if fileExists("/etc/hosts") {
			log.Printf("Warning: %s", ErrOS)
			log.Println("/etc/hosts file found. (I assume) it's a unix system")
			return "/etc/hosts"
		}
		return ""
	}
}

func editHosts() error {
	if err := backupHosts(); err != nil {
		return err
	}
	hostsFile := hostsFilePath()
	if hostsFile == "" {
		return ErrOS
	}
	log.Printf("Editing %s...", hostsFile)
	bytes, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		return err
	}
	bytes = append(bytes, []byte(fmt.Sprintf("\n127.0.0.1 %s\n", *host))...)
	return ioutil.WriteFile(hostsFile, bytes, 0644)
}

func checkHosts() bool {
	hostsFile := hostsFilePath()
	if hostsFile == "" {
		log.Fatalf(ErrOS.Error())
	}
	log.Printf("Checking %s...", hostsFile)
	bytes, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		if strings.Contains(line, *host) {
			log.Printf("Hosts file already contains %s", *host)
			return true
		}
	}
	log.Printf("Hosts file does not contain %s", *host)
	return false
}

func uiElevate() (string, []string) {
	if fileExists("/usr/bin/gksudo") {
		return "/usr/bin/gksudo", []string{}
	}
	if fileExists("/usr/bin/kdesudo") {
		return "/usr/bin/kdesudo", []string{}
	}
	if fileExists("/usr/bin/pkexec") {
		return "/usr/bin/pkexec", []string{}
	}
	return "sudo", []string{}
}

func main() {
	flag.Parse()
	restateCommand := []string{os.Args[0], "--host", *host, "--port", *port, "--directory", *directory}
	ispFile := checkThunderbirdIsp()
	if ispFile == 1 {
		if os.Geteuid() != 0 {
			uiElevate, uiArgs := uiElevate()
			log.Printf("Elevating to %s...", "root")
			log.Printf("%s %s %s :%s", uiElevate, uiArgs, restateCommand, thunderbirdIspPath())
			exec.Command(uiElevate, append(uiArgs, restateCommand...)...).Run()
			return
		}
		log.Printf("Creating %s...", thunderbirdIspXMLFile())
		err := ioutil.WriteFile(thunderbirdIspXMLFile(), []byte(ispXML), 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Creating %s...", thunderbirdIspXMLFileAlias())
		err = ioutil.WriteFile(thunderbirdIspXMLFileAlias(), []byte(ispAliasXML), 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else if ispFile == 0 {
		os.Exit(ispFile)
	} else if ispFile == -1 {
		if !checkHosts() {
			if os.Geteuid() == 0 {
				runAs := os.Getenv("SUDO_USER")
				if err := editHosts(); err != nil {
					log.Fatal(err)
				}
				if runAs == "" {
					log.Fatal("SUDO_USER not set")
				}
				exec.Command("sudo", "-u", runAs, os.Args[0], "--host", *host, "--port", *port, "--directory", *directory).Run()
				os.Exit(0)
			} else {
				uiElevate, uiArgs := uiElevate()
				exec.Command(uiElevate, append(uiArgs, os.Args[0], "--host", *host, "--port", *port, "--directory", *directory)...)
				log.Fatal("You must be root to edit the hosts file")
			}
		}
		fs := http.FileServer(http.Dir(*directory))

		address := net.JoinHostPort(*host, *port)
		log.Printf("Listening on %s...", address)
		log.Printf("Serving %s...", *directory)
		log.Printf("Args were %s, %s, %s", *port, *host, *directory)
		err := http.ListenAndServe(address, fs)
		if err != nil {
			log.Fatal(err)
		}
	}
}
