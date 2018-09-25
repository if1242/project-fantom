Install PBC (Ubuntu):

tar -xvf pbc-0.5.14.tar.gz
sudo apt-get install libgmp-dev
sudo apt-get install build-essential flex bison
./configure
make
sudo make install
sudo ldconfig

For test:
go build 1S.go
./1S
