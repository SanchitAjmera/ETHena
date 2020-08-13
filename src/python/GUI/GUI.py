import PySimpleGUIWeb as sg
import os

# initialise the GUI
form = sg.FlexForm('ETHena GUI')  # begin with a blank form

# make it look nice
sg.theme('Default')

# define the layout of the GUI
layout = [
    [sg.Text(
        'Please enter your name and use the sliders to adjust the weighting of each strategy. Set the slider to zero to not use the strategy. Then add the time period you wish to trade upon. Finally add if you would like to trade live or on historical data')],
    [sg.Text('Name', size=(15, 1)), sg.InputText('name')],
    [sg.Text('Strategies:')],
    [sg.Text('RSI'), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('MACD'), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('Candlestick'), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('Offset'), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('TimeInterval(seconds)'), sg.InputText('30')],
    [sg.Checkbox('Trade Live')],
    [sg.Submit(), sg.Cancel()]]

# get the values from the GUI
button, values = form.Layout(layout).Read()

# if cancel button was pressed end code
if button == 'Cancel':
    exit()

# initialise the binary number to choose the strategy
strategy_chooser = ''
# get the number from the tick boxes
for i in range(1, 5, 1):
    if values[i]:
        strategy_chooser = strategy_chooser + '1'
    else:
        strategy_chooser = strategy_chooser + '0'
name = values[0]
timeinterval = values[5]
live = values[6]
if live:
    live = "1"
else:
    live = "0"
# form the command
command = "nohup go run TradingHackathon/src/go/main/main.go " + name + " " + strategy_chooser + " " + timeinterval + " " + live + "& "
# run the comamnd
os.system(command)
