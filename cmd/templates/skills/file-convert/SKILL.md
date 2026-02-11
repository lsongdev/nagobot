---
name: file-convert
description: Convert between file formats (documents, images, audio, video).
tags: [convert, files, cross-platform]
---
# File Format Conversion

Convert between document, image, audio, and video formats using cross-platform CLI tools.

## Document Conversion (Pandoc)

Markdown to HTML:
```
exec: pandoc /path/to/input.md -o /path/to/output.html
```

Markdown to PDF (requires LaTeX):
```
exec: pandoc /path/to/input.md -o /path/to/output.pdf
```

Markdown to Word (.docx):
```
exec: pandoc /path/to/input.md -o /path/to/output.docx
```

Word to Markdown:
```
exec: pandoc /path/to/input.docx -o /path/to/output.md
```

HTML to Markdown:
```
exec: pandoc /path/to/input.html -t markdown -o /path/to/output.md
```

Word to PDF:
```
exec: pandoc /path/to/input.docx -o /path/to/output.pdf
```

With table of contents:
```
exec: pandoc /path/to/input.md --toc -o /path/to/output.pdf
```

With custom CSS (HTML output):
```
exec: pandoc /path/to/input.md -c /path/to/style.css --standalone -o /path/to/output.html
```

## Image Conversion (ImageMagick)

PNG to JPEG:
```
exec: convert /path/to/input.png /path/to/output.jpg
```

Resize image:
```
exec: convert /path/to/input.png -resize 800x600 /path/to/output.png
```

Resize keeping aspect ratio:
```
exec: convert /path/to/input.png -resize 800x /path/to/output.png
```

Batch convert:
```
exec: for f in /path/to/dir/*.png; do convert "$f" "${f%.png}.jpg"; done
```

Compress JPEG:
```
exec: convert /path/to/input.jpg -quality 75 /path/to/output.jpg
```

Create thumbnail:
```
exec: convert /path/to/input.png -thumbnail 200x200 /path/to/thumb.png
```

PDF to images:
```
exec: convert /path/to/input.pdf /path/to/output-%03d.png
```

Images to PDF:
```
exec: convert /path/to/img1.png /path/to/img2.png /path/to/output.pdf
```

Get image info:
```
exec: identify /path/to/image.png
```

## Image Conversion (sips — macOS built-in)

Convert format:
```
exec: sips -s format jpeg /path/to/input.png --out /path/to/output.jpg
```

Resize:
```
exec: sips --resampleWidth 800 /path/to/input.png --out /path/to/output.png
```

Get properties:
```
exec: sips -g pixelWidth -g pixelHeight -g format /path/to/image.png
```

## Audio/Video Conversion (FFmpeg)

Video format conversion:
```
exec: ffmpeg -i /path/to/input.mov -c:v libx264 -c:a aac /path/to/output.mp4
```

Extract audio from video:
```
exec: ffmpeg -i /path/to/input.mp4 -vn -acodec libmp3lame /path/to/output.mp3
```

Convert audio format:
```
exec: ffmpeg -i /path/to/input.wav /path/to/output.mp3
```

Video to GIF:
```
exec: ffmpeg -i /path/to/input.mp4 -vf "fps=10,scale=480:-1" /path/to/output.gif
```

Trim video:
```
exec: ffmpeg -i /path/to/input.mp4 -ss 00:00:30 -t 00:01:00 -c copy /path/to/output.mp4
```

Resize video:
```
exec: ffmpeg -i /path/to/input.mp4 -vf "scale=1280:720" /path/to/output.mp4
```

Get media info:
```
exec: ffprobe -v quiet -print_format json -show_format -show_streams /path/to/media.mp4
```

## Text Encoding

Convert encoding:
```
exec: iconv -f GB2312 -t UTF-8 /path/to/input.txt > /path/to/output.txt
```

Detect encoding:
```
exec: file -I /path/to/file.txt
```

## Base64

Encode:
```
exec: base64 /path/to/file > /path/to/file.b64
```

Decode:
```
exec: base64 -d /path/to/file.b64 > /path/to/file
```

## Notes

- **Pandoc**: Install with `brew install pandoc` or `apt install pandoc`. PDF output needs LaTeX (`brew install basictex`).
- **ImageMagick**: Install with `brew install imagemagick` or `apt install imagemagick`. On macOS, `sips` is built-in.
- **FFmpeg**: Install with `brew install ffmpeg` or `apt install ffmpeg`.
- `iconv` and `base64` are pre-installed on macOS and Linux.
- All tools work cross-platform (macOS, Linux, Windows WSL).
