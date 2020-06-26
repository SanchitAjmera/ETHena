import process
import analysis


class Trader:

    def __init__(self, initialFund, initialInventory, givenData, testData):
        self.funds = initialFund
        self.inventory = initialInventory
        self.historicalData = givenData
        self.testData = testData

    def buy(self):
        # TODO: need a mechanism for working out how much to buy
        price = self.testData[1][0]
        buyableStock = self.funds / price

        if buyableStock != 0:
            self.inventory.append([buyableStock, price])
            self.funds -= buyableStock * price

        print(".")
        print("                     Bought %d stocks at %d " %
              (buyableStock, price))

    def sell(self):
        # TODO: need a mechanism for working out how much to sell
        # TODO: need to work out which priced stocks to sell at that time
        price = self.testData[1][0]

        sold = 0

        if len(self.inventory) == 0:
            print(".")
            print("                     Sold %d stocks at %d " %
                  (sold, price))
            return

        for item in self.inventory:
            sold += item[0]
            self.funds += (item[0] * price)

        # right now it sells everything
        self.inventory = []
        print(".")
        print("                     Sold %d stocks at %d " %
              (sold, price))

    def trade(self, first):
        action = analysis.checkPosition(
            self.historicalData[1], self.testData[1][0])
        print("                     Current Price:  £", self.testData[1][0])
        print("                     Moving average: £",
              analysis.calculateSMA(self.historicalData[1]))

        if action > 0 or first:
            self.buy()

        elif action < 0:
            self.sell()

        self.historicalData[1].append(self.testData[1][0])
        self.testData[1].remove(self.testData[1][0])

        return


trader1 = Trader(100000, [], process.historicalData, process.testData)

day = 0

while(len(trader1.testData[1]) != 1):

    print("                     Day:    ", trader1.testData[0][day])
    print("                     Trade:  ", day)
    print("                     Funds:   £", trader1.funds)
    print("                     inventory: ", trader1.inventory)
    trader1.trade(day == 0)
    day += 1
    print(".")
    print(".")
    print(".")


# selling all stocks at the end
trader1.sell()

print(".")
print(".")
print(".")
print(".")
print(".")
print("                     initial funds were  : £ 100000 ")
print("                     final funds are     : £", trader1.funds)
print(".")
print(".")
print(".")
print(".")
print(".")
print(".")
print(".")
print(".")
