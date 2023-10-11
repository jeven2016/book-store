
### install
```shell
go get github.com/gocolly/colly/v2 latest
msgp 序列化
```
https://github.com/prometheus/client_golang/blob/main/prometheus/examples_test.go
https://github.com/lao-siji/lao-siji

``
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
``

### Build
Using taskfile to build

```shell
# install 
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

chmod -R 777 ./bin/task
mv ./bin/task /usr/local/bin/

task -version

```


libs
* github.com/go-co-op/gocron
* chromedp(浏览器抓取)