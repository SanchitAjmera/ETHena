import urllib.request
import json
import data


# function for processing data by their curreny


def processCurrency(listOfDictionaries):
    dictionaryOfCurrency = {}
    # looping through all the elements in given list
    for i in range(len(listOfDictionaries)):
        # adding currency to dictionary
        currencyName = listOfDictionaries[i].get('currency')
        dictionaryOfCurrency[currencyName] = listOfDictionaries[i]

    return dictionaryOfCurrency


# function for parsing time string


def parseTime(time):
    # time in the form "2018-04-03T16:00:00Z"
    return time


# function to process times and prices of currencies


def processTimes(dictionaryOfCurrency, currencyID):
    dictionaryOfTimes = dictionaryOfCurrency[currencyID]
    timestamps = dictionaryOfTimes['timestamps']
    prices = dictionaryOfTimes['prices']
    for i in range(len(prices)):
        # converting string to int
        prices[i] = eval(prices[i])
    array = [timestamps, prices]
    return array


# function which prints out processed currency data

def displayProccessedCurrency(currencyID, dicts):
    # printing processed currency data
    print(processCurrency(dicts[0])[currencyID])
    print(processCurrency(dicts[1])[currencyID])
    print(processCurrency(dicts[2])[currencyID])
    print(processCurrency(dicts[4])[currencyID])
    return


# function which returns a list of loaded lists of monthly data


def loadData():
    # loading data into list variables
    dict_1 = json.loads(urllib.request.urlopen(data.URL_APRIL).read())
    dict_2 = json.loads(urllib.request.urlopen(data.URL_MAY).read())
    dict_3 = json.loads(urllib.request.urlopen(data.URL_JUNE).read())
    dict_4 = json.loads(urllib.request.urlopen(data.URL_JULY).read())
    return [dict_1, dict_2, dict_3, dict_4]


dataValues = loadData()
#displayProccessedCurrency('BTC', data)
historicalData = processTimes(processCurrency(dataValues[1]), 'BTC')
testData = processTimes(processCurrency(dataValues[2]), 'BTC')
