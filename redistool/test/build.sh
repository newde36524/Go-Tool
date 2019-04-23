mkdir release
cd client01
go build 
cp client01.exe ../release
rm client01.exe
cd ../
cd client02
go build 
cp client02.exe ../release
rm client02.exe
cd ../
cd client03
go build 
cp client03.exe ../release
rm client03.exe
cd ../
cd master
go build 
cp master.exe ../release
rm master.exe
