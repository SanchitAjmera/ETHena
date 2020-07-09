import main
import bots
import utils


def calcMaxMin(asks, bids):
    LARGEST_INT = 1844674407370955200
    maxBid = 0
    maxAsk = 0
    minBid = LARGEST_INT
    minAsk = LARGEST_INT
    for i in range(len(asks)):
        if asks[i] > maxAsk:
            maxAsk = asks[i]
        if asks[i] < minAsk:
            minAsk = asks[i]
        if bids[i] > maxBid:
            maxBid = bids[i]
        if bids[i] < minBid:
            minBid = bids[i]
    return maxAsk, minAsk, maxBid, minBid


def calcCandleStick(bot):
    asks = []
    bids = []
    for i in range(bot.tradingPeriod):
        asks.append(utils.getAsk(bot.currRow - i))
        bids.append(utils.getBid(bot.currRow - i))

    bot.candleStickInfo["closeAsk"][2] = asks[0]
    bot.candleStickInfo["closeBid"][2] = bids[0]
    bot.candleStickInfo["openAsk"][2] = asks[bot.tradingPeriod - 1]
    bot.candleStickInfo["openBid"][2] = bids[bot.tradingPeriod - 1]
    (bot.candleStickInfo["maxAsk"][2],
     bot.candleStickInfo["minAsk"][2],
     bot.candleStickInfo["maxBid"][2],
     bot.candleStickInfo["minBid"][2]) = calcMaxMin(asks, bids)
