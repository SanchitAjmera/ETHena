import main
import bots
import xlrd
import matplotlib.pyplot as plt


# global variable
wb = xlrd.open_workbook("recentAPIdata.xlsx")
sheet = wb.sheet_by_index(0)


def getAsk(currRow):
    res = (sheet.cell_value(currRow, 7))

    if res == "NaN":
        return getBid(currRow - 1)
    else:
        return float(res)


def getBid(currRow):
    res = (sheet.cell_value(currRow, 7))

    if res == "NaN":
        return getBid(currRow - 1)
    else:
        return float(res)


def plotGraph(ys, bot):
    xs = []
    for i in range(len(ys)):
        xs.append(i)

    plt.plot(xs, ys)
    plt.plot(bot.buys[0], bot.buys[1], 'o', color="orange")
    plt.show()
    plt.close()
