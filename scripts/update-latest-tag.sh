#!/usr/bin/env bash

git tag -d latest
git tag -a latest -m "Latest"
git push --force origin --tags
