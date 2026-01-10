#!/usr/bin/env python3

import argparse
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable

from playwright.sync_api import sync_playwright


@dataclass(frozen=True)
class Viewport:
    name: str
    width: int
    height: int


VIEWPORTS: tuple[Viewport, ...] = (
    Viewport("desktop-1920x1080", 1920, 1080),
    Viewport("desktop-1366x768", 1366, 768),
    Viewport("desktop-1536x864", 1536, 864),
    Viewport("desktop-1440x900", 1440, 900),
    Viewport("desktop-1280x720", 1280, 720),
    Viewport("tablet-1024x768", 1024, 768),
    Viewport("tablet-834x1194", 834, 1194),
    Viewport("mobile-430x932", 430, 932),
    Viewport("mobile-414x896", 414, 896),
    Viewport("mobile-390x844", 390, 844),
    Viewport("mobile-375x667", 375, 667),
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Capture screenshots for popular viewport sizes using Playwright.",
    )
    parser.add_argument(
        "--url",
        default="http://localhost:8000",
        help="Base URL of the running Asteroids web build.",
    )
    parser.add_argument(
        "--output-dir",
        default="screenshots",
        help="Directory to write screenshots into.",
    )
    parser.add_argument(
        "--wait-ms",
        type=int,
        default=2000,
        help="Milliseconds to wait after page load to allow the canvas to render.",
    )
    return parser.parse_args()


def capture_screenshots(url: str, output_dir: Path, wait_ms: int) -> Iterable[Path]:
    output_dir.mkdir(parents=True, exist_ok=True)
    written: list[Path] = []

    with sync_playwright() as playwright:
        browser = playwright.chromium.launch()
        try:
            for viewport in VIEWPORTS:
                context = browser.new_context(
                    viewport={"width": viewport.width, "height": viewport.height},
                    device_scale_factor=1,
                )
                page = context.new_page()
                page.goto(url, wait_until="domcontentloaded")
                try:
                    page.wait_for_selector("canvas", timeout=5000)
                except Exception:
                    pass
                page.wait_for_timeout(wait_ms)
                output_path = output_dir / f"{viewport.name}.png"
                page.screenshot(path=str(output_path), full_page=True)
                written.append(output_path)
                context.close()
        finally:
            browser.close()

    return written


def main() -> None:
    args = parse_args()
    output_dir = Path(args.output_dir)
    written = capture_screenshots(args.url, output_dir, args.wait_ms)
    for path in written:
        print(path)


if __name__ == "__main__":
    main()
