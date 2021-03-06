# Thunderbird I2P Auto Configuration

**NOTE: *Untested on Windows***

Automatic configuration of Mozilla Thunderbird for Postman's
mail services. These services are configured by default on all
current I2P routers, but using them with an email client has
an additional barrier to configuration that other mail services
do not have. Most email clients are aware of email services
and provide some auto-configuration which is detected by domain
name. They achieve this by maintaining a directory of XML files
that describe those services[link](https://wiki.mozilla.org/Thunderbird:Autoconfiguration:MozillaWebservicePublish).
This page presents those same XML files but for Postman's I2P servics.

Using the Files from Disk:
--------------------------

This is the easiest and least disruptive way to install the Thunderbird
ISP file, and the default if it is determined to be possible. The
`i2pmail-OSTYPE-OSARCH` executable, or the `i2pmail-OSTYPE-OSARCH.su3`
plugin will ask for privileges once, and install the `.xml` file to
the default location for Thunderbird for your architecture. If you just
want to use this repository and the `Makefile`, then:

```bash
git clone https://github.com/eyedeekay/Thunderbird-I2P-Auto-Configuration
cd Thunderbird-I2P-Auto-Configuration
sudo make install
```

## Plugin Links

 - [Mac OSX Intel](http://idk.i2p/Thunderbird-I2P-Auto-Configuration/i2pmail-darwin-amd64.su3)
 - [Mac OSX M1](http://idk.i2p/Thunderbird-I2P-Auto-Configuration/i2pmail-darwin-arm64.su3)
 - [FreeBSD amd64](http://idk.i2p/Thunderbird-I2P-Auto-Configuration/i2pmail-freebsd-amd64.su3)
 - [Linux 386](http://idk.i2p/Thunderbird-I2P-Auto-Configuration/i2pmail-linux-386.su3)
 - [Linux amd64](http://idk.i2p/Thunderbird-I2P-Auto-Configuration/i2pmail-linux-amd64.su3)
 - [Linux arm64](http://idk.i2p/Thunderbird-I2P-Auto-Configuration/i2pmail-linux-arm64.su3)
 - [OpenBSD amd64](http://idk.i2p/Thunderbird-I2P-Auto-Configuration/i2pmail-openbsd-amd64.su3)
 - [Windows 386](http://idk.i2p/Thunderbird-I2P-Auto-Configuration/i2pmail-windows-386.su3)
 - [Windows amd64](http://idk.i2p/Thunderbird-I2P-Auto-Configuration/i2pmail-windows-amd64.su3)

Using the Files from Server:
----------------------------

Is not entirely straightforward yet either. You will need to have:

1. An entry in your `hosts` file(`/etc/hosts`) for `i2pmail.org` which
 points at `127.0.0.1` or `localhost`.
2. A web server listening for requests on `i2pmail.org` which is serving
 the xml file under a URL Thunderbird recognizes(`mail/config-v1.1.xml`).

In order to do that on Linux, you can:

```bash
echo 127.0.0.1  i2pmail.org | sudo tee -a /etc/hosts
```

and then configure a server of your choice to listen as i2pmail.org. For
the sake of making things easy, TODO: plugin
