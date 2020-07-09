import main
import bots
import xlrd
import plotly.graph_objects as go
import matplotlib.pyplot as plt


# global variable
wb = xlrd.open_workbook(
    "../../../go/src/ticker/data_7to8_July/tickerData09072020.xlsx")
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

    fig = go.Figure(data=[go.Candlestick(x=xs,
                                         open=bot.candleStickInfo["openAsk"][2],
                                         high=bot.candleStickInfo["maxAsk"][2],
                                         low=bot.candleStickInfo["minAsk"][2],
                                         close=bot.candleStickInfo["closeAsk"][2])])

    fig.show()
    plt.plot(xs, ys)
    plt.plot(bot.buys[0], bot.buys[1], 'o', color="orange")
    plt.show()
    plt.close()
