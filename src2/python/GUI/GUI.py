import PySimpleGUI as sg
import os

# make it look nice
sg.change_look_and_feel('Dark Blue 3')

text_size = (20, 3)

# define the layout of the GUI
layout = [
    [sg.Text('Name', size=text_size),
     sg.Combo(['devam', 'luqman', 'manuj', 'sanchit', 'shivam'], size=(19, 3), pad=(0, 1))],
    [sg.Text('Strategies:', size=text_size), sg.Text('Weightings:')],
    [sg.Text('RSI', size=text_size), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('MACD', size=text_size), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('Candlestick', size=text_size), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('Offset', size=text_size), sg.Slider(range=(0, 9), orientation='h', default_value=0)],
    [sg.Text('TimeInterval(seconds)', size=text_size), sg.InputText('30')],
    [sg.Text('', size=text_size), sg.Radio('Live', "RADIO1", default=True, size=(19, 3)),
     sg.Radio("Offline", "RADIO1")],
    [sg.T('Path to main.go file', size=text_size), sg.In()],
    [sg.Text('', size=text_size), sg.FileBrowse(target=(-1, 1))],
    [sg.OK(button_text="Run"), sg.Cancel()]]

# form the window
window = sg.Window('ETHena', layout)
# get the values from the GUI

while True:
    button, values = window.read()
    if button == sg.WIN_CLOSED or button == 'Cancel':  # if user closes window or clicks cancel
        break
        exit()
    elif button == 'Run':
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
        command = "go run " + values[8] + " " + name + " " + strategy_chooser + " " + timeinterval + " " + live + ""
        # run the comamnd
        os.system(command)
        break
