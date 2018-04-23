## This is a simply program based on beego web framework! It aims to monitor your specified bitcoin sites for any changes. But it only supports BTC and BCH at now, and you can edit this repository to meet your requirement.

[![Go Report Card](https://goreportcard.com/badge/github.com/qshuai/WatchBitcoinAddress)](https://goreportcard.com/report/github.com/qshuai/WatchBitcoinAddress)

### Features

- easily deploy(just modify app.conf file and it will work)
- clear debug output(benefit by beego)

### How to use

1. clone this reporsity
```
git clone https://github.com/qshuai/WatchBitcoinAddress.git
```
2. install dependence
```
cd $GOPATH/src/WatchBitcoinAddress.git
glide install
```
3. run 
```
cd ./WatchBitcoinAddress
//if you have installed bee tool
bee run

//optional scheme
go build -o tracktx main.go
./tracktx >> debug.log &
```

### Notice:

If you use new BCH address now, please convert it to the older one via following tool online: [https://bch.btc.com/tools/address-converter](https://bch.btc.com/tools/address-converter)
