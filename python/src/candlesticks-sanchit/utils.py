import main
import bots
import matplotlib.pyplot as plt


def getAsk(currRow):
    return sheet.cell_value(currRow, 7)


def getBid(currRow):
    return sheet.cell_value(currRow, 7)


def plotGraph(ys):
    xs = []
    for i in range(len(ys)):
        xs.append(i)

    plt.plot(xs, ys)
    plt.show()
    plt.close()
