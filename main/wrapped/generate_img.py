import sqlite3
import matplotlib.pyplot as plt
import matplotlib.patches as patches
import numpy as np
from matplotlib import rcParams
from matplotlib.colors import LinearSegmentedColormap
import random

# ---------- CONFIG ----------
DB_PATH = "data/poop_tracker.db"   # <-- Update this to your DB path
OUTPUT_PATH = "poop_wrapped_2025.png"

YEAR = 2025  # Adjust as needed

# Emoji-friendly font
rcParams['font.family'] = 'Segoe UI Emoji'  # Windows; Mac: 'Apple Color Emoji'

# ---------- FUNCTIONS ----------
def fetch_yearly_stats(db_path, year):
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    cursor.execute("""
    SELECT username, COUNT(*) as poop_count
    FROM poop_tracker
    WHERE strftime('%Y', timestamp) = ?
    GROUP BY user_id
    ORDER BY poop_count DESC;
    """, (str(year),))
    data = cursor.fetchall()
    conn.close()
    return data

def add_confetti(ax, top_index, max_value, n=50):
    """Add subtle confetti behind top pooper"""
    for _ in range(n):
        y = top_index + 0.3 + random.uniform(-0.15, 0.15)
        x = random.uniform(0, 0.7 * max_value)
        size = random.uniform(10, 60)
        color = random.choice(['#FFD27F', '#FFB6B9', '#A0E7E5', '#B5EAD7', '#FFDAC1'])
        ax.scatter(x, y, s=size, color=color, alpha=0.5, zorder=1)

def generate_wrapped_chart(data, output_path, title="💩 Poop Wrapped"):
    if not data:
        print("No data to plot!")
        return

    users = [row[0] for row in data]
    poops = [row[1] for row in data]
    top_index = np.argmax(poops)
    max_value = max(poops)

    # ---------- Colors ----------
    bar_colors = ["#FFD27F" if i == top_index else "#A0E7E5" for i in range(len(users))]

    # ---------- Plot ----------
    fig, ax = plt.subplots(figsize=(12, 7))
    fig.patch.set_facecolor('#FFF0F5')  # pastel pink background
    ax.set_facecolor('#FFF0F5')

    # Rounded bars manually
    for i, (user, count) in enumerate(zip(users, poops)):
        bar = patches.FancyBboxPatch(
            (0, i - 0.3),  # x, y
            count,  # width
            0.6,  # height
            boxstyle="round,pad=0.02",
            linewidth=0,
            facecolor=bar_colors[i],
            zorder=2
        )
        ax.add_patch(bar)
        ax.text(
            count + max_value*0.01, i, f"{count} 💩",
            va='center', ha='left', fontsize=14, fontweight='bold', zorder=3
        )

    # ---------- Confetti ----------
    add_confetti(ax, top_index, max_value)

    # ---------- Aesthetics ----------
    ax.set_yticks(np.arange(len(users)))
    ax.set_yticklabels(users, fontsize=12)
    ax.set_xlabel("Total Poops", fontsize=16, color='#555')
    ax.set_xlim(0, max_value*1.3)
    ax.invert_yaxis()  # largest on top
    ax.set_title(f"{title} {YEAR} 💩", fontsize=26, fontweight='bold', pad=30, color='#333')

    # Remove spines
    for spine in ax.spines.values():
        spine.set_visible(False)

    # Subtle grid
    ax.xaxis.grid(True, linestyle='--', alpha=0.2, zorder=0)

    plt.tight_layout()
    plt.savefig(output_path, dpi=150, facecolor=fig.get_facecolor())
    plt.close()
    print(f"Modern Spotify-Wrapped style chart saved to {output_path}")

# ---------- MAIN ----------
if __name__ == "__main__":
    yearly_data = fetch_yearly_stats(DB_PATH, YEAR)
    generate_wrapped_chart(yearly_data, OUTPUT_PATH)