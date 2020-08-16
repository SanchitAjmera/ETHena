# ETHena
========

ETHena is an algorithmic trading bot which executes trades on the Luno exchange using multiple strategies. This was developed as a solution for the Spark Blockchain Hackathon's Luno trading challenge.

## Features
 - **Luno API:** Integrated functionality from the Luno GO SDK, with ability to fetch live prices of six cryptocurrency paris as well as place and cancel post-limit-orders
 - **Multiple Strategies:** ETHena can trade using information from 4 (user-weighted) strategies including RSI, MACD, Candlesticks, Offset-Ema and Trailing Stoploss for risk managment.
 - **Email Notification:** A convenient manner to update the user on ETHena's trade histories, status, and daily performance summaries.
 - **Performance Report:** Utility feature which automatically generates a graph to help pinpoint where ETHena decided to buy or sell over the course of one day.
 

## Installation
### Dependencies
This system must have the Go language pre-installed. To find out how to install Go, please follow the steps on this page - https://golang.org/doc/install. Python is also used throughout this program, please find the step-by-step guide to install Python here - https://www.python.org/downloads/

To install further dependacies used within the project, navigate through src/go/utils/setup folder and run the file setup.go to install all dependancies:

```go
go run setup.go
```

If you recieve an error while running this file, then you can also links and instructions to manually install the dependancies in src/go/utils/setup/ReadMe.md

### Luno Setup
To access the Luno market, sign up and verify your Luno account here - https://www.luno.com/en/login. Once you've been verified, deposit money into your wallet and initialise an API key. The API key should be kept private as it will provide access to your Luno account remotely. 

Once you've authorised an API key, please insert them into src/go/utils/apiKeys.go in this format:

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

Note: Please ensure the "NAME" is all uppercase.

## Running ETHena
### GUI
Congratulations on completing the setup. To run ETHena go to src/python/GUI and enter the following command:

```python3
python3 GUI.py
```
the ETHena GUI will open and you can enter your trading preferences. Please enter your name which should be the same one you entered into the apiKeys.go file. 

ETHena allows you to weight 4 different strategies from 1 to 9. The weighting of each strategy will determine their importance when making a decision to whether buy or sell. Move the slider to the desired level for each strategy - 0 weighting means it will not be used to execute trades.

Please enter a time interval, which will be how often ETHena decided to execute a trade. The recommended settings are 1 for RSI at 20 second intervals. 

Finally, browse your files to select the main.go file within the src/go/main/ directory and click run to start ETHena.

<p align="center">
  <img src="https://github.com/SanchitAjmera/ETHena/blob/master/docs/images/GUI-Image.png" width="60%">
</p>

You will be directed to ETHena's TUI where you can monitor the ask, bid price and keep track of previous buy and sell orders.

# TODO:
  - Clean up ReadMe
  - Finish acknowledgements in report
  - submit


DONE:
  - Clean up github repo
  - Go through ema and fix
  - Implement ema-offset bot
  - Set up email notifications
  - Email notifications with graphed data
  - Research and develop security for api-keys
  - Prepare submission report, presentation
  
"There's a problem here, a real problem. We should be able to work together"
                                              - Sanchit 3 weeks into the competition
                                              
"If it isn't broken, don't fix it" 
                                              - Shivam 1 week before the deadline
