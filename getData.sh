#!/bin/bash

# For 720p
for img in dataset/720p/*.png; do
	./bin/png2go $img result/720p/$(basename $img .png) >> data720p.csv
done


# For 1080p
for img in dataset/1080p/*.png; do
	./bin/png2go $img result/1080p/$(basename $img .png) >> data1080p.csv
done
