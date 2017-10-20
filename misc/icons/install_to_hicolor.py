#!/usr/bin/python3
import sys
import os
import os.path
import argparse
import subprocess
from shutil import copyfile

available_sizes = '16,22,24,32,48,64,96,128'

parser = argparse.ArgumentParser()
parser.add_argument('target', help='input svg file or directory')

default_output = 'icons'
default_dir = 'apps'

parser.add_argument('-d', '--directory',
        dest='directory',
        default=default_dir,
        help='the directory to export icon, default: %s' % default_dir
        )

parser.add_argument('-o', '--output',
        dest='output',
        default= default_output,
        help=('the output file path, default: %s') % default_output
        )

args = parser.parse_args()
target = args.target
output = args.output
directory = args.directory
sizes=available_sizes.split(',')

def get_target_files():
    files = []

    if os.path.exists(target) == False:
        raise Exception('The input target %s not found' % target)
        return files

    if os.path.isfile(target):
        files.append(target)
        return files

    for f in os.listdir(target):
        files.append(os.path.join(target,f))

    return files

def svg2png(svg_file):
    name = os.path.splitext(os.path.basename(svg_file))[0]
    for size in sizes:
        png_file =  os.path.join(output, "hicolor", size + "x" + size, directory, name+".png")
        print('Will convert file %s to %s' % (svg_file, png_file))
        os.makedirs( os.path.dirname(png_file), mode=0o755, exist_ok=True)
        subprocess.run(['rsvg-convert', '-w', str(size), '-h', str(size), '-o', png_file, svg_file ])

def copy_svg_file(svg_file):
    target_file = os.path.join(output, "hicolor", "scalable", directory, os.path.basename(svg_file))
    print('Will copy file %s to %s' % (svg_file, target_file))
    os.makedirs( os.path.dirname(target_file), mode=0o755, exist_ok=True)
    copyfile(svg_file, target_file)

if __name__ == '__main__':
    files = get_target_files()
    for f in files:
        svg2png(f)
        copy_svg_file(f)
