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
        buyableStock = price % self.funds

        self.inventory.append([buyableStock, price])
        self.funds -= buyableStock * price
        self.historicalData[1].append(self.testData[1][0])
        self.testData[1].remove(self.testData[1][0])

    def sell(self):
        # TODO: need a mechanism for working out how much to sell
        # TODO: need to work out which priced stocks to sell at that time
        price = self.testData[1][0]

        if len(self.inventory) == 0:
            return

        for item in self.inventory:
            self.funds += (item[0] * price)

        # right now it sells everything
        self.inventory = []
        self.historicalData[1].append(self.testData[1][0])
        self.testData[1].remove(self.testData[1][0])

    def trade(self):
        action = analysis.checkPosition(
            self.historicalData[1], self.testData[1][0])

        if action > 0:
            self.buy()

        if action < 1:
            self.sell()

        return


trader1 = Trader(100000, [], process.historicalData, process.testData)

day = 0

while(len(trader1.testData[1]) != 1):

    print("                     Day: ", trader1.testData[0][day])
    print("")

    trader1.trade()
    day += 1
# selling all stocks at the end
trader1.sell()
print("final funds are: ", trader1.funds)
