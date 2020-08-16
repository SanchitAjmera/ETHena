import main
import candlesticks
import utils


class Bot:

    def __init__(self, tradingPeriod, tradesMade, funds, stock, currRow):
        LARGEST_INT = 1844674407370955200

        self.tradingPeriod = tradingPeriod
        self.tradesMade = tradesMade
        self.numOfDecisions = int(800 / tradingPeriod)
        self.funds = funds
        self.stock = stock
        self.currRow = currRow
        self.buys = [[], []]
        self.sells = [[], []]
        self.candleSticks = [[], [], [], [], [], [], [], []]
        self.stopLoss = 0
        # dictionary of candlestick info
        self.candleStickInfo = {
            "maxAsk":   [0, 0, 0],
            "minAsk":   [LARGEST_INT, LARGEST_INT, LARGEST_INT],
            "openAsk":  [0, 0, 0],
            "closeAsk": [0, 0, 0],
            "maxBid":   [0, 0, 0],
            "minBid":   [LARGEST_INT, LARGEST_INT, LARGEST_INT],
            "openBid":  [0, 0, 0],
            "closeBid": [0, 0, 0]
        }

    def buy(self):
        currAsk = utils.getAsk(self.currRow)
        buyable = self.funds / currAsk
        if self.funds == 0:
            return
        if self.funds >= buyable * currAsk:
            self.stock += buyable
            self.funds -= currAsk * buyable
            self.buys[0].append(self.currRow - 1)
            self.buys[1].append(currAsk)
            self.stopLoss = currAsk
            print("BUY AT ", currAsk)
            print("funds: ", self.funds)
            print("stock: ", self.stock)
            print(" ")

    def sell(self):
        currAsk = utils.getAsk(self.currRow)

        if self.stock == 0:  # or self.stopLoss > currAsk:
            return
        else:
            self.funds += self.stock * currAsk
            self.stock = 0
            self.sells[0].append(self.currRow - 1)
            self.sells[1].append(utils.getAsk(self.currRow - 1))
            print("SELL AT ", currAsk)
            print("funds: ", self.funds)
            print("stock: ", self.stock)
            print(" ")

    def prediction(self, info):
        cond1 = info["maxAsk"][1] > info["maxAsk"][0]
        cond2 = info["maxAsk"][1] > info["maxAsk"][2]
        cond3 = info["minAsk"][1] > info["minAsk"][0]
        cond4 = info["minAsk"][1] > info["minAsk"][2]
        cond5 = info["closeAsk"][0] > info["openAsk"][0]
        cond6 = info["openAsk"][2] > info["closeAsk"][2]
        cond7 = info["closeAsk"][1] > info["openAsk"][1]

        if cond1 and cond2 and cond3 and cond4 and cond5 and cond6 and cond7:
            return "sell"
        elif not cond1 and not cond2 and not cond3 and not cond4 and not cond5 and not cond6 and not cond7:
            return "buy"
        else:
            return ""

    def trade(self):
        i = 0
        for key in self.candleStickInfo:
            self.candleSticks[i].append(self.candleStickInfo[key][0])
            self.candleStickInfo[key][0] = self.candleStickInfo[key][1]
            self.candleStickInfo[key][1] = self.candleStickInfo[key][2]
            i += 1

        candlesticks.calcCandleStick(self)
        predict = self.prediction(self.candleStickInfo)
        if predict == "sell":
            self.sell()
        elif predict == "buy":
            self.buy()
        # print(predict)
        self.currRow += self.tradingPeriod
