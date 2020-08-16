#!/bin/bash
clear

echo "Installing PySimpleGUI"
pip3 install PySimpleGUI

clear
echo "PySimpleGUI installed"


echo "Installing pandas"
pip3 install pandas

clear
echo "pandas installed"

echo "Installing matplotlib"
pip3 install matplotlib

clear
echo "matplotlib installed"

echo "Installing datetime"
pip3 install datetime

clear
echo "datetime installed"

echo "Installing excelize"
go get github.com/360EntSecGroup-Skylar/excelize

clear
echo "excelize installed"

echo "Installing tealeg"
go get github.com/tealeg/xlsx

clear
echo "tealeg installed"

echo "Installing gomail"
go get gopkg.in/gomail.v2

clear
echo "gomail installed"

echo "Installing Luno"
go get github.com/luno/luno-go

clear
echo "Luno installed"

echo "All Packages Installed"