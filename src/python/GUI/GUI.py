import PySimpleGUI as sg
import os

# initialise the GUI
form = sg.FlexForm('Simple data entry form')  # begin with a blank form

# make it look nice
sg.theme('Default')

# define the layout of the GUI
layout = [
    [sg.Text(
        'Please enter your name and tick any strategies you would like to use. Finally adjust the slider based on what risk you are willing to take')],
    [sg.Text('Name', size=(15, 1)), sg.InputText('name')],
    [sg.Text('Strategies:')],
    [sg.Checkbox('RSI')],
    [sg.Checkbox('MACD')],
    [sg.Checkbox('Candlestick')],
    [sg.Checkbox('Offset')],
    [sg.Text('Risk'), sg.Slider(range=(1, 100), orientation='v', size=(5, 20), default_value=50)],
    [sg.Submit(), sg.Cancel()]
]

# get the values from the GUI
button, values = form.Layout(layout).Read()

# initialise the binary number to choose the strategy
strategy_chooser = ''
# get the number from the tick boxes
for i in range(2, 6, 1):
    if values[i]:
        strategy_chooser = strategy_chooser + '1'
    else:
        strategy_chooser = strategy_chooser + '0'
# form the command
command = "nohup go run TradingHackathon/src/go/main/main.go " + values[0] + strategy_chooser + "&"
# run the comamnd
os.system(command)
