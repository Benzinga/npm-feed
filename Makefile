cross: \
	npmfeed-windows.zip \
	npmfeed-darwin.tar.gz \
	npmfeed-freebsd.tar.gz \
	npmfeed-linux.tar.gz

npmfeed-windows.exe:
	GOOS=windows GOARCH=amd64 go build -o $@

npmfeed-windows.zip: npmfeed-windows.exe
	zip $@ $<

npmfeed-darwin:
	GOOS=darwin GOARCH=amd64 go build -o $@

npmfeed-darwin.tar.gz: npmfeed-darwin
	tar -cvzf $@ $<

npmfeed-freebsd:
	GOOS=freebsd GOARCH=amd64 go build -o $@

npmfeed-freebsd.tar.gz: npmfeed-freebsd
	tar -cvzf $@ $<

npmfeed-linux:
	GOOS=linux GOARCH=amd64 go build -o $@

npmfeed-linux.tar.gz: npmfeed-linux
	tar -cvzf $@ $<
