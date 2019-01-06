"""Stitches a directory of images together into one image."""
import argparse
import os
import logging
import math
import sys
from PIL import Image

def Stitch(image_dir='images', output_image='output.jpg', images_per_row=6):
  images = [Image.open(os.path.join(image_dir, img))
            for img in os.listdir(image_dir)
            if img != '.DS_Store'
  ]
  if not images_per_row:
    images_per_row = int(math.floor(math.sqrt(len(images))))

  image_width = images[0].size[0]
  image_height = images[0].size[1]
  logging.info('Using image_width, image_height = (%d,%d)' % (image_width,
                                                              image_height))
  
  width = images_per_row * image_width
  height = (len(images) / images_per_row) * image_height
  logging.info('New image width, height = (%d,%d)' % (width,height))
  new_im = Image.new('RGB', (width, height))
  
  for i, im in enumerate(images):
    row = i / images_per_row
    col = i % images_per_row
    x = col * image_width
    y = row * image_height
    logging.info('row=%d col=%d x=%d y=%d', row, col, x ,y)
    new_im.paste(im, (x, y))
  new_im.save(output_image)
  logging.info('Output to %s', output_image)

def main(argv):
  logging.basicConfig(level=logging.INFO)
  parser = argparse.ArgumentParser()
  parser.add_argument('--images', help='Directory containing images')
  parser.add_argument('--output', help='Path of the output file')
  parser.add_argument('--images_per_row', default=0, nargs='?', type=int,
                      help=('Set this to override the number of images per row. '
                            'If not specified it defaults to '
                            'floor(sqrt(len(images)))'))
  args = parser.parse_args()
  Stitch(args.images, args.output, args.images_per_row)

if __name__ == '__main__':
  main(sys.argv)
