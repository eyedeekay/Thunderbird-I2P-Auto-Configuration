
echo:

include plugin.mk

index:
	@echo "<head>"
	@echo "<title>Makefile</title>"
	@echo "<link rel="stylesheet" type="text/css" href="/style.css" />"
	@echo "</head>"
	@echo "<body>"
	markdown README.md >> index.html
	@echo "<h2>XML Files</h2>"
	@echo "<h3>mail.i2p.xml</h3>"
	@echo "<pre>"
	@echo -n "`cat mail.i2p.xml`"
	@echo "</pre>"
	@echo "<h3>i2pmail.org.xml</h3>"
	@echo "<pre>"
	@echo -n "`cat i2pmail.org.xml`"
	@echo "</pre>"
	@echo "</body>"
	cp index.html conf/www/
	cp ../style.css conf/www/

install:
	cp mail.i2p.xml /usr/share/thunderbird/i2p/mail.i2p.xml
	cp i2pmail.org.xml /usr/share/thunderbird/i2p/i2pmail.org.xml

uninstall:
	rm /usr/share/thunderbird/i2p/mail.i2p.xml /usr/share/thunderbird/i2p/i2pmail.org.xml