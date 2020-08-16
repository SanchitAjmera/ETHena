import PySimpleGUI as sg
import os

# make it look nice
sg.change_look_and_feel('DefaultNoMoreNagging')

text_size = (22, 3)
drop_down_menu_size = (21, 3)
slider_size = (20.25, 14)
box_size = (23, 3)

window = sg.Window('Columns')
# define the layout of the GUI
pre_column = [sg.Image('Logo.png', pad=(200, 0))]
col1 = [
    [sg.Text('Name', size=text_size),
     sg.Combo(['devam', 'luqman', 'manuj', 'sanchit', 'shivam'], size=drop_down_menu_size)],
    [sg.Text('Strategies:', size=text_size), sg.Text('Weightings:')],
    [sg.Text('   RSI', size=text_size), sg.Slider(range=(0, 9), orientation='h', default_value=0, size=slider_size)],
    [sg.Text('   MACD', size=text_size), sg.Slider(range=(0, 9), orientation='h', default_value=0, size=slider_size)],
    [sg.Text('   Candlestick', size=text_size),
     sg.Slider(range=(0, 9), orientation='h', default_value=0, size=slider_size)],
    [sg.Text('   Offset', size=text_size), sg.Slider(range=(0, 9), orientation='h', default_value=0, size=slider_size)]]

col2 = [
    [sg.Text('TimeInterval(seconds)', size=text_size), sg.InputText('30', size=box_size)],
    [sg.Text('', size=text_size), sg.Radio('Live', "RADIO1", default=True),
     sg.Radio("Offline", "RADIO1")],
    [sg.T('Path to main.go file', size=text_size),
     sg.In(default_text='Click browse button', size=box_size)],
    [sg.Text('', size=text_size), sg.FileBrowse(target=(-1, 1))]]

last_line = [sg.Text('', size=(83, 0)), sg.OK(button_text="Run"), sg.Cancel()]

# form the window
layout = [pre_column,
          [sg.Column(col1),
           sg.VerticalSeparator(pad=None, color='Black'),
           sg.Column(col2)],
          last_line]
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
        for i in range(2, 6, 1):
            strategy_chooser = strategy_chooser + str(int(values[i]))
        name = values[1]
        timeinterval = values[7]
        live = values[8]
        if live:
            live = "1"
        else:
            live = "0"
        # form the command
        command = "go run " + values[10] + " " + name + " " + strategy_chooser + " " + timeinterval + " " + live + ""
        # run the comamnd
        os.system(command)
        break
