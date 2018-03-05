#!/usr/bin/python
# -*- coding: UTF-8 -*-

from io import StringIO
from bypy import ByPy
import sys
import logging

if len(sys.argv) <= 1:
    logging.error("No filename specified")
    exit()

fn = sys.argv[1]
folder = "/"

if len(sys.argv) > 2:
    folder = sys.argv[2] + folder

bp=ByPy(verbose=1, debug=True)
bp.upload(fn, folder + fn)

sys.stdout = mystdout = StringIO()

bp.meta(folder + fn, '$i')

sys.stdout = sys.__stdout__

mystdout.seek(0)
mystdout.readline()
fid = mystdout.readline()

print("fid:\"" + str(fid).strip() + "\"")