for img in ../1080p/*.png; do
    ffmpeg -i "$img" -vf "scale=1280:720" "./$(basename "$img")"
done
