#!/usr/bin/python
# -*- coding: UTF-8 -*-

from youtube_dl import YoutubeDL
from io import StringIO
from bypy import ByPy
import sys
import logging

if len(sys.argv) <= 1:
    logging.error("No url specified")
    exit()

url = sys.argv[1]
folder = "/"

if len(sys.argv) > 2:
    folder = sys.argv[2] + folder

print(url)

with YoutubeDL({format: "best", 'outtmpl': '%(title)s.%(ext)s'}) as ydl:
    info_dict = ydl.extract_info(url, download=True)
    fn = ydl.prepare_filename(info_dict)

print(fn)

bp=ByPy()
bp.verbose=True
bp.upload(fn, folder + fn)

sys.stdout = mystdout = StringIO()

bp.meta(folder + fn, '$i')

sys.stdout = sys.__stdout__

mystdout.seek(0)
mystdout.readline()
fid = mystdout.readline()

print("fid:\"" + str(fid).strip() + "\"")