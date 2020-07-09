import main
import bots
import xlrd
import matplotlib.pyplot as plt
import plotly.graph_objects as go


# global variable
wb = xlrd.open_workbook(
    "../../../go/src/backtesting-sanchit/recentAPIdata2018v2.xlsx")
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
    xs2 = []
    for i in range(len(ys)):
        xs.append(i)

    for i in range(bot.numOfDecisions):
        xs2.append(i)

    fig = go.Figure(data=[go.Candlestick(x=xs2,
                                         open=bot.candleSticks[2],
                                         high=bot.candleSticks[0],
                                         low=bot.candleSticks[1],
                                         close=bot.candleSticks[3])])

    fig.show()

    plt.plot(xs, ys)
    plt.plot(bot.buys[0], bot.buys[1], 'o', color="orange")
    plt.plot(bot.sells[0], bot.sells[1], 'o', color="green")
    plt.show()
    plt.close()
