# ETHena

A trading bot which executes trades on the Luno exchange using multiple strategies. This was developed as a solution for the Spark Blockchain Hackathon's Luno challenge.

## Installation and setup
### dependencies
This system must have the Go language pre-installed. To find out how to install Go, please follow the steps on this page - https://golang.org/doc/install. Python is also used throughout this program, please find the step-by-step guide to install Python here - https://www.python.org/downloads/

To install further dependacies used within the project, navigate through src/go/utils/setup folder and run the file setup.go to install all dependancies:

```go
go run setup.go
```

If you recieve an error while running this file, then you can also links and instructions to manually install the dependancies in src/go/utils/setup/ReadMe.md

### Luno setup
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
this will open up the ETHena GUI where your trading preferences can be entered.

![ETHena GUI](https://github.com/SanchitAjmera/ETHena/blob/master/docs/images/GUI-Image.png | width=100)

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
