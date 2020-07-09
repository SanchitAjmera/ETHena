
import matplotlib.pyplot as plt
import bots
import utils
import candlesticks


def test(bot):
    for i in range(bot.numOfDecisions):
        bot.trade()


def rotate(bot):
    for key in bot.candleStickInfo:
        bot.candleStickInfo[key][0] = bot.candleStickInfo[key][1]
        bot.candleStickInfo[key][1] = bot.candleStickInfo[key][2]


def main():
    tradingPeriod = 10
    bot = bots.Bot(tradingPeriod, 0, 100, 0, tradingPeriod + 1)

    candlesticks.calcCandleStick(bot)
    rotate(bot)
    candlesticks.calcCandleStick(bot)
    rotate(bot)
    candlesticks.calcCandleStick(bot)
    test(bot)

    ys = []
    # print(bot.sells)
    # print(bot.buys)
    for i in range(1, 50000):
        ys.append(utils.getAsk(i))

    utils.plotGraph(ys, bot)

    print("trades made :", len(bot.buys[0]) + len(bot.sells[0]))
    print(" ")
    print("profit: ", bot.funds + (bot.stock * ys[len(ys) - 1]) - 100)


if __name__ == "__main__":
    main()
