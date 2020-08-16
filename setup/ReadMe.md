## Prerequisites:

- This script is designed for Linux Operating Systems
- Have python3 installed
- Have pip3 installed
- Have Golang installed
- Type the following commands

```shell
chmod +x setup.sh\ 
./setup.sh
```

## Troubleshooting:

Update to the latest python3 version

```shell
sudo apt install python3
```

Update to the latest version of pip3

```shell
sudo apt install python3-pip
```

Update to the latest golang version by replacing the 15 with the latest on golang's website.

```shell
wget https://dl.google.com/go/go1.15.linux-amd64.tar.gz
```

You may inherit some problems with the PySimpleGUI if you do not have the tkinter module

```shell
sudo apt-get install python3-tk
```

One can verify this by running the following command. Again replacing the 15 with the latest version

```
sha256sum go1.15.linux-amd64.tar.gz
```

Make sure go is configured correctly with the following commands

```shell
cd ~
```

Replace 15 with the latest version

```shell
sudo tar -C /usr/local -xzf go1.15.linux-amd64.tar.gz
```

```shell
export PATH=$PATH:/usr/local/go/bin
```

```shell
mkdir ~/go
```

```shell
cd ~/go
```

```git
git clone https://github.com/SanchitAjmera/ETHena.git
```

Now try and run the shell command. If there are still any issues contact any of the develops below:

- **Sanchit Ajmera:** <sanchitajmera2017@gmail.com>
- **Shivam Patel:** <shivpatel1306@gmail.com>
- **Manuj Mishra:** <manujmishra2000@gmail.com>
- **Luqman Liaquat:** <luqman.liaquat90@gmail.com>
- **Devam Savjani:** <devamsavjani@rocketmail.com>
