#! /bin/bash
# Build script for Cloudflare pages

python -m pip install mkdocs mkdocs-material 
python -m mkdocs build --clean