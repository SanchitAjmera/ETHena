import process
import analysis


class Trader:

    # initialising trader state

    def __init__(self, initialFund, initialInventory, givenData, testData):
        self.funds = initialFund
        self.inventory = initialInventory
        self.historicalData = givenData
        self.testData = testData

    # function to buy stocks

    def buy(self):
        # TODO: need a mechanism for working out how much to buy
        price = self.testData[1][0]
        buyableStock = self.funds / price
        # check to make sure trader has enough funds to buy stock
        if buyableStock != 0:
            self.inventory.append([buyableStock, price])
            self.funds -= buyableStock * price

        print(".")
        print("                     Bought %d stocks at %d " %
              (buyableStock, price))

    # function to sell stocks

    def sell(self):
        # TODO: need a mechanism for working out how much to sell
        # TODO: need to work out which priced stocks to sell at that time
        price = self.testData[1][0]
        sold = 0
        # check for if inventory is empty so cannot sell anything
        if len(self.inventory) == 0:
            print(".")
            print("                     Sold %d stocks at %d " %
                  (sold, price))
            return
        # selling all items in the inventory
        for item in self.inventory:
            sold += item[0]
            self.funds += (item[0] * price)
        # reseting inventory
        self.inventory = []
        print(".")
        print("                     Sold %d stocks at %d " %
              (sold, price))

    # function to carry out trade actions

    def trade(self, first):
        # using function in analysis to check what action to conduct
        action = analysis.checkPosition(
            self.historicalData[1], self.testData[1][0])
        print("                     Current Price:  £", self.testData[1][0])
        print("                     Moving average: £",
              analysis.calculateSMA(self.historicalData[1]))

        if action > 0 or first:
            self.buy()

        elif action < 0:
            self.sell()
        # removing current day stats and appending them to historical data
        # this refershes the data for a new moving average to be calculated
        self.historicalData[1].append(self.testData[1][0])
        self.testData[1].remove(self.testData[1][0])

        return


# initialising trader
trader1 = Trader(100000, [], process.historicalData, process.testData)

day = 0

# trader trades until last day
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


# selling all stocks at the end on the last day
trader1.sell()


# printing out final funds
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
