#!/usr/bin/env bash
# Creating an avatar or a purchase uploads an image to a local file system
# these images are located at /images/ folder.
# This script removes all temporary uploaded files and all avatars and purchases
cd images/avatars/b/
rm *.jpg
rm *.png

cd ../s
rm *.jpg
rm *.png

cd ../../tmp
find . -type f  ! -name "*.*"  -delete