#!/usr/bin/env python3
"""Read-only QA for deterministic motion renders; never calls Kie.ai.

Usage:
  python qa.py --frames ./out/frames --fps 24 --duration 6 \
    --endpoint ./out/seam-endpoint.png --out ./out/qa

Requires Pillow. Input files must be ordered PNG sequence 0000.png ... N-1.png.
The separately-rendered endpoint at t=duration is for QA only, NOT the video.
"""
import argparse
import json
from pathlib import Path
from PIL import Image, ImageChops, ImageStat


def mean_difference(a, b):
    """Mean absolute channel difference for equal-sized RGB frames (0..255)."""
    if a.size != b.size:
        raise ValueError("Seam frames have inconsistent dimensions")
    rgb_a = a.convert("RGB")
    rgb_b = b.convert("RGB")
    return sum(ImageStat.Stat(ImageChops.difference(rgb_a, rgb_b)).mean) / 3.0


def audit(frames_dir, fps, duration, endpoint=None, out_dir="qa", seam_threshold=1.0):
    if fps <= 0 or duration <= 0:
        raise ValueError("fps and duration must be positive")
    expected = round(fps * duration)
    if expected < 2:
        raise ValueError("An animation needs at least two encoded frames")
    frames_dir = Path(frames_dir)
    out = Path(out_dir)
    out.mkdir(parents=True, exist_ok=True)
    files = sorted(frames_dir.glob("*.png"))
    expected_names = [f"{i:04d}.png" for i in range(expected)]
    names_ok = [p.name for p in files] == expected_names
    if len(files) != expected or not names_ok:
        result = {"status": "fail", "reason": "Frame count/naming mismatch",
                  "expected_frames": expected, "found_frames": len(files),
                  "names_ok": names_ok}
        (out / "qa-report.json").write_text(json.dumps(result, indent=2) + "\n")
        raise ValueError("Expected zero-padded sequential frames 0000.png .. N-1.png")

    first = Image.open(files[0]).convert("RGB")
    size = first.size
    all_sizes_ok = all(Image.open(f).size == size for f in files)

    # Sample 12 frames across the entire loop. The contact sheet is not a
    # substitute for checking every keyframe boundary where events occur.
    samples = sorted(set(round(i * (expected - 1) / 11) for i in range(12)))
    thumb_w = 360
    thumb_h = round(thumb_w * size[1] / size[0])
    sheet = Image.new("RGB", (4 * thumb_w, 3 * thumb_h), (15, 18, 24))
    for n, idx in enumerate(samples):
        with Image.open(files[idx]) as im:
            thumb = im.convert("RGB").resize((thumb_w, thumb_h), Image.Resampling.LANCZOS)
            sheet.paste(thumb, ((n % 4) * thumb_w, (n // 4) * thumb_h))
    sheet.save(out / "contact-sheet.jpg", quality=88)

    # A video includes frame 0..N-1. Seam correctness must compare the
    # separately requested frame at EXACT t=duration against frame 0.
    seam_diff = None
    seam_ok = None
    if endpoint is not None:
        with Image.open(endpoint) as end:
            seam_diff = mean_difference(first, end)
        seam_ok = seam_diff <= seam_threshold

    report = {
        "status": "pass" if (all_sizes_ok and names_ok and seam_ok is not False) else "fail",
        "expected_frames": expected, "encoded_frames": len(files),
        "frames_are_sequential": names_ok, "frame_dimensions": list(size),
        "dimensions_consistent": all_sizes_ok,
        "sample_frames": samples, "contact_sheet": str(out / "contact-sheet.jpg"),
        "seam_check": {
            "tested": endpoint is not None,
            "comparison": "renderAt(0) vs separate renderAt(duration); not frame N-1",
            "mean_absolute_pixel_difference": seam_diff,
            "threshold": seam_threshold, "pass": seam_ok
        },
        "remaining_manual_checks": [
            "Review event frames and +/- 1 and 2 frames",
            "Inspect text safe areas and contrast in every shot",
            "Review encoded audio/video synchronization and codec padding",
            "Confirm Kie provider task IDs, media provenance and usage rights",
            "Record human creative and publication approval"
        ]
    }
    (out / "qa-report.json").write_text(json.dumps(report, indent=2) + "\n")
    return report


if __name__ == "__main__":
    p = argparse.ArgumentParser()
    p.add_argument("--frames", required=True)
    p.add_argument("--fps", required=True, type=float)
    p.add_argument("--duration", required=True, type=float)
    p.add_argument("--endpoint", default=None)
    p.add_argument("--out", default="qa")
    p.add_argument("--seam-threshold", default=1.0, type=float)
    args = p.parse_args()
    result = audit(args.frames, args.fps, args.duration, args.endpoint,
                   args.out, args.seam_threshold)
    print(json.dumps(result, indent=2))
    if result["status"] != "pass":
        raise SystemExit(1)
