# dynamic-message-client
gRPC Dynamic Message client
## Requires
`-OS: Ubuntu 16.04`

`-Go version 15`
## Install Go version 15
### Install
~~~
root@ubuntu:~# sudo apt-get update
root@ubuntu:~# apt-get upgrade
root@ubuntu:~# wget https://dl.google.com/go/go1.15.6.linux-amd64.tar.gz
root@ubuntu:~# sudo tar -xvf go1.15.6.linux-amd64.tar.gz
root@ubuntu:~# sudo mv go /usr/local
~~~
### Config
~~~
root@ubuntu:~# mkdir admin
root@ubuntu:~# cd admin/
root@ubuntu:~/admin# mkdir go
root@ubuntu:/# nano ~/.profile
~~~
ADD at End of file ~/.profile
~~~
export GOROOT=/usr/local/go
export GOPATH=$HOME/admin/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
root@ubuntu:/# source ~/.profile
~~~
Check version
~~~
root@ubuntu:/# go version
~~~
### Test
~~~
root@master:~# cd admin/go/
root@master:~/admin/go# mkdir src
root@master:~/admin/go# cd src/
root@master:~/admin/go/src# mkdir github.com
root@master:~/admin/go/src# cd github.com/
root@master:~/admin/go/src/github.com# mkdir thanhlam
root@master:~/admin/go/src/github.com# cd thanhlam/
root@master:~/admin/go/src/github.com/thanhlam#
~~~
Clone source code in directory thanhlam
~~~
Get dependencies
root@master:~/admin/go/src/github.com/thanhlam/grpc-client# go get .
Run code
root@master:~/admin/go/src/github.com/thanhlam/grpc-client# go run .
~~~
