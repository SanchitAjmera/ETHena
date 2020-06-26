
# need a function which cleverly calculates the offset
OFFSET = 2000
AVERAGE_LENGTH = 5
# function to calculate the simple moving average over the last n days


def calculateSMA(historicalData):
    length = len(historicalData)

    return sum(historicalData[length - AVERAGE_LENGTH:]) / AVERAGE_LENGTH


def calculateOffset():
    return

    # function to calculate the exponential moving average
    # places higher weightings on recent data


def calculateEMA(historicalData):
    return


def checkPosition(historicalData, testPrice):
    average = calculateSMA(historicalData)

    if average + OFFSET < testPrice:
        # sell
        return -1
    elif average - OFFSET > testPrice:
        # buy
        return 1
    else:
        # speculate
        return 0
