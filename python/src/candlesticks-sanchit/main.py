import xlrd
import matplotlib.pyplot as plt
import bots
import utils

# global variable
wb = xlrd.open_workbook("recentAPIdata.xlsx")
sheet = wb.sheet_by_index(0)


def test(bot):
    for i in range(bot.numOfDecisions):
        bot.trade()


def main():
    bot = bots.Bot(10, 0, 100, 0, 1)
    test(bot)
    ys = []
    # for i in range(1, 50000):
    #    ys.append(float(getAsk(i)))

    # plotGraph(ys)


if __name__ == "__main__":
    main()
