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
from matplotlib.colors import LinearSegmentedColormap, BoundaryNorm
from matplotlib.patheffects import withStroke
from matplotlib.offsetbox import OffsetImage, AnnotationBbox
from PIL import Image
import numpy as np
import datetime

# New Visual Identity - Vibrant & Polished
# Color Palette
BG_GRADIENT_START = '#0f0c29'  # Deep dark purple-blue
BG_GRADIENT_END = '#1a1a2e'    # Dark blue-purple
ACCENT_PRIMARY = '#6366f1'     # Indigo (main accent)
ACCENT_SECONDARY = '#8b5cf6'  # Purple
ACCENT_TERTIARY = '#ec4899'   # Pink
ACCENT_HIGHLIGHT = '#f59e0b'  # Amber (for highlights)
ACCENT_SUCCESS = '#10b981'     # Emerald (for positive stats)
TEXT_WHITE = '#FFFFFF'
TEXT_LIGHT = '#E5E7EB'
TEXT_MUTED = '#9CA3AF'
SLIDE_WIDTH = 1200
SLIDE_HEIGHT = 800
DPI = 150
MARGIN = 0.1  # Margin for all slides (increased for better spacing)
EMOJI_PATH = os.path.join(os.path.dirname(__file__), '..', '..', 'utils', 'Poop_Emoji.webp')
MARGIN = 0.1  # Margin for all slides (increased for better spacing)
EMOJI_PATH = os.path.join(os.path.dirname(__file__), '..', '..', 'utils', 'Poop_Emoji.webp')

# Font configuration - more playful
rcParams['font.family'] = 'sans-serif'
rcParams['font.sans-serif'] = ['Arial', 'Helvetica', 'DejaVu Sans']

def create_gradient_background(ax, color1, color2, xlim=None, ylim=None):
    """Create a gradient background"""
    if xlim is None:
        xlim = ax.get_xlim()
    if ylim is None:
        ylim = ax.get_ylim()
    
    x = np.linspace(xlim[0], xlim[1], 100)
    y = np.linspace(ylim[0], ylim[1], 100)
    X, Y = np.meshgrid(x, y)
    Z = Y  # Gradient from top to bottom
    
    # Create gradient
    colors = [color1, color2]
    n_bins = 100
    cmap = LinearSegmentedColormap.from_list('gradient', colors, N=n_bins)
    ax.imshow(Z, extent=[xlim[0], xlim[1], ylim[0], ylim[1]], aspect='auto', cmap=cmap, alpha=1.0, zorder=0)

def add_rounded_rect(ax, x, y, width, height, color, alpha=0.15, zorder=1, edgecolor=None, linewidth=0):
    """Add a rounded rectangle with proper padding"""
    rect = patches.FancyBboxPatch(
        (x - width/2, y - height/2), width, height,
        boxstyle="round,pad=0.015", 
        facecolor=color, 
        edgecolor=edgecolor if edgecolor else 'none',
        linewidth=linewidth,
        alpha=alpha,
        zorder=zorder
    )
    ax.add_patch(rect)
    return rect

def generate_title_slide(year, output_dir):
    """Generate title slide with gradient and polished design"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100))
    ax.set_xlim(MARGIN, 1 - MARGIN)
    ax.set_ylim(MARGIN, 1 - MARGIN)
    ax.axis('off')
    
    # Gradient background
    create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END, 
                              xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
    
    # Subtle rounded container for title (with proper padding, centered vertically)
    add_rounded_rect(ax, 0.5, 0.6, 0.6, 0.18, ACCENT_PRIMARY, alpha=0.2, zorder=2)
    
    # Title with stroke effect (centered within margins)
    title_text = ax.text(0.5, 0.6, f'Your Poop Wrapped {year}', 
            ha='center', va='center', fontsize=50, fontweight='bold', 
            color=TEXT_WHITE, zorder=3)
    title_text.set_path_effects([withStroke(linewidth=3, foreground=BG_GRADIENT_START)])
    
    # Emoji image
    if os.path.exists(EMOJI_PATH):
        try:
            emoji_img = Image.open(EMOJI_PATH)
            # Convert to RGBA if needed
            if emoji_img.mode != 'RGBA':
                emoji_img = emoji_img.convert('RGBA')
            # Resize emoji to appropriate size (about 15% of slide height)
            emoji_size_pixels = int(SLIDE_HEIGHT * 0.15)
            emoji_img.thumbnail((emoji_size_pixels, emoji_size_pixels), Image.Resampling.LANCZOS)
            im = OffsetImage(emoji_img, zoom=1.0)
            ab = AnnotationBbox(im, (0.5, 0.4), frameon=False, box_alignment=(0.5, 0.5), zorder=3)
            ax.add_artist(ab)
        except Exception as e:
            # Fallback to text if image fails
            ax.text(0.5, 0.4, '💩', ha='center', va='center', fontsize=120, zorder=3)
    else:
        ax.text(0.5, 0.4, '💩', ha='center', va='center', fontsize=120, zorder=3)
    
    output_path = os.path.join(output_dir, 'slide_01_title.png')
    plt.savefig(output_path, dpi=DPI, bbox_inches='tight', facecolor=BG_GRADIENT_START)
    plt.close()
    return output_path

def generate_total_slide(stats, output_dir):
    """Generate total count slide with polished design"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100))
    ax.set_xlim(MARGIN, 1 - MARGIN)
    ax.set_ylim(MARGIN, 1 - MARGIN)
    ax.axis('off')
    
    # Gradient background
    create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END,
                              xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
    
    # Rounded card for the number (centered vertically)
    add_rounded_rect(ax, 0.5, 0.65, 0.35, 0.18, ACCENT_PRIMARY, alpha=0.2, zorder=2)
    
    # Total count
    count_text = ax.text(0.5, 0.65, f'{stats["TotalPoops"]}', 
            ha='center', va='center', fontsize=112, fontweight='bold', 
            color=ACCENT_HIGHLIGHT, zorder=3)
    count_text.set_path_effects([withStroke(linewidth=4, foreground=BG_GRADIENT_START)])
    
    # Label
    ax.text(0.5, 0.45, 'Total dumps', 
            ha='center', va='center', fontsize=40, color=TEXT_WHITE, 
            fontweight='bold', zorder=3)
    
    # Comparison if available
    if stats.get("GroupRank"):
        rank = stats["GroupRank"]
        percentage = rank["Percentage"]
        add_rounded_rect(ax, 0.5, 0.25, 0.55, 0.1, ACCENT_SECONDARY, alpha=0.15, zorder=2)
        ax.text(0.5, 0.25, f'That\'s more than {percentage:.1f}% of the group', 
                ha='center', va='center', fontsize=24, color=TEXT_LIGHT, 
                style='italic', zorder=3)
    
    output_path = os.path.join(output_dir, 'slide_02_total.png')
    plt.savefig(output_path, dpi=DPI, bbox_inches='tight', facecolor=BG_GRADIENT_START)
    plt.close()
    return output_path

def generate_streak_slide(stats, output_dir):
    """Generate streak slide"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100))
    ax.set_xlim(MARGIN, 1 - MARGIN)
    ax.set_ylim(MARGIN, 1 - MARGIN)
    ax.axis('off')
    
    # Gradient background
    create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END,
                              xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
    
    # Rounded card
    add_rounded_rect(ax, 0.5, 0.7, 0.3, 0.16, ACCENT_TERTIARY, alpha=0.2, zorder=2)
    
    streak_text = ax.text(0.5, 0.7, f'{stats["MaxStreak"]}', 
            ha='center', va='center', fontsize=112, fontweight='bold', 
            color=ACCENT_HIGHLIGHT, zorder=3)
    streak_text.set_path_effects([withStroke(linewidth=4, foreground=BG_GRADIENT_START)])
    
    ax.text(0.5, 0.5, 'Longest streak', 
            ha='center', va='center', fontsize=40, color=TEXT_WHITE, 
            fontweight='bold', zorder=3)
    
    add_rounded_rect(ax, 0.5, 0.3, 0.5, 0.1, ACCENT_SECONDARY, alpha=0.15, zorder=2)
    ax.text(0.5, 0.3, 'Your bowels respect routine', 
            ha='center', va='center', fontsize=24, color=TEXT_LIGHT, 
            style='italic', zorder=3)
    
    output_path = os.path.join(output_dir, 'slide_03_streak.png')
    plt.savefig(output_path, dpi=DPI, bbox_inches='tight', facecolor=BG_GRADIENT_START)
    plt.close()
    return output_path

def generate_extreme_day_slide(stats, output_dir):
    """Generate extreme day slide"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100))
    ax.set_xlim(MARGIN, 1 - MARGIN)
    ax.set_ylim(MARGIN, 1 - MARGIN)
    ax.axis('off')
    
    # Gradient background
    create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END,
                              xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
    
    # Rounded card for number
    add_rounded_rect(ax, 0.5, 0.75, 0.25, 0.14, ACCENT_TERTIARY, alpha=0.2, zorder=2)
    
    count_text = ax.text(0.5, 0.75, f'{stats["MostPoopsCount"]}', 
            ha='center', va='center', fontsize=112, fontweight='bold', 
            color=ACCENT_HIGHLIGHT, zorder=3)
    count_text.set_path_effects([withStroke(linewidth=4, foreground=BG_GRADIENT_START)])
    
    add_rounded_rect(ax, 0.5, 0.55, 0.6, 0.1, ACCENT_SECONDARY, alpha=0.15, zorder=2)
    ax.text(0.5, 0.55, f'poops on {stats["DayWithMostPoops"]}', 
            ha='center', va='center', fontsize=32, color=TEXT_WHITE, 
            fontweight='bold', zorder=3)
    
    ax.text(0.5, 0.35, 'Your wildest day', 
            ha='center', va='center', fontsize=28, color=TEXT_LIGHT, 
            style='italic', zorder=3)
    
    output_path = os.path.join(output_dir, 'slide_04_extreme.png')
    plt.savefig(output_path, dpi=DPI, bbox_inches='tight', facecolor=BG_GRADIENT_START)
    plt.close()
    return output_path

def generate_ranking_slide(stats, output_dir):
    """Generate ranking slide"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100))
    ax.set_xlim(MARGIN, 1 - MARGIN)
    ax.set_ylim(MARGIN, 1 - MARGIN)
    ax.axis('off')
    
    # Gradient background
    create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END,
                              xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
    
    if stats.get("GroupRank"):
        rank = stats["GroupRank"]
        # Rounded card for rank
        add_rounded_rect(ax, 0.5, 0.7, 0.22, 0.16, ACCENT_PRIMARY, alpha=0.2, zorder=2)
        
        rank_text = ax.text(0.5, 0.7, f'#{rank["Rank"]}', 
                ha='center', va='center', fontsize=112, fontweight='bold', 
                color=ACCENT_HIGHLIGHT, zorder=3)
        rank_text.set_path_effects([withStroke(linewidth=4, foreground=BG_GRADIENT_START)])
        
        add_rounded_rect(ax, 0.5, 0.45, 0.45, 0.1, ACCENT_SECONDARY, alpha=0.15, zorder=2)
        ax.text(0.5, 0.45, f'out of {rank["TotalUsers"]} poopers', 
                ha='center', va='center', fontsize=32, color=TEXT_WHITE, 
                fontweight='bold', zorder=3)
    else:
        ax.text(0.5, 0.5, 'Ranking unavailable', 
                ha='center', va='center', fontsize=32, color=TEXT_WHITE)
    
    output_path = os.path.join(output_dir, 'slide_05_ranking.png')
    plt.savefig(output_path, dpi=DPI, bbox_inches='tight', facecolor=BG_GRADIENT_START)
    plt.close()
    return output_path

def generate_personality_slide(stats, output_dir):
    """Generate personality slide"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100))
    ax.set_xlim(MARGIN, 1 - MARGIN)
    ax.set_ylim(MARGIN, 1 - MARGIN)
    ax.axis('off')
    
    # Gradient background
    create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END,
                              xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
    
    personality = stats.get("Personality", {})
    personality_type = personality.get("Type", "Unknown")
    description = personality.get("Description", "")
    
    description_lines = description.split('\n')
    
    # Rounded card for personality type (centered vertically, with proper margins)
    add_rounded_rect(ax, 0.5, 0.7, 0.55, 0.12, ACCENT_PRIMARY, alpha=0.2, zorder=2)
    
    personality_text = ax.text(0.5, 0.7, personality_type, 
            ha='center', va='center', fontsize=56, fontweight='bold', 
            color=ACCENT_HIGHLIGHT, zorder=3)
    personality_text.set_path_effects([withStroke(linewidth=4, foreground=BG_GRADIENT_START)])
    
    # Display main description - wrap if too long
    if len(description_lines) > 0:
        main_desc = description_lines[0]
        # Break long descriptions into multiple lines (max ~40 chars per line)
        if len(main_desc) > 40:
            words = main_desc.split()
            wrapped_lines = []
            current_line = []
            current_len = 0
            for word in words:
                if current_len + len(word) + 1 > 40 and current_line:
                    wrapped_lines.append(' '.join(current_line))
                    current_line = [word]
                    current_len = len(word)
                else:
                    current_line.append(word)
                    current_len += len(word) + 1
            if current_line:
                wrapped_lines.append(' '.join(current_line))
            main_desc = '\n'.join(wrapped_lines)
        
        # Adjust container height based on number of lines
        num_lines = main_desc.count('\n') + 1
        container_height = 0.06 + (num_lines - 1) * 0.03
        
        add_rounded_rect(ax, 0.5, 0.5, 0.55, container_height, ACCENT_SECONDARY, alpha=0.15, zorder=2)
        ax.text(0.5, 0.5, main_desc, 
                ha='center', va='center', fontsize=24, color=TEXT_WHITE, 
                fontweight='bold', zorder=3)
    
    # Display subtext (explanation) if present - ensure it doesn't touch bottom (with margin)
    if len(description_lines) > 1:
        subtext = description_lines[1]
        # Wrap subtext if too long (max ~50 chars per line)
        if len(subtext) > 50:
            words = subtext.split()
            wrapped_lines = []
            current_line = []
            current_len = 0
            for word in words:
                if current_len + len(word) + 1 > 50 and current_line:
                    wrapped_lines.append(' '.join(current_line))
                    current_line = [word]
                    current_len = len(word)
                else:
                    current_line.append(word)
                    current_len += len(word) + 1
            if current_line:
                wrapped_lines.append(' '.join(current_line))
            subtext = '\n'.join(wrapped_lines)
        
        # Position text above the bottom margin to avoid touching edges
        num_lines = subtext.count('\n') + 1
        y_pos = MARGIN + 0.1 + (num_lines - 1) * 0.03
        ax.text(0.5, y_pos, subtext, 
                ha='center', va='center', fontsize=20, color=TEXT_LIGHT, 
                style='italic', zorder=3)
    
    output_path = os.path.join(output_dir, 'slide_06_personality.png')
    plt.savefig(output_path, dpi=DPI, bbox_inches='tight', facecolor=BG_GRADIENT_START)
    plt.close()
    return output_path

def generate_time_pattern_slide(stats, output_dir):
    """Generate hour-of-day histogram with polished design"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100))
    ax.set_xlim(MARGIN, 1 - MARGIN)
    ax.set_ylim(MARGIN, 1 - MARGIN)
    
    hour_dist = stats.get("HourDistribution", [])
    if not hour_dist:
        ax.axis('off')
        create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END,
                                  xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
        ax.text(0.5, 0.5, 'No time pattern data available', 
                ha='center', va='center', fontsize=32, color=TEXT_WHITE)
    else:
        # Create gradient background
        create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END,
                                  xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
        
        # Rounded container for chart
        add_rounded_rect(ax, 0.5, 0.5, 0.75, 0.55, ACCENT_PRIMARY, alpha=0.1, zorder=1)
        
        # Create subplot for chart - even larger size
        chart_ax = fig.add_axes([0.06 + MARGIN, 0.12 + MARGIN, 0.88 - MARGIN*2, 0.76 - MARGIN*2])
        chart_ax.set_facecolor('none')
        
        hours = [h["Hour"] for h in hour_dist]
        counts = [h["PoopCount"] for h in hour_dist]
        
        # Create bars with consistent gradient color
        max_count = max(counts) if counts else 1
        bars = chart_ax.bar(hours, counts, 
                           color=ACCENT_PRIMARY, alpha=0.7, 
                           edgecolor=ACCENT_SECONDARY, linewidth=1.2)
        
        chart_ax.set_xlabel('Hour of Day', color=TEXT_WHITE, fontsize=18, fontweight='bold')
        chart_ax.set_ylabel('Poop Count', color=TEXT_WHITE, fontsize=18, fontweight='bold')
        chart_ax.set_title('Your Poop Rhythm', color=TEXT_WHITE, fontsize=28, 
                          fontweight='bold', pad=20)
        chart_ax.set_xticks(range(0, 24, 2))
        chart_ax.tick_params(colors=TEXT_WHITE, labelsize=12)
        chart_ax.spines['bottom'].set_color(TEXT_WHITE)
        chart_ax.spines['left'].set_color(TEXT_WHITE)
        chart_ax.spines['top'].set_visible(False)
        chart_ax.spines['right'].set_visible(False)
        chart_ax.set_facecolor('none')
        chart_ax.grid(True, alpha=0.15, color=TEXT_WHITE, linestyle='--', axis='y')
    
    output_path = os.path.join(output_dir, 'slide_07_time_pattern.png')
    plt.savefig(output_path, dpi=DPI, bbox_inches='tight', facecolor=BG_GRADIENT_START)
    plt.close()
    return output_path

def generate_monthly_chart(stats, output_dir):
    """Generate monthly distribution bar chart"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100))
    ax.set_xlim(MARGIN, 1 - MARGIN)
    ax.set_ylim(MARGIN, 1 - MARGIN)
    
    monthly_stats = stats.get("MonthlyStats", [])
    if not monthly_stats:
        ax.axis('off')
        create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END,
                                  xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
        ax.text(0.5, 0.5, 'No monthly data available', 
                ha='center', va='center', fontsize=32, color=TEXT_WHITE)
    else:
        # Gradient background
        create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END,
                                  xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
        
        # Rounded container
        add_rounded_rect(ax, 0.5, 0.5, 0.75, 0.55, ACCENT_SECONDARY, alpha=0.1, zorder=1)
        
        # Create subplot for chart - even larger size
        chart_ax = fig.add_axes([0.06 + MARGIN, 0.12 + MARGIN, 0.88 - MARGIN*2, 0.76 - MARGIN*2])
        chart_ax.set_facecolor('none')
        
        # Create a map of month number to count
        month_map = {}
        for ms in monthly_stats:
            month_str = ms.get("Month", "")
            if len(month_str) >= 7:
                month_num = int(month_str[5:7])  # Extract MM from YYYY-MM
                month_map[month_num] = ms.get("PoopCount", 0)
        
        # Create arrays for all 12 months
        months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 
                  'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
        monthly_counts = [month_map.get(i+1, 0) for i in range(12)]
        
        # Use consistent gradient color scheme
        bars = chart_ax.bar(months, monthly_counts, color=ACCENT_SECONDARY, alpha=0.7, 
                           edgecolor=ACCENT_PRIMARY, linewidth=1.2)
        
        chart_ax.set_xlabel('Month', color=TEXT_WHITE, fontsize=18, fontweight='bold')
        chart_ax.set_ylabel('Poop Count', color=TEXT_WHITE, fontsize=18, fontweight='bold')
        chart_ax.set_title('Monthly Breakdown', color=TEXT_WHITE, fontsize=28, 
                          fontweight='bold', pad=20)
        chart_ax.tick_params(colors=TEXT_WHITE, labelsize=12)
        chart_ax.spines['bottom'].set_color(TEXT_WHITE)
        chart_ax.spines['left'].set_color(TEXT_WHITE)
        chart_ax.spines['top'].set_visible(False)
        chart_ax.spines['right'].set_visible(False)
        chart_ax.set_facecolor('none')
        chart_ax.grid(True, alpha=0.15, color=TEXT_WHITE, linestyle='--', axis='y')
    
    output_path = os.path.join(output_dir, 'slide_08_monthly.png')
    plt.savefig(output_path, dpi=DPI, bbox_inches='tight', facecolor=BG_GRADIENT_START)
    plt.close()
    return output_path

def generate_heatmap(stats, output_dir):
    """Generate heatmap with GitHub-style green colors"""
    fig, ax = plt.subplots(figsize=(SLIDE_WIDTH/100, SLIDE_HEIGHT/100))
    
    week_day_stats = stats.get("WeekDayStats", [])
    year = stats.get("Year", 2025)
    
    # Calculate total_weeks first (needed for axis limits)
    if week_day_stats:
        jan_1 = datetime.date(year, 1, 1)
        dec_31 = datetime.date(year, 12, 31)
        days_until_sunday = (6 - jan_1.weekday()) % 7
        first_sunday = jan_1 + datetime.timedelta(days=days_until_sunday)
        days_since_first_sunday_dec31 = (dec_31 - first_sunday).days
        last_week = (days_since_first_sunday_dec31 // 7) + 1 if days_since_first_sunday_dec31 >= 0 else 1
        total_weeks = last_week
    else:
        total_weeks = 52
    
    if not week_day_stats:
        ax.set_xlim(MARGIN, 1 - MARGIN)
        ax.set_ylim(MARGIN, 1 - MARGIN)
        ax.axis('off')
        create_gradient_background(ax, BG_GRADIENT_START, BG_GRADIENT_END,
                                  xlim=(MARGIN, 1 - MARGIN), ylim=(MARGIN, 1 - MARGIN))
        ax.text(0.5, 0.5, 'No heatmap data available', 
                ha='center', va='center', fontsize=32, color=TEXT_WHITE)
    else:
        # Set axis limits first
        ax.set_xlim(-0.5, total_weeks + 5)
        ax.set_ylim(-0.5, 7)
        
        # Dark background for heatmap (GitHub-style)
        ax.set_facecolor('#0d1117')  # GitHub dark background
        
        # Create a 2D grid: 7 days (rows) x total_weeks (columns)
        heatmap_data = np.zeros((7, total_weeks))
        
        # Fill in the data
        max_count = 0
        for wd in week_day_stats:
            week_num = wd.get("WeekNumber", 0)
            day_of_week = wd.get("DayOfWeek", 0)
            count = wd.get("PoopCount", 0)
            if week_num > 0 and week_num <= total_weeks and day_of_week >= 0 and day_of_week < 7:
                heatmap_data[day_of_week, week_num - 1] = count
                max_count = max(max_count, count)
        
        # GitHub-style green color scheme
        def get_color_for_count(count):
            if count == 0:
                return '#161b22'  # Dark (matches GitHub)
            elif count == 1:
                return '#0e4429'  # Very dark green
            elif count <= 3:
                return '#006d32'  # Medium green
            elif count <= 5:
                return '#26a641'  # Bright green
            else:
                return '#39d353'  # Lightest green
        
        # Draw individual squares with subtle rounded corners
        square_size = 0.85
        for day in range(7):
            for week in range(total_weeks):
                count = heatmap_data[day, week]
                color = get_color_for_count(count)
                
                # Draw rounded rectangle
                rect = patches.FancyBboxPatch(
                    (week, day), square_size, square_size,
                    boxstyle="round,pad=0.02",
                    facecolor=color, 
                    edgecolor='#ffffff', 
                    linewidth=0.15,
                    alpha=0.95,
                    zorder=2
                )
                ax.add_patch(rect)
        
        # Set title
        ax.set_title('Your Poop Activity Heatmap', color=TEXT_WHITE, fontsize=28, 
                    fontweight='bold', pad=20)
        
        # Set day labels on the left
        day_labels = ['S', 'M', 'T', 'W', 'T', 'F', 'S']
        ax.set_yticks(np.arange(7) + 0.5)
        ax.set_yticklabels(day_labels, color=TEXT_WHITE, fontsize=12, fontweight='bold')
        
        # Calculate month positions
        month_labels = []
        month_positions = []
        months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 
                  'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
        
        seen_weeks = set()
        for month_idx in range(12):
            first_day = datetime.date(year, month_idx + 1, 1)
            jan_1 = datetime.date(year, 1, 1)
            days_until_sunday = (6 - jan_1.weekday()) % 7
            first_sunday = jan_1 + datetime.timedelta(days=days_until_sunday)
            
            days_since_first_sunday = (first_day - first_sunday).days
            if days_since_first_sunday >= 0:
                week_pos = days_since_first_sunday // 7
            else:
                week_pos = 0
            
            if week_pos < total_weeks and week_pos not in seen_weeks:
                month_labels.append(months[month_idx])
                month_positions.append(week_pos + 0.5)
                seen_weeks.add(week_pos)
        
        # Set x-axis labels
        ax.set_xticks(month_positions)
        ax.set_xticklabels(month_labels, color=TEXT_WHITE, fontsize=12, 
                          rotation=0, ha='left', fontweight='bold')
        
        # Add vertical lines at month boundaries
        for week_pos in month_positions:
            if week_pos > 0:
                ax.axvline(x=week_pos - 0.5, color=ACCENT_SECONDARY, linewidth=1.2, 
                          alpha=0.4, linestyle='-', zorder=0)
        
        # Set axis limits
        ax.set_xlim(-0.5, total_weeks + 5)
        ax.set_ylim(-0.5, 7)
        ax.invert_yaxis()
        
        # Remove spines
        ax.spines['top'].set_visible(False)
        ax.spines['right'].set_visible(False)
        ax.spines['bottom'].set_visible(False)
        ax.spines['left'].set_visible(False)
        ax.tick_params(axis='x', which='both', bottom=False, top=False, length=0)
        ax.tick_params(axis='y', which='both', left=True, right=False, length=0, labelsize=12)
        
        # Add color legend
        color_legend_x = total_weeks + 1.5
        legend_y_start = 2.0
        
        legend_colors = ['#161b22', '#0e4429', '#006d32', '#26a641', '#39d353']
        legend_labels = ['0', '1', '2-3', '4-5', '6+']
        
        ax.text(color_legend_x, legend_y_start + 2.0, 'Poops', color=TEXT_WHITE, 
                fontsize=13, ha='left', weight='bold')
        
        for i, (color, label) in enumerate(zip(reversed(legend_colors), reversed(legend_labels))):
            y_pos = legend_y_start + 1.3 - i * 0.35
            rect = patches.FancyBboxPatch(
                (color_legend_x, y_pos - 0.12), 0.35, 0.35,
                boxstyle="round,pad=0.02",
                facecolor=color, 
                edgecolor='#ffffff', 
                linewidth=0.15
            )
            ax.add_patch(rect)
            ax.text(color_legend_x + 0.5, y_pos, label, color=TEXT_WHITE, 
                   fontsize=10, va='center', ha='left', fontweight='bold')
    
    output_path = os.path.join(output_dir, 'slide_09_heatmap.png')
    plt.savefig(output_path, dpi=DPI, bbox_inches='tight', facecolor=BG_GRADIENT_START)
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
    image_paths.append(generate_personality_slide(stats, output_dir))
    image_paths.append(generate_time_pattern_slide(stats, output_dir))
    image_paths.append(generate_monthly_chart(stats, output_dir))
    image_paths.append(generate_heatmap(stats, output_dir))
    
    # Write image paths to JSON
    images_file = stats_file.replace('_stats.json', '_images.json')
    with open(images_file, 'w') as f:
        json.dump(image_paths, f, indent=2)
    
    print(f"Generated {len(image_paths)} slides")

if __name__ == "__main__":
    main()
