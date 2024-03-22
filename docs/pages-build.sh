#! /bin/bash
# Build script for Cloudflare pages

python -m pip install mkdocs mkdocs-material mkdocs-table-reader-plugin
python -m mkdocs build --clean