import process
import matplotlib.pyplot as plt

timearray = []

for i in range(len(process.dataTimesPrices[1])):
    timearray.append(i)

plt.plot(timearray, process.dataTimesPrices[1])

plt.show()
