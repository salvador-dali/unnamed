#!/usr/bin/env bash
# Creating an avatar or a purchase uploads an image to a local file system
# these images are located at /images/ folder.
# This script removes all temporary uploaded files and all avatars and purchases

# cryptic line means: Remove all files which do not end with .md extension or is not for tests

# big avatars
cd images/avatars/b/
find . -type f  ! -name "*.md" ! -name "*isForTests.jpg" -delete

# small avatars
cd ../s/
find . -type f  ! -name "*.md" ! -name "*isForTests.jpg" -delete

# big purchases
cd ../../purchases/b/
find . -type f  ! -name "*.md" ! -name "*isForTests.jpg" -delete

# medium purchases
cd ../m/
find . -type f  ! -name "*.md" ! -name "*isForTests.jpg" -delete

# temporary files
cd ../../tmp
find . -type f  ! -name "*.md"  -delete