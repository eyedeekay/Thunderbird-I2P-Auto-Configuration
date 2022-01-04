
echo:

include plugin.mk

index:
	@echo "<head>" > index.html
	@echo "<title>Makefile</title>" >> index.html
	@echo "<link rel="stylesheet" type="text/css" href="/style.css" />" >> index.html
	@echo "</head>" >> index.html
	@echo "<body>" >> index.html
	pandoc README.md >> index.html >> index.html
	@echo "<h2>XML Files</h2>" >> index.html
	@echo "<h3>mail.i2p.xml</h3>" >> index.html
	@echo "<pre>" >> index.html
	@cat mail.i2p.xml >> index.html
	@echo "</pre>" >> index.html
	@echo "<h3>i2pmail.org.xml</h3>" >> index.html
	@echo "<pre>" >> index.html
	@cat i2pmail.org.xml >> index.html
	@echo "</pre>" >> index.html
	@echo "</body>" >> index.html
	cp index.html conf/www/
	cp ../style.css conf/www/

install:
	cp mail.i2p.xml /usr/share/thunderbird/isp/mail.i2p.xml
	cp i2pmail.org.xml /usr/share/thunderbird/isp/i2pmail.org.xml

uninstall:
	rm /usr/share/thunderbird/isp/mail.i2p.xml /usr/share/thunderbird/isp/i2pmail.org.xml