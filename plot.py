import sys
import json
import matplotlib.pyplot as plt
import numpy as np
from scipy.stats import norm

from pprint import pprint

def minimum_radar_compliance_percentage(distances):
    plt.figure()

    # Compliance
    compliant = len([d for d in distances if d > 3])
    values = [
        compliant,
        len(distances) - compliant,
    ]

    wedges, _, _ = plt.pie(values, autopct="%1.1f%%")
    plt.legend(wedges, ["Compliant flights", "Not compliant flights"],
    )
    plt.title('Minimum radar compliance percentage')

def plot_normal_dist(data, ax, xlabel, ylabel):
    num_bins = int(max(data) - min(data)) + 1
    ax.hist(data, bins=num_bins, density=True, alpha=0.5, color='b', edgecolor='black')

    mu, std = norm.fit(data)

    xmin, xmax = ax.get_xlim()
    x = np.linspace(xmin, xmax, 100)
    p = norm.pdf(x, mu, std)
    ax.plot(x, p, 'k', linewidth=2)

    ax.set_title("Fit results: mu = %.2f,  std = %.2f" % (mu, std))
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)

def minimum_radar_normal_distributions(minDistances):
    # Normal distribution
    #_, (ax1, ax2) = plt.subplots(nrows=1, ncols=2)

    fig, ax = plt.subplots()
    plot_normal_dist(minDistances, ax, "Minimum distance between consecutive flights [NM]", "Frequency")

def plot_statistics(json_file):
    with open(json_file, 'r') as file:
        results = json.load(file)
    

    minDistances = [x["MinDistance"] for x in results["Results"]]

    for result in results["Results"]:
        if result["MinDistance"] < 3:
            pprint(result)

    minimum_radar_compliance_percentage(minDistances)
    minimum_radar_normal_distributions(minDistances)

    plt.show()

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python plot.py <json_file>")
        sys.exit(1)

    json_file = sys.argv[1]
    plot_statistics(json_file)