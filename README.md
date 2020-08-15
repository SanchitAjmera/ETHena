# ETHena

ETHena is a bot which executes trades on the Luno exchange based on multiple strategies. This was developed as a solution for the Spark Blockchain Hackathon's Luno challenge.

## Installation and setup

This system must have the Go language pre-installed. To find out how to install Go, please follow the steps on this page - https://golang.org/doc/install. Python is also used throughout this program, please find the step-by-step guide to install Python here - https://www.python.org/downloads/

To install further dependacies used within the project, navigate through src -> go -> utils -> setup folder and run the file setup.go to install all dependancies:

'''
go run setup.go
'''

If you recieve an error while running this file, then you can also links to manually install the dependancies in src -> go -> util -> setup -> ReadMe.md

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
