import pandas as pd 
import seaborn as sns
import matplotlib.pyplot as plt 
import os

# read multiple csv data
data_lists = []
dir_name="../eval/"

for f in os.listdir(dir_name):
	filename = os.fsdecode(f)
	if filename.endswith(".csv"):
		data_lists.append(pd.read_csv(dir_name+filename))

data = pd.concat(data_lists, axis=0)
		

# data = pd.read_csv("../eval.csv")
plt.rcParams["font.family"] = "Ubuntu Condensed"
fig, ax = plt.subplots(figsize=(6, 3))


pallete = sns.color_palette("husl", 8)
my_pal = {0: pallete[0], 1: pallete[4]}

plt.grid(color="#605B56", linestyle="dotted", linewidth=1, alpha=0.8)
# plot = sns.lineplot(data=data, x="idx", y="val", hue="type", ax=ax)
plot = sns.boxplot(data=data, x="rngRange", y="val", hue="type", ax=ax)


plot.tick_params(left=False, bottom=False, labelsize=14)
ax.set_xlabel("Random number range", fontsize=18)
ax.set_ylabel("Entropy(bit)", fontsize=18)
handles, labels = ax.get_legend_handles_labels()
ax.legend(handles=handles, bbox_to_anchor=(0.5,1.1), loc='center', frameon=False, fontsize=16, ncol=len(labels))

ax.spines['top'].set_visible(False)
ax.spines['right'].set_visible(False)
ax.spines['bottom'].set_visible(False)
ax.spines['left'].set_visible(False)

fig.savefig("./eval.pdf", bbox_inches='tight', pad_inches=0)
