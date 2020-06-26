
AVERAGE_LENGTH = 5
# function to calculate the simple moving average over the last n days


def calculateSMA(historicalData):
    length = len(historicalData)
    return sum(historicalData[0][length - AVERAGE_LENGTH - 1:]) / AVERAGE_LENGTH


def calculateOffset():
    return

    # function to calculate the exponential moving average
    # places higher weightings on recent data


def calculateEMA(historicalData):
    return


def checkPosition(historicalData, testPrice):
    if calculateSMA(historicalData) >:
        pass
