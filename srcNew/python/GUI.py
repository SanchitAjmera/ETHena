import tkinter as tk

events_list = []

# Create an event handler


def handle_keypress(event):
    """Print the character associated to the key pressed"""
    print(event.char)


window = tk.Tk()

window.title("QuantX")

greeting = tk.Label(text="Hello, Tkinter")

greeting.pack()

label = tk.Label(
    text="add API key ID",
    fg="white",
    bg="black",
    height=10,
    width=10)


def handle_click(event):
    window.destroy()


button = tk.Button(
    text="RSI",
    width=25,
    height=5,
    bg="blue",
    fg="yellow",
)

button.bind("<Button-1>", handle_click)


entry = tk.Entry(fg="yellow", bg="red", width=50)

input = entry.get()

entry.pack()
label.pack()
button.pack()

# Bind keypress event to handle_keypress()
window.bind("<Key>", handle_keypress)

window.mainloop()

print("The button was clicked!")
