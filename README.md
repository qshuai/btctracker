## btctracker aims to monitor your specified bitcoin sites for any changes. But it only supports BTC and BCH at now, and you can edit this repository to meet your requirement.

[![Go Report Card](https://goreportcard.com/badge/github.com/qshuai/btctracker)](https://goreportcard.com/report/github.com/qshuai/btctracker)

### Features:

- easily deploy(just modify app.conf file and it will work)
- clear debug output(benefit by beego)

### Usage:

1. clone this reporsity
```
git clone https://github.com/qshuai/btctracker.git
```
2. install dependence
```
cd $GOPATH/src/btctracker.git
glide install
```
3. run: please configurate your `app.conf` file
```
cd ./btctracker
//if you have installed bee tool
bee run

//optional scheme
go build -o tracktx main.go
./tracktx >> debug.log &
```

### Notice:

If you use new BCH address now, please convert it to the older one via following tool online: [https://bch.btc.com/tools/address-converter](https://bch.btc.com/tools/address-converter)
