import process
import matplotlib.pyplot as plt


timearray = []

for i in range(len(process.historicalData[1])):
    timearray.append(i)

plt.plot(timearray, process.historicalData[1])

plt.show()
plt.close()
