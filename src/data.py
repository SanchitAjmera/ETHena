import urllib.request
import json

# monthly data from 2019
URL_APRIL = "https://api.nomics.com/v1/currencies/sparkline?key=c7154c3c700596fd7c4d234f71d3feb8&ids=BTC,ETH,XRP&start=2019-04-14T00%3A00%3A00Z&end=2019-05-14T00%3A00%3A00Z"
URL_MAY = "https://api.nomics.com/v1/currencies/sparkline?key=c7154c3c700596fd7c4d234f71d3feb8&ids=BTC,ETH,XRP&start=2019-05-15T00%3A00%3A00Z&end=2019-06-14T00%3A00%3A00Z"
URL_JUNE = "https://api.nomics.com/v1/currencies/sparkline?key=c7154c3c700596fd7c4d234f71d3feb8&ids=BTC,ETH,XRP&start=2019-06-15T00%3A00%3A00Z&end=2019-07-14T00%3A00%3A00Z"


# example template of information inside the API
"""
{
    "currency": "BTC",
    "timestamps": [
        "2018-04-03T16:00:00Z"
    ],
    "prices": [
        "7436.82313"
    ]
}
"""

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
    dictionaryOfTimes = dictionaryOfCurrency['currencyID']
    timestamps = dictionaryOfTimes['timestamps']
    prices = dictionaryOfTimes['prices']
    array = [timestamps, prices]
    return array


# function which prints out processed currency data

def displayProccessedCurrency(currencyID, dicts):
    # printing processed currency data
    print(processCurrency(dicts[0])[currencyID])
    print(processCurrency(dicts[1])[currencyID])
    print(processCurrency(dicts[2])[currencyID])
    return


# function which returns a list of loaded lists of monthly data


def loadData():
    # loading data into list variables
    dict_1 = json.loads(urllib.request.urlopen(URL_APRIL).read())
    dict_2 = json.loads(urllib.request.urlopen(URL_MAY).read())
    dict_3 = json.loads(urllib.request.urlopen(URL_JUNE).read())
    return [dict_1, dict_2, dict_3]


data = loadData()
displayProccessedCurrency('BTC', data)

# print(urllib.request.urlopen(url).read())
