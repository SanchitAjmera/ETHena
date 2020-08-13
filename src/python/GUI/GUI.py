import PySimpleGUI as sg
import os

# make it look nice
sg.change_look_and_feel('Dark Blue 3')

# define the layout of the GUI
layout = [
    [sg.Text('Name', size=(16, 1)), sg.Spin(['devam','luqman','manuj','sanchit','shivam','name'], initial_value='name')],
    [sg.Text('Strategies:', size=(16, 1)), sg.Text('Weightings:')],
    [sg.Text('RSI', size=(16, 1)), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('MACD', size=(16, 1)), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('Candlestick', size=(16, 1)), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('Offset', size=(16, 1)), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('TimeInterval(seconds)', size=(16, 1)), sg.InputText('30')],
    [sg.Text('Trade Live', size=(16, 1)), sg.Checkbox('')],
    [sg.T('main.go file')],
    [sg.In()],
    [sg.FileBrowse(target=(-1, 0))],
    [sg.OK(), sg.Cancel()]]

# form the window
window = sg.Window('ETHena GUI', layout)
# get the values from the GUI

while True:
    button, values = window.read()
    if button == sg.WIN_CLOSED or button == 'Cancel':  # if user closes window or clicks cancel
        break
        exit()
    elif button == 'OK':
        window.close()
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
        command = "go run " + values[7] + " " + name + " " + strategy_chooser + " " + timeinterval + " " + live + ""
        # run the comamnd
        os.system(command)
        break