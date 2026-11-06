#!/usr/bin/env python3
"""
Visualise a list of 3D direction vectors read from stdin.

Input format (any whitespace/newlines allowed):
    [(0.69,0.69,0.195),(-0.69,0.69,0.195),...]

Examples:
    cat vectors.txt | ./plot_vectors.py                 # renders vectors.png
    cat vectors.txt | ./plot_vectors.py -i              # interactive 3D window
    cat vectors.txt | ./plot_vectors.py -i --arrows     # quiver arrows instead of lines
    echo "[(1,0,0),(0,1,0)]" | ./plot_vectors.py -o xy.png --color-by elevation
"""

import argparse
import ast
import math
import re
import sys

import numpy as np
import matplotlib

NUM_RE = re.compile(r"[-+]?\d*\.?\d+(?:[eE][-+]?\d+)?")


# ---------------------------------------------------------------- input parsing

def parse_vectors(text: str) -> np.ndarray:
    """Turn the stdin blob into an (N, 3) float array."""
    text = text.strip()
    if not text:
        raise ValueError("stdin is empty")

    try:
        data = ast.literal_eval(text)
    except (ValueError, SyntaxError):
        # Fallback: pick every number and regroup them by threes.
        nums = [float(x) for x in NUM_RE.findall(text)]
        if len(nums) < 3 or len(nums) % 3:
            raise ValueError("could not parse the input as triples of numbers")
        data = np.asarray(nums, dtype=float).reshape(-1, 3)

    if isinstance(data, (int, float)):
        raise ValueError("expected a list of 3D vectors, got a single number")

    arr = np.asarray(data, dtype=float)

    if arr.ndim == 1:                                  # flat list of numbers
        if arr.size < 3 or arr.size % 3:
            raise ValueError("flat input must contain a multiple of 3 numbers")
        arr = arr.reshape(-1, 3)
    if arr.ndim != 2 or arr.shape[1] != 3:
        raise ValueError(f"expected shape (N, 3), got {arr.shape}")
    if not np.isfinite(arr).all():
        raise ValueError("input contains NaN or infinite values")

    return arr


# ------------------------------------------------------------------- geometry

def rotation_from_z(axis: np.ndarray) -> np.ndarray:
    """Rodrigues rotation mapping +z onto the unit vector `axis`."""
    z = np.array([0.0, 0.0, 1.0])
    v = np.cross(z, axis)
    c = float(np.dot(z, axis))
    if np.linalg.norm(v) < 1e-12:
        return np.eye(3) if c > 0 else np.diag([1.0, -1.0, -1.0])
    vx = np.array([[0.0, -v[2], v[1]],
                   [v[2], 0.0, -v[0]],
                   [-v[1], v[0], 0.0]])
    return np.eye(3) + vx + vx @ vx * (1.0 / (1.0 + c))


def hemisphere(pole: np.ndarray, n_phi: int = 48, n_theta: int = 12):
    """Mesh of the hemisphere whose pole is `pole`, plus its equator ring."""
    theta = np.linspace(0.0, math.pi / 2.0, n_theta)
    phi = np.linspace(0.0, 2.0 * math.pi, n_phi)
    T, P = np.meshgrid(theta, phi)
    pts = np.stack([np.sin(T) * np.cos(P),
                    np.sin(T) * np.sin(P),
                    np.cos(T)], axis=-1)
    pts = pts @ rotation_from_z(pole).T

    eq = np.stack([np.cos(phi), np.sin(phi), np.zeros_like(phi)], axis=-1)
    return pts, eq @ rotation_from_z(pole).T


# ----------------------------------------------------------------------- main

def main() -> int:
    ap = argparse.ArgumentParser(
        description="Plot 3D direction vectors (read from stdin) on a hemisphere.")
    ap.add_argument("-i", "--interactive", action="store_true",
                    help="open an interactive window (drag to rotate, scroll to zoom)")
    ap.add_argument("-o", "--output", default="vectors.png",
                    help="image file for non-interactive mode (default: vectors.png)")
    ap.add_argument("--arrows", action="store_true",
                    help="draw quiver arrows (slower than lines above ~200 vectors)")
    ap.add_argument("--no-mesh", action="store_true",
                    help="skip the translucent hemisphere")
    ap.add_argument("--color-by", choices=("azimuth", "elevation", "index", "none"),
                    default="azimuth")
    ap.add_argument("--point-size", type=float, default=16.0)
    ap.add_argument("--width", type=float, default=9.0)
    ap.add_argument("--height", type=float, default=8.0)
    ap.add_argument("--dpi", type=int, default=150)
    ap.add_argument("--elev", type=float, default=22.0, help="camera elevation (deg)")
    ap.add_argument("--azim", type=float, default=-58.0, help="camera azimuth (deg)")
    ap.add_argument("--no-axes", action="store_true", help="clean look, no frame/labels")
    ap.add_argument("--quiet", action="store_true", help="do not print statistics")
    args = ap.parse_args()

    if not args.interactive:
        matplotlib.use("Agg")
    import matplotlib.pyplot as plt
    from matplotlib import cm
    from mpl_toolkits.mplot3d.art3d import Line3DCollection

    try:
        vecs = parse_vectors(sys.stdin.read())
    except (ValueError, SyntaxError) as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 2

    n = len(vecs)
    norms = np.linalg.norm(vecs, axis=1)
    unit = vecs / norms[:, None]

    mean = unit.mean(axis=0)
    mean_norm = np.linalg.norm(mean)
    pole = mean / mean_norm if mean_norm > 1e-9 else unit[np.argmax(unit[:, 2])]
    max_angle = math.degrees(np.arccos(np.clip(unit @ pole, -1.0, 1.0)).max())

    if not args.quiet:
        print(f"{n} vectors | norms in [{norms.min():.4f}, {norms.max():.4f}]",
              file=sys.stderr)
        print(f"mean direction: ({pole[0]:+.4f}, {pole[1]:+.4f}, {pole[2]:+.4f}) | "
              f"resultant length {mean_norm:.3f} | max angle from mean {max_angle:.1f} deg"
              + ("  -> spans more than a hemisphere" if max_angle > 90 else ""),
              file=sys.stderr)

    # ---- figure
    fig = plt.figure(figsize=(args.width, args.height))
    ax = fig.add_subplot(111, projection="3d")
    ax.set_box_aspect((1, 1, 1))

    # ---- hemisphere shell + equator
    if not args.no_mesh:
        mesh, equator = hemisphere(pole)
        ax.plot_surface(mesh[..., 0], mesh[..., 1], mesh[..., 2],
                        color="#4f8fd0", alpha=0.10, linewidth=0, antialiased=False)
        for row in (mesh[::max(1, len(mesh) // 8), :, :],
                    mesh[:, ::max(1, len(mesh[0]) // 12), :]):
            ax.plot(row[..., 0], row[..., 1], row[..., 2],
                    color="#4f8fd0", alpha=0.18, linewidth=0.5)
        ax.plot(equator[:, 0], equator[:, 1], equator[:, 2],
                color="#2b5f92", alpha=0.6, linewidth=1.2)

    # ---- mean direction
    ax.quiver(0, 0, 0, *(2.05 * pole), color="#d1495b", linewidth=2,
              arrow_length_ratio=0.08, label="mean direction")

    # ---- the vectors themselves
    if args.arrows:
        ax.quiver(np.zeros(n), np.zeros(n), np.zeros(n),
                  vecs[:, 0], vecs[:, 1], vecs[:, 2],
                  color="#333333", linewidth=0.9, arrow_length_ratio=0.06)
    else:
        segs = np.concatenate([np.zeros((n, 1, 3)), vecs[:, None, :]], axis=1)
        ax.add_collection3d(Line3DCollection(segs, colors="#4a4a4a",linewidths=0.8, alpha=0.75))

    # ---- tips
    if args.color_by == "azimuth":
        values, cmap, vmin, vmax = (np.arctan2(unit[:, 1], unit[:, 0]) / (2 * np.pi)
                                    + 0.5), cm.hsv, 0.0, 1.0
    elif args.color_by == "elevation":
        values, cmap, vmin, vmax = unit[:, 2], cm.viridis, -1.0, 1.0
    elif args.color_by == "index":
        values, cmap, vmin, vmax = np.arange(n), cm.plasma, 0, max(n - 1, 1)
    else:
        values = cmap = None

    if values is None:
        ax.scatter(vecs[:, 0], vecs[:, 1], vecs[:, 2], s=args.point_size,
                   c="#111111", depthshade=False)
    else:
        sc = ax.scatter(vecs[:, 0], vecs[:, 1], vecs[:, 2], s=args.point_size,
                        c=values, cmap=cmap, vmin=vmin, vmax=vmax,
                        depthshade=False, edgecolors="white", linewidths=0.3)
        if args.color_by != "index":
            cbar = fig.colorbar(sc, ax=ax, shrink=0.55, pad=0.08)
            cbar.set_label("azimuth (0-360 deg)" if args.color_by == "azimuth"
                           else "z component")        
    # ---- framing
    lim = float(norms.max()) * 1.12
    ax.set_xlim(-lim, lim); ax.set_ylim(-lim, lim); ax.set_zlim(-lim, lim)
    ax.set_xlabel("x"); ax.set_ylabel("y"); ax.set_zlabel("z")
    ax.set_title(f"{n} directions  |  max angle from mean {max_angle:.1f}°")
    ax.view_init(elev=args.elev, azim=args.azim)
    if args.no_axes:
        ax.set_axis_off()
    ax.legend(loc="upper left", frameon=False)

    fig.tight_layout()

    if args.interactive:
        ax.set_title(f"{n} directions  |  drag to rotate, scroll to zoom", pad=14)
        plt.show()
    else:
        fig.savefig(args.output, dpi=args.dpi, bbox_inches="tight")
        print(f"saved: {args.output}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main())
