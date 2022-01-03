# Thunderbird I2P Auto Configuration

Automatic configuration of Mozilla Thunderbird for Postman's
mail services. These services are configured by default on all
current I2P routers, but using them with an email client has
an additional barrier to configuration that other mail services
do not have. Most email clients are aware of email services
and provide some auto-configuration which is detected by domain
name. They achieve this by maintaining a directory of XML files
that describe those services[link](https://wiki.mozilla.org/Thunderbird:Autoconfiguration:MozillaWebservicePublish).
This page presents those same XML files but for Postman's I2P servics.

Using the Files:
----------------

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