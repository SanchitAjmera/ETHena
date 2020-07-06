import main


class Bot:

    LARGEST_INT

    def __init__(self, tradingPeriod, tradesMade, funds, stock, currRow):
        self.tradingPeriod = tradingPeriod
        self.tradesMade = tradesMade
        self.numOfDecisions = int(50000 / tradingPeriod)
        self.funds = funds
        self.stock = stock
        self.currRow = currRow
        # dictionary of candlestick info
        self.candleStickInfo = {
            "maxAsk":   [0, 0, 0],
            "minAsk":   [LARGEST_INT, LARGEST_INT, LARGEST_INT],
            "maxBid":   [0, 0, 0],
            "minBid":   [LARGEST_INT, LARGEST_INT, LARGEST_INT],
            "openAsk":  [0, 0, 0],
            "openBid":  [0, 0, 0],
            "closeAsk": [0, 0, 0],
            "closeBid": [0, 0, 0]
        }

    def calcCandleStick(self):

    def trade(self):
        currAsk = utils.getAsk(self.currRow)

        self.currRow += 1
