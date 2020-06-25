import urllib.request
import json
url = "https://api.nomics.com/v1/currencies/ticker?key=c7154c3c700596fd7c4d234f71d3feb8"


# example template of information inside the API
"""
{
    "currency": "BTC",
    "id": "BTC",
    "price": "8451.36516421",
    "price_date": "2019-06-14T00:00:00Z",
    "price_timestamp": "2019-06-14T12:35:00Z",
    "symbol": "BTC",
    "circulating_supply": "17758462",
    "max_supply": "21000000",
    "name": "Bitcoin",
    "logo_url": "https://s3.us-east-2.amazonaws.com/nomics-api/static/images/currencies/btc.svg",
    "market_cap": "150083247116.70",
    "transparent_market_cap": "150003247116.70",
    "rank": "1",
    "high": "19404.81116899",
    "high_timestamp": "2017-12-16",
    "1d": {
        "price_change": "269.75208019",
        "price_change_pct": "0.03297053",
        "volume": "1110989572.04",
        "volume_change": "-24130098.49",
        "volume_change_pct": "-0.02",
        "market_cap_change": "4805518049.63",
        "market_cap_change_pct": "0.03",
        "transparent_market_cap_change": "4800518049.00",
        "transparent_market_cap_change_pct": "0.0430",
        "volume_transparency": [
            {
                "grade": "A",
                "volume": "2144455081.37",
                "volume_change": "-235524941.08",
                "volume_change_pct": "-0.10"
            },
            {
                "grade": "B",
                "volume": "15856762.85",
                "volume_change": "-6854329.88",
                "volume_change_pct": "-0.30"
            }
        ]
    }
}
"""

# loading data into variable
dict_1 = json.loads(urllib.request.urlopen(url).read())

# printing out curreny name
for i in range(len(dict_1)):
    print(dict_1[i].get('currency'))
# print(urllib.request.urlopen(url).read())
