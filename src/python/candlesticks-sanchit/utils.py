import main
import bots
import xlrd
import plotly.graph_objects as go
import matplotlib.pyplot as plt
import plotly.graph_objects as go


# global variable
wb = xlrd.open_workbook(
<<<<<<< HEAD
    "../../../go/src/ticker/data_7to8_July/tickerData09072020.xlsx")
=======
    "../../../go/src/backtesting-sanchit/recentAPIdata2018v2.xlsx")
>>>>>>> c69df557ad4a17931d1064a37e86ec44bb0e021a
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

<<<<<<< HEAD
    fig = go.Figure(data=[go.Candlestick(x=xs,
                                         open=bot.candleStickInfo["openAsk"][2],
                                         high=bot.candleStickInfo["maxAsk"][2],
                                         low=bot.candleStickInfo["minAsk"][2],
                                         close=bot.candleStickInfo["closeAsk"][2])])

    fig.show()
=======
    for i in range(bot.numOfDecisions):
        xs2.append(i)

    fig = go.Figure(data=[go.Candlestick(x=xs2,
                                         open=bot.candleSticks[2],
                                         high=bot.candleSticks[0],
                                         low=bot.candleSticks[1],
                                         close=bot.candleSticks[3])])

    fig.show()

>>>>>>> c69df557ad4a17931d1064a37e86ec44bb0e021a
    plt.plot(xs, ys)
    plt.plot(bot.buys[0], bot.buys[1], 'o', color="orange")
    plt.plot(bot.sells[0], bot.sells[1], 'o', color="green")
    plt.show()
    plt.close()
