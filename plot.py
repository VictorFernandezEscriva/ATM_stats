import sys
import json
import matplotlib.pyplot as plt
import numpy as np
from scipy.stats import norm

from pprint import pprint

def plot_compliance(ax, compliant, total, title):
    values = [
        compliant,
        total - compliant,
    ]

    wedges, _, _ = ax.pie(values, autopct=lambda pct: f"{int(round(pct/100*total))} ({pct:1.1f}%)")
    ax.legend(wedges, ["Compliant flights", "Not compliant flights"])
    ax.set_title(title)

def plot_by_airline(radar: dict, wake: dict, loa: dict, results):
    fig, ax = plt.subplots()

    airlines = list(set(radar.keys()) | set(wake.keys()) | set(loa.keys()))

    not_compliances = {
        "LoA": np.array([loa[x] if x in loa else 0 for x in airlines]),
        "Radar": np.array([radar[x] if x in radar else 0 for x in airlines]),
        "Wake": np.array([wake[x] if x in wake else 0 for x in airlines]),
    }
    width = 0.5

    bottom = np.zeros(len(airlines))

    for type_compliance, not_compliances in not_compliances.items():
        ax.bar(airlines, not_compliances, width, label=type_compliance, bottom=bottom)
        bottom += not_compliances

    ax.set_title("Not compliances by airline")
    ax.legend(loc="upper right")

    fig, ax = plt.subplots()

    values = []
    imp_airlines = []
    other = 0
    for i in range(len(airlines)):
        if bottom[i] > 1:
            values.append(bottom[i])
            imp_airlines.append(airlines[i])
        else:
            other += bottom[i]

    values.append(other)
    imp_airlines.append("Other")
    total = sum(values)

    wedges, _, _ = ax.pie(values, autopct=lambda pct: f"{int(round(pct/100*total))} ({pct:1.1f}%)")
    ax.legend(wedges, imp_airlines)

    fig, ax = plt.subplots()

    values = []
    imp_airlines = []
    other = 0
    airlines = {}
    for result in results:
        airline = result["Second"]["Company"]
        if airline not in airlines:
            airlines[airline] = 0
        airlines[airline] += 1
    
    other = 0
    for v in airlines.values():
        if v <= 2:
            other += v
    airlines = {k: v for k, v in airlines.items() if v > 2}
    airlines["Other"] = other
    print(airlines)
    total = sum(list(airlines.values()))

    wedges, _ = ax.pie(list(airlines.values()))
    ax.legend(wedges, list(airlines.keys()))
    ax.set_title("Total flights by airline")

def plot_by_wake(wake: dict):
    fig, ax = plt.subplots()

    wake_categories = list(wake.keys())
    values = np.array(list(wake.values()))
    total = sum(values)
    wedges, _, _ = ax.pie(values, autopct=lambda pct: f"{int(round(pct/100*total))} ({pct:1.1f}%)")
    ax.legend(wedges, wake_categories)

def plot_by_class(loa: dict):
    print(loa)
    fig, ax = plt.subplots()

    loa_classes = [k for k, v in loa.items() if v != 0]
    values = [v for k, v in loa.items() if k in loa_classes]
    total = sum(values)
    wedges, _, _ = ax.pie(values, autopct=lambda pct: f"{int(round(pct/100*total))} ({pct:1.1f}%)")
    ax.legend(wedges, list(loa.keys()))

def plot_by_sid(sid: dict):
    print(sid)
    fig, ax = plt.subplots()

    loa_classes = [k for k, v in sid.items() if v != 0]
    values = [v for k, v in sid.items() if k in loa_classes]
    total = sum(values)
    wedges, _, _ = ax.pie(values, autopct=lambda pct: f"{int(round(pct/100*total))} ({pct:1.1f}%)")
    ax.legend(wedges, list(sid.keys()))

def minimum_radar_compliance_percentage(results):
    fig, ax = plt.subplots()

    # Compliance
    compliant = 0
    by_airline = {}
    for result in results:
        if result["MinDistance"] > 3:
            compliant += 1
        else:
            if not result["Second"]["Company"] in by_airline:
                by_airline[result["Second"]["Company"]] = 0
            by_airline[result["Second"]["Company"]] += 1
        
    plot_compliance(ax, compliant, len(results), 'Minimum radar compliance percentage')
    return by_airline

def wake_compliance_percentage(results) -> (dict, dict):
    fig, ax = plt.subplots()

    # Compliance
    compliant = 0
    by_airline = {}
    by_wake = {}
    for result in results:
        sep_minima = 0
        key = (result["First"]["Wake"], result["Second"]["Wake"])
        match key:
            case ("Super heavy", "Heavy"):
                sep_minima = 6
            case ("Super heavy", "Medium"):
                sep_minima = 7
            case ("Super heavy", "Light"):
                sep_minima = 8
            case ("Heavy", "Heavy"):
                sep_minima = 4
            case ("Heavy", "Medium"):
                sep_minima = 5
            case ("Heavy", "Light"):
                sep_minima = 6
            case ("Medium", "Light"):
                sep_minima = 5
        
        key = key[0] + "-" + key[1]
        if not key in by_wake:
            by_wake[key] = 0
        if result["MinDistance"] > sep_minima:
            compliant += 1
        else:
            if not result["Second"]["Company"] in by_airline:
                by_airline[result["Second"]["Company"]] = 0
            by_airline[result["Second"]["Company"]] += 1
            by_wake[key] += 1
        
    plot_compliance(ax, compliant, len(results), 'Wake compliance percentage')
    return by_airline, by_wake

def loa_compliance_percentage(results) -> (dict, dict, dict):
    fig, ax = plt.subplots()

    # Compliance
    compliant = 0
    by_airline = {}
    by_class = {}
    by_sid = {}
    for result in results:
        sep_minima = 0
        same_sid = result["First"]["SidGroup"] == result["Second"]["SidGroup"]
        key = (result["First"]["Class"], result["Second"]["Class"], same_sid)
        match key:
            case ("HP", "HP", True):
                sep_minima = 5
            case ("HP", "HP", False):
                sep_minima = 3
            case ("HP", "R", True):
                sep_minima = 5
            case ("HP", "R", False):
                sep_minima = 3
            case ("HP", "LP", True):
                sep_minima = 5
            case ("HP", "LP", False):
                sep_minima = 3
            case ("HP", "NR+", _):
                sep_minima = 3
            case ("HP", "NR-", _):
                sep_minima = 3
            case ("HP", "NR", _):
                sep_minima = 3
            case ("R", "HP", True):
                sep_minima = 7
            case ("R", "HP", False):
                sep_minima = 5
            case ("R", "R", True):
                sep_minima = 5
            case ("R", "R", False):
                sep_minima = 3
            case ("R", "LP", True):
                sep_minima = 5
            case ("R", "LP", False):
                sep_minima = 3
            case ("R", "NR+", _):
                sep_minima = 3
            case ("R", "NR-", _):
                sep_minima = 3
            case ("R", "NR", _):
                sep_minima = 3
            case ("LP", "HP", True):
                sep_minima = 8
            case ("LP", "HP", False):
                sep_minima = 6
            case ("LP", "R", True):
                sep_minima = 6
            case ("LP", "R", False):
                sep_minima = 4
            case ("LP", "LP", True):
                sep_minima = 5
            case ("LP", "LP", False):
                sep_minima = 3
            case ("LP", "NR+", _):
                sep_minima = 3
            case ("LP", "NR-", _):
                sep_minima = 3
            case ("LP", "NR", _):
                sep_minima = 3
            case ("NR+", "HP", True):
                sep_minima = 11
            case ("NR+", "HP", False):
                sep_minima = 8
            case ("NR+", "R", True):
                sep_minima = 9
            case ("NR+", "R", False):
                sep_minima = 6
            case ("NR+", "LP", True):
                sep_minima = 9
            case ("NR+", "LP", False):
                sep_minima = 6
            case ("NR+", "NR+", True):
                sep_minima = 5
            case ("NR+", "NR+", False):
                sep_minima = 3
            case ("NR+", "NR-", _):
                sep_minima = 3
            case ("NR+", "NR", _):
                sep_minima = 3
            case ("NR-", "HP", _):
                sep_minima = 9
            case ("NR-", "R", _):
                sep_minima = 9
            case ("NR-", "LP", _):
                sep_minima = 9
            case ("NR-", "NR+", True):
                sep_minima = 9
            case ("NR-", "NR+", False):
                sep_minima = 6
            case ("NR-", "NR-", True):
                sep_minima = 5
            case ("NR-", "NR-", False):
                sep_minima = 3
            case ("NR-", "NR", _):
                sep_minima = 3
            case ("NR", "HP", _):
                sep_minima = 9
            case ("NR", "R", _):
                sep_minima = 9
            case ("NR", "LP", _):
                sep_minima = 9
            case ("NR", "NR+", _):
                sep_minima = 9
            case ("NR", "NR-", _):
                sep_minima = 9
            case ("NR", "NR", True):
                sep_minima = 5
            case ("NR", "NR", False):
                sep_minima = 3

        class_key = key[0] + " | " + key[1]
        if not class_key in by_class:
            by_class[class_key] = 0
        sid_key =  result["First"]["SidGroup"] + " | " + result["Second"]["SidGroup"]
        if not sid_key in by_sid:
            by_sid[sid_key] = 0
        if result["MinDistance"] > sep_minima:
            compliant += 1
        else:
            if not result["Second"]["Company"] in by_airline:
                by_airline[result["Second"]["Company"]] = 0
            by_airline[result["Second"]["Company"]] += 1
            by_class[class_key] += 1
            by_sid[sid_key] += 1

    plot_compliance(ax, compliant, len(results), 'LoA compliance percentage')
    return by_airline, by_class, by_sid

def plot_normal_dist(data, ax, xlabel, ylabel):
    num_bins = int(max(data) - min(data)) + 1
    ax.hist(data, bins=num_bins, density=True, alpha=0.5, color='b', edgecolor='black')

    mu, std = norm.fit(data)
    perc95 = np.percentile(data, 95)

    xmin, xmax = ax.get_xlim()
    x = np.linspace(xmin, xmax, 100)
    p = norm.pdf(x, mu, std)
    ax.plot(x, p, 'k', linewidth=2)

    ax.set_title("Fit results: mu = %.2f,  std = %.2f, percentile 95 = %.2f" % (mu, std, perc95))
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
            #pprint(result)
            pass

    by_airline_radar = minimum_radar_compliance_percentage(results["Results"])
    by_airline_wake, by_wake = wake_compliance_percentage(results["Results"])
    by_airline_loa, by_class, by_sid = loa_compliance_percentage(results["Results"])

    minimum_radar_normal_distributions(minDistances)
    plot_by_airline(by_airline_radar, by_airline_wake, by_airline_loa, results["Results"])
    plot_by_wake(by_wake)
    plot_by_class(by_class)
    plot_by_sid(by_sid)

    plt.show()

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python plot.py <json_file>")
        sys.exit(1)

    json_file = sys.argv[1]
    plot_statistics(json_file)