Install PBC (Ubuntu):

tar -xvf pbc-0.5.14.tar.gz  </br>
sudo apt-get install libgmp-dev sudo ldconfig </br>
sudo apt-get install build-essential flex bison </br>
./configure </br>
make </br>
sudo make install </br>
sudo ldconfig </br>
</br>
For test: </br>
go build 1S.go </br>
./1S </br>
