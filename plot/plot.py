import pandas as pd 
import seaborn as sns
import matplotlib.pyplot as plt 

data = pd.read_csv("../eval.csv")
plt.rcParams["font.family"] = "Ubuntu Condensed"
fig, ax = plt.subplots(figsize=(4.5, 3))

'''
pallete = sns.color_palette("husl", 8)
my_pal = {"0": pallete[0], "1": pallete[4]}
'''
plt.grid(color="#605B56", linestyle="dotted", linewidth=1, alpha=0.8)
plot = sns.lineplot(data=data, x="idx", y="val", hue="type")


plot.tick_params(left=False, bottom=False, labelsize=14)
ax.set_xlabel("Time (per entropy trial)", fontsize=18)
ax.set_ylabel("Evaluated entropy", fontsize=18)
handles, labels = ax.get_legend_handles_labels()
ax.legend(handles=handles, labels=["LRNG3", "DIFF"], bbox_to_anchor=(0.5,1.1), loc='center', frameon=False, fontsize=16, ncol=2)

ax.spines['top'].set_visible(False)
ax.spines['right'].set_visible(False)
ax.spines['bottom'].set_visible(False)
ax.spines['left'].set_visible(False)

fig.savefig("./eval.pdf", bbox_inches='tight', pad_inches=0)
