import main
import bots
import utils


def calcMaxMin(asks, bids)


def calcCandleStick(bot):
    asks = []
    bids = []
    for i in range(bot.tradingPeriod):
        asks.append(getAsk(bot.currRow))
        bids.append(getBid(bot.currRow))
    bot.candleStickInfo["openAsk"][2] = asks[0]
    bot.candleStickInfo["openBid"][2] = bids[0]
    bot.candleStickInfo["closeAsk"][2] = asks[bot.tradingPeriod - 1]
    bot.candleStickInfo["closebid"][2] = bids[bot.tradingPeriod - 1]
    bot.candleStickInfo["maxAsk"][2],
    bot.candleStickInfo["maxAsk"][2],
    bot.candleStickInfo["maxAsk"][2],
    bot.candleStickInfo["maxAsk"][2] = calcMaxMin(asks, bids)
