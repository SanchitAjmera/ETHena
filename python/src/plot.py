import process
import matplotlib.pyplot as plt


timearray = []

for i in range(len(process.testData[1])):
    timearray.append(i)

plt.plot(timearray, process.testData[1])

plt.show()
plt.close()
