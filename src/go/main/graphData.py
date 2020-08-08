import pandas as pd
import matplotlib.pyplot as plt
from datetime import date

lower_limit = 25
date = str(date.today())
df = pd.read_excel(date + ".xlsx")

lower_limit_list = []
for i in df['Sr No.']:
    lower_limit_list.append(lower_limit)

plt.figure()
plt.subplot(3, 1, (1, 2))
plt.plot(df['Sr No.'], df['Ready To Buy Price'], color='r', label='Sold')
plt.plot(df['Sr No.'], df['Ready To Sell Price'], color='g', label='Bought')
plt.grid(b=True, which='both', axis='both')
plt.legend()
plt.ylabel('Price')
plt.title('Summary of ' + date)

plt.subplot(3, 1, 3)
plt.plot(df['Sr No.'], df['RSI'], color='blue', label='RSI')
plt.plot(df['Sr No.'], lower_limit_list, color='yellow')
plt.grid(b=True, which='both', axis='both')
plt.legend()
plt.xlabel('minute')
plt.ylabel('RSI')
plt.savefig('graph.png')
