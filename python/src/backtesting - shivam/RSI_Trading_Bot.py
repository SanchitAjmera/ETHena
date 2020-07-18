from RSI_Calculator import *
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np

def RSI_calculator(RSI_period,df, pair):


    close_list = list(df[pair+'_ask'])

    if str(close_list[0]) == 'nan' and str(close_list[1]) != 'nan':
        close_list[0] = close_list[1]

    for i in range(1, len(close_list)):
        if str(close_list[i]) == 'nan':
            close_list[i] = close_list[i - 1]

    RSI_list = [50]
    for i in range(1,RSI_period):
        positive_list=[]
        negative_list=[]
        for j in range(1,i+1):
            percentage_difference = (close_list[j]-close_list[j-1])/close_list[j]
            if percentage_difference>=0:
                positive_list.append(percentage_difference)
            if percentage_difference<0:
                negative_list.append(percentage_difference)
        average_gain = np.mean(positive_list)
        absolute_average_loss = abs(np.mean(negative_list))
        if str(average_gain) =='nan':
            rsi = 50
        elif str(absolute_average_loss) =='nan':
            rsi = 50
        else:
            rsi = 100-(100/(1+(average_gain/absolute_average_loss)))
        RSI_list.append(rsi)

    for i in range(RSI_period,len(close_list)):
        positive_list=[]
        negative_list=[]
        for j in range(0,RSI_period):
            percentage_difference = (close_list[i-RSI_period+j]-close_list[i-RSI_period+j-1])/close_list[i-RSI_period+j]
            if percentage_difference>=0:
                positive_list.append(percentage_difference)
            if percentage_difference<0:
                negative_list.append(percentage_difference)
        average_gain = np.mean(positive_list)
        absolute_average_loss = abs(np.mean(negative_list))
        if str(average_gain) =='nan':
            rsi = 50
        elif str(absolute_average_loss) =='nan':
            absolute_average_loss = 1
            rsi = 100 - (100 / (1 + (average_gain / absolute_average_loss)))
        else:
            rsi = 100-(100/(1+(average_gain/absolute_average_loss)))
        RSI_list.append(rsi)

    df['RSI' + str(RSI_period)] = RSI_list

RSI_period = 14

pair_list = ['XBTGBP', 'ETHXBT', 'XRPXBT', 'XRPZAR', 'BCHXBT', 'LTCXBT', 'XBTZAR']

# initialise the lists

XBTGBP_bid = []
XBTGBP_ask = []

ETHXBT_bid = []
ETHXBT_ask = []

XRPXBT_bid = []
XRPXBT_ask = []

XRPZAR_bid = []
XRPZAR_ask = []

BCHXBT_bid = []
BCHXBT_ask = []

LTCXBT_bid = []
LTCXBT_ask = []

XBTZAR_bid = []
XBTZAR_ask = []

for i in range(0, 15):
    df = pd.read_excel("tickerData" + str(i) + ".xlsx")
    for i in range(0, len(df["XBTGBP bid"])):
        XBTGBP_bid.append(df["XBTGBP bid"][i])
        XBTGBP_ask.append(df["XBTGBP ask"][i])

        ETHXBT_bid.append(df["ETHXBT bid"][i])
        ETHXBT_ask.append(df["ETHXBT ask"][i])

        XRPXBT_bid.append(df["XRPXBT bid"][i])
        XRPXBT_ask.append(df["XRPXBT ask"][i])

        XRPZAR_bid.append(df["XRPZAR bid"][i])
        XRPZAR_ask.append(df["XRPZAR ask"][i])

        BCHXBT_bid.append(df["BCHXBT bid"][i])
        BCHXBT_ask.append(df["BCHXBT ask"][i])

        LTCXBT_bid.append(df["LTCXBT bid"][i])
        LTCXBT_ask.append(df["LTCXBT ask"][i])

        XBTZAR_bid.append(df["XBTZAR bid"][i])
        XBTZAR_ask.append(df["XBTZAR ask"][i])

data = {'XBTGBP_bid': XBTGBP_bid, 'XBTGBP_ask': XBTGBP_ask, 'ETHXBT_bid': ETHXBT_bid, 'ETHXBT_ask': ETHXBT_ask,
        'XRPXBT_bid': XRPXBT_bid,
        'XRPXBT_ask': XRPXBT_ask, 'XRPZAR_bid': XRPZAR_bid, 'XRPZAR_ask': XRPZAR_ask, 'BCHXBT_bid': BCHXBT_bid,
        'BCHXBT_ask': BCHXBT_ask,
        'LTCXBT_bid': LTCXBT_bid, 'LTCXBT_ask': LTCXBT_ask, 'XBTZAR_bid': XBTZAR_bid, 'XBTZAR_ask': XBTZAR_ask}

df = pd.DataFrame(data)

xlist = np.arange(0, len(XBTGBP_bid))
list_lists = [XBTGBP_bid, XBTGBP_ask, ETHXBT_bid, ETHXBT_ask, XRPXBT_bid, XRPXBT_ask, XRPZAR_bid, XRPZAR_ask,
              BCHXBT_bid, BCHXBT_ask, LTCXBT_bid, LTCXBT_ask, XBTZAR_bid, XBTZAR_ask]

pair = 'XRPXBT'

RSI_calculator(14, df, pair)
buy_days = []
sell_days = []

close_list = list(df[pair+'_ask'])

if str(close_list[0]) == 'nan' and str(close_list[1]) != 'nan':
    close_list[0] = close_list[1]

for i in range(1, len(close_list)):
    if str(close_list[i]) == 'nan':
        close_list[i] = close_list[i - 1]

open_list = list(df[pair+'_bid'])

if str(open_list[0]) == 'nan' and str(open_list[1]) != 'nan':
    open_list[0] = open_list[1]

for i in range(1, len(open_list)):
    if str(open_list[i]) == 'nan':
        open_list[i] = open_list[i - 1]

RSI_list = list(df['RSI'+str(RSI_period)])

state = 'wait'
for i in range(0,len(RSI_list)):
    if state == 'wait':
        if RSI_list[i]<20:
            state = 'hold'
            buy_days.append(i)
    if state == 'hold':
        if open_list[i]>close_list[buy_days[-1]]:
            state = 'wait'
            sell_days.append(i)

a = 1
b = 1
change_list = []
buy_prices=[]
sell_prices=[]

if len(buy_days)-1 == len(sell_days):
    for i in range(0,len(buy_days)-1):
        buy_day = buy_days[i]
        sell_day = sell_days[i]
        buy_price = open_list[buy_day]
        sell_price = open_list[sell_day]
        buy_prices.append(buy_price)
        sell_prices.append(sell_price)
        absolute_change = sell_price-buy_price
        percentage_change = absolute_change/buy_price
        change_list.append(percentage_change)
        a = a*(1+percentage_change)
else:
    for i in range(0,len(buy_days)):
        buy_day = buy_days[i]
        sell_day = sell_days[i]
        buy_price = 0.999999*open_list[buy_day]
        sell_price = 1.000001*open_list[sell_day]
        buy_prices.append(buy_price)
        sell_prices.append(sell_price)
        absolute_change = sell_price-buy_price
        percentage_change = absolute_change/buy_price
        change_list.append(percentage_change)
        a = a*(1+percentage_change)

plt.figure()
plt.subplot(211)
plt.plot(df['XRPXBT_ask'], label = 'Buy Price')
plt.ylabel('XRPXBT Price')
plt.grid(True)

plt.subplot(212)
plt.plot(df['RSI'+str(RSI_period)], label = 'RSI'+ str(RSI_period))
plt.ylabel('RSI')
plt.grid(True)

plt.show()

print(change_list)
print(a)
print((a-b)/b)
print(buy_prices)
print(sell_prices)
