
import matplotlib.pyplot as plt
import bots
import utils


def test(bot):
    for i in range(bot.numOfDecisions):
        bot.trade()


def main():
    bot = bots.Bot(10, 0, 100, 0, 11)
    test(bot)
    ys = []
    print(bot.sells)
    print(bot.buys)
    for i in range(1, 50000):
        ys.append(utils.getAsk(i))

    utils.plotGraph(ys, bot)


if __name__ == "__main__":
    main()
