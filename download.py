#!/usr/bin/python
# -*- coding: UTF-8 -*-

from youtube_dl import YoutubeDL
import sys
import logging

if len(sys.argv) <= 1:
    logging.error("No url specified")
    exit()

url = sys.argv[1]

print(url)

with YoutubeDL({'format': "bestvideo[ext=mp4]+bestaudio[ext=m4a]/bestvideo+bestaudio",  'outtmpl': '%(title)s.%(ext)s'}) as ydl:
    info_dict = ydl.extract_info(url, download=True)
    fn = ydl.prepare_filename(info_dict)

print("fn:\"" + fn + "\"")