#!/usr/bin/env python3
import json
import sys
import os
from pathlib import Path

# Add parent directory to path for imports
sys.path.insert(0, str(Path(__file__).parent.parent.parent))

import matplotlib.pyplot as plt
import matplotlib.patches as patches
from matplotlib import rcParams
import numpy as np

# Configuration
SPOTIFY_GREEN = '#1DB954'
SPOTIFY_DARK = '#191414'
SPOTIFY_WHITE = '#FFFFFF'
SLIDE_WIDTH = 1200
SLIDE_HEIGHT = 800
DPI = 150

# Font configuration
rcParams['font.family'] = 'sans-serif'
rcParams['font.sans-serif'] = ['Arial', 'Helvetica', 'DejaVu Sans']

def generate_title_slide(year, output_dir):
    """Generate title slide"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100), facecolor=SPOTIFY_DARK)
    ax.set_facecolor(SPOTIFY_DARK)
    ax.axis('off')
    
    # Title
    ax.text(0.5, 0.6, f'Your Poop Wrapped {year}', 
            ha='center', va='center', fontsize=48, fontweight='bold', 
            color=SPOTIFY_WHITE)
    ax.text(0.5, 0.4, '💩', ha='center', va='center', fontsize=72)
    
    output_path = os.path.join(output_dir, 'slide_01_title.png')
    plt.savefig(output_path, dpi=DPI, facecolor=SPOTIFY_DARK, bbox_inches='tight')
    plt.close()
    return output_path

def generate_total_slide(stats, output_dir):
    """Generate total count slide"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100), facecolor=SPOTIFY_DARK)
    ax.set_facecolor(SPOTIFY_DARK)
    ax.axis('off')
    
    # Total count
    ax.text(0.5, 0.7, f'{stats["TotalPoops"]}', 
            ha='center', va='center', fontsize=72, fontweight='bold', 
            color=SPOTIFY_GREEN)
    ax.text(0.5, 0.5, 'Total dumps', 
            ha='center', va='center', fontsize=32, color=SPOTIFY_WHITE)
    
    # Comparison if available
    if stats.get("GroupRank"):
        rank = stats["GroupRank"]
        percentage = rank["Percentage"]
        ax.text(0.5, 0.3, f'That\'s more than {percentage:.1f}% of the group', 
                ha='center', va='center', fontsize=24, color=SPOTIFY_WHITE, style='italic')
    
    output_path = os.path.join(output_dir, 'slide_02_total.png')
    plt.savefig(output_path, dpi=DPI, facecolor=SPOTIFY_DARK, bbox_inches='tight')
    plt.close()
    return output_path

def generate_streak_slide(stats, output_dir):
    """Generate streak slide"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100), facecolor=SPOTIFY_DARK)
    ax.set_facecolor(SPOTIFY_DARK)
    ax.axis('off')
    
    ax.text(0.5, 0.7, f'{stats["MaxStreak"]}', 
            ha='center', va='center', fontsize=72, fontweight='bold', 
            color=SPOTIFY_GREEN)
    ax.text(0.5, 0.5, 'Longest streak', 
            ha='center', va='center', fontsize=32, color=SPOTIFY_WHITE)
    ax.text(0.5, 0.3, 'Your bowels respect routine', 
            ha='center', va='center', fontsize=24, color=SPOTIFY_WHITE, style='italic')
    
    output_path = os.path.join(output_dir, 'slide_03_streak.png')
    plt.savefig(output_path, dpi=DPI, facecolor=SPOTIFY_DARK, bbox_inches='tight')
    plt.close()
    return output_path

def generate_extreme_day_slide(stats, output_dir):
    """Generate extreme day slide"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100), facecolor=SPOTIFY_DARK)
    ax.set_facecolor(SPOTIFY_DARK)
    ax.axis('off')
    
    ax.text(0.5, 0.7, f'{stats["MostPoopsCount"]}', 
            ha='center', va='center', fontsize=72, fontweight='bold', 
            color=SPOTIFY_GREEN)
    ax.text(0.5, 0.5, f'poops on {stats["DayWithMostPoops"]}', 
            ha='center', va='center', fontsize=32, color=SPOTIFY_WHITE)
    ax.text(0.5, 0.3, 'Your wildest day', 
            ha='center', va='center', fontsize=24, color=SPOTIFY_WHITE, style='italic')
    
    output_path = os.path.join(output_dir, 'slide_04_extreme.png')
    plt.savefig(output_path, dpi=DPI, facecolor=SPOTIFY_DARK, bbox_inches='tight')
    plt.close()
    return output_path

def generate_ranking_slide(stats, output_dir):
    """Generate ranking slide"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100), facecolor=SPOTIFY_DARK)
    ax.set_facecolor(SPOTIFY_DARK)
    ax.axis('off')
    
    if stats.get("GroupRank"):
        rank = stats["GroupRank"]
        ax.text(0.5, 0.7, f'#{rank["Rank"]}', 
                ha='center', va='center', fontsize=72, fontweight='bold', 
                color=SPOTIFY_GREEN)
        ax.text(0.5, 0.5, f'out of {rank["TotalUsers"]} poopers', 
                ha='center', va='center', fontsize=32, color=SPOTIFY_WHITE)
    else:
        ax.text(0.5, 0.6, 'Ranking unavailable', 
                ha='center', va='center', fontsize=32, color=SPOTIFY_WHITE)
    
    output_path = os.path.join(output_dir, 'slide_05_ranking.png')
    plt.savefig(output_path, dpi=DPI, facecolor=SPOTIFY_DARK, bbox_inches='tight')
    plt.close()
    return output_path

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 generate_personal_wrapped.py <stats_json_file>")
        sys.exit(1)
    
    stats_file = sys.argv[1]
    
    # Read stats
    with open(stats_file, 'r') as f:
        stats = json.load(f)
    
    # Create output directory
    output_dir = os.path.join(os.path.dirname(stats_file), f'user_{stats["UserID"]}_slides')
    os.makedirs(output_dir, exist_ok=True)
    
    # Generate slides
    image_paths = []
    image_paths.append(generate_title_slide(stats["Year"], output_dir))
    image_paths.append(generate_total_slide(stats, output_dir))
    image_paths.append(generate_streak_slide(stats, output_dir))
    image_paths.append(generate_extreme_day_slide(stats, output_dir))
    image_paths.append(generate_ranking_slide(stats, output_dir))
    
    # Write image paths to JSON
    images_file = stats_file.replace('_stats.json', '_images.json')
    with open(images_file, 'w') as f:
        json.dump(image_paths, f, indent=2)
    
    print(f"Generated {len(image_paths)} slides")

if __name__ == "__main__":
    main()