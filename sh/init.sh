
mkdir $GOPATH/src/github.com
cd $GOPATH/src/github.com

mkdir $GOPATH/src/golang.org/x
cd $GOPATH/src/golang.org/x
git clone https://github.com/golang/net.git
git clone https://github.com/golang/sync.git
git clone https://github.com/golang/sys.git
git clone https://github.com/golang/tools.git

mkdir $GOPATH/src/google.golang.org
cd $GOPATH/src/google.golang.org

