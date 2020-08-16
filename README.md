
<img src="https://github.com/SanchitAjmera/ETHena/blob/master/docs/images/ETHenaLogo.png" width="30%">
 

ETHena is an algorithmic trading bot which executes trades on the Luno exchange using multiple strategies. This was developed as a solution for the Spark Blockchain Hackathon's Luno trading challenge.

## Features
 - **Luno API:** Integrated functionality from the Luno GO SDK, with the ability to fetch live prices of six cryptocurrency paris as well as place or cancel post-limit-orders
 - **Multiple Strategies:** ETHena can trade using information from 4 (user-weighted) strategies including RSI, MACD, Candlesticks, Offset-Ema and Trailing Stoploss for risk managment.
 - **Email Notification:** A convenient manner to update the user on ETHena's trade histories, status, and daily performance summaries.
 - **Performance Report:** Utility feature which automatically generates a graph to help pinpoint where ETHena decided to buy or sell over the course of one day.
 

## Installation
### Dependencies
This system is supported the Go language and Python 3.8 . To find out how to install Go, and Python 3.8 follow the step-by-step guides here - https://golang.org/doc/install, https://www.python.org/downloads/

To install further dependacies used within the project, navigate through the ```src/go/utils/setup``` folder and run the file ```setup.go``` to install all dependancies:

```go
go run setup.go
```

**Note:** If you recieve an error while running this file, then follow the links and instructions to manually install the dependancies in ```src/go/utils/setup/ReadMe.md```

### Luno Setup
To access the Luno market, sign up and verify your Luno account here - https://www.luno.com/en/login. After you've been verified, deposit money into your wallet and initialise an API key. The API key should be kept private as it will provide access to your Luno account remotely. 

Insert the authorised API key and key ID into ```src/go/utils/apiKeys.go``` in this format:

**Note:** Please ensure the ```<NAME>``` is all uppercase.

```go
func InitialiseKeys() {
	ApiKeys = map[string]([]string){
		"<NAME>": []string{
			"<KEY_ID>",
			"<KEY>",
		},
	}
}
```

## Quickstart
### GUI
Congratulations on completing the setup. To run ETHena, go to ```src/python/GUI``` and enter the following command:

```python3
python3 GUI.py
```
The ETHena GUI will open and you can enter your trading preferences.
 - Enter your name which should be the same one you entered into the apiKeys.go file.
 - ETHena allows you to weight 4 different strategies from 1 to 9. The weighting of each strategy will determine their importance when making a decision to whether buy or sell. Move the slider to the desired level for each strategy - 0 weighting means it will not be used to execute trades.
 - Specify a time interval, which will be how often ETHena decided to execute a trade. The recommended settings are 1 for RSI at 20 second intervals. 
 - Finally, browse your files to select the ```main.go``` file within the ```src/go/main/``` directory and click run to start ETHena.

<p align="center">
  <img src="https://github.com/SanchitAjmera/ETHena/blob/master/docs/images/GUI-Image.png" width="60%">
</p>

You will be directed to ETHena's TUI where you can monitor the ask, bid price and keep track of previous buy and sell orders.

## Questions?
If you don't understand something or find an issue with the program, please create an issue for this repository on Github or contact one of the contributors - our contact details are:
 - **Sanchit Ajmera:** <sanchitajmera2017@gmail.com>
 - **Shivam Patel:** <shivpatel1306@gmail.com>
 - **Manuj Mishra:** <manujmishra2000@gmail.com>
 - **Luqman Liaquat:** <luqman.liaquat90@gmail.com>
 - **Devam Savjani:** <devamsavjani@rocketmail.com>
