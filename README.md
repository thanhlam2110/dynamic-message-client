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
