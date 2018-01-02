###This is a simply program based on beego web framework! It aims to monitor your specified bitcoin sites for any changes. But it only supports BTC and BCH at now, and you can edit this reposity to meet your requirement.

#####Features

- easily deploy(just modify app.conf file and it will work)
- clear debug output(benefit by beego)

######How to use

clone this reposity
```
git clone https://github.com/qshuai/WatchBitcoinAddress.git
```
install dependence
```
cd $GOPATH/src
go get github.com/astaxie/beego
go get github.com/tidwall/gjson
```
run 
```
//if you have install bee tool
bee run

//optional scheme
go build -o tracktx main.go
./tracktx >> debug.log &
```
