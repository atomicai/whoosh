import geopandas
import folium
import pandas
from shapely.geometry import Polygon
import numpy
import math
from numpy.random import random
import csv

x_min = 37.5308 - 0.005
y_min = 55.7029 - 0.005

size = 0.0001
h = size / math.sqrt(2)

x_step = size + h
y_step = 2 * h

num = 100
x_max = x_min + x_step * num
y_max = y_min + y_step * num

rows = list(numpy.arange(y_min, y_max, y_step))
cols = list(numpy.arange(x_min, x_max, x_step))

polygons = []
for j, y in enumerate(rows):
  for i, x in enumerate(cols):
    if (i == 0):
      pass
    elif (i % 2 == 1):
        y = y + h
    elif (i % 2 == 0):
        y = y - h
    polygons.append(
      Polygon([
        (x, y),
        (x + h, y + h),
        (x + size + h, y + h),
        (x + size + h + h, y),
        (x + size + h, y - h),
        (x + h, y - h),
      ])
    )

graph = geopandas.GeoDataFrame({'geometry':polygons})
graph['centroid'] = graph.centroid

def setX(index, array):
  item = array.iloc[index]
  i = index % 100
  j = (index - i) // 100

  x = 1
  if ((i == 0) or (i == 99)):
    x = 1
  elif ((j == 0) or (j == 99)):
    x = 1
  else:
    x = random()

  return x

def getObject(index, array, xs):
  item = array.iloc[index]
  i = index % 100
  j = (index - i) //  100

  if ((i == 0) or (i == 99)):
    return [
      item.centroid.x,
      item.centroid.y,
      int(-1), 1,
      int(-1), 1,
      int(-1), 1,
      int(-1), 1,
      int(-1), 1,
      int(-1), 1,
    ]
  elif ((j == 0) or (j == 99)):
    return [
      item.centroid.x,
      item.centroid.y,
      int(-1), 1,
      int(-1), 1,
      int(-1), 1,
      int(-1), 1,
      int(-1), 1,
      int(-1), 1,
    ]

  # 1
  id1 = index - 100;
  w1 = max(xs[index], xs[id1]);

  # 2
  id2 = index - 100 + 1;
  w2 = max(xs[index], xs[id2]);

  # 3
  id3 = index + 1;
  w3 = max(xs[index], xs[id3]);

  # 4
  id4 = index + 100;
  w4 = max(xs[index], xs[id4]);

  # 5
  id5 = index - 1;
  w5 = max(xs[index], xs[id5]);

  # 6
  id6 = index -100 - 1;
  w6 = max(xs[index], xs[id6]);

  return [
    item.centroid.x,
    item.centroid.y,
    int(id1), w1,
    int(id2), w2,
    int(id3), w3,
    int(id4), w4,
    int(id5), w5,
    int(id6), w6,
  ]

parkings = []
with open("parkings.csv", 'r') as file:
  csvreader = csv.reader(file)
  for row in csvreader:
    row = row[2:]
    row[0] = float(row[0])
    row[1] = float(row[1])
    if 37.5258 <= row[0] <= 37.542 and 55.6977 <= row[1] <= 55.712:
      parkings.append(row)

xs = []
arr = []

ids = list(numpy.arange(0, num*num, 1))

for index, item in enumerate(ids):
  xs.append(setX(index, graph))

for index, item in enumerate(ids):
  arr.append(getObject(index, graph, xs))

cnt = 0
for park in parkings:
  mnDist = 100
  index = -1
  for i in range(len(arr)):
    x = arr[i][0]
    y = arr[i][1]
    dist = (x - park[0])**2 + (y - park[1])**2
    if dist <= 0.0001**2 and dist < mnDist:
      mnDist = dist
      index = i
  if index != -1 and len(arr[index]) < 15:
    arr[index].append(1)
    cnt+=1
for i in range(len(arr)):
  if len(arr[i]) < 15:
    arr[i].append(0)

# Import hexes
hexes = []
with open("hexes.csv", 'r') as file:
  csvreader = csv.reader(file)
  next(csvreader, None)
  for row in csvreader:
    row = row[1:]
    row[0] = int(float(row[0]))
    row[1] = float(row[1])
    row[2] = float(row[2])
    if x_min <= row[1] <= x_max and y_min <= row[2] <= y_max:
      hexes.append(row)

# Join hexes id with graph order
inds = []
for i in range(len(arr)):
  x = arr[i][0]
  y = arr[i][1]
  min_dist = 1
  min_ind = 0
  for hex in hexes:
    ind = hex[0]
    x_0 = hex[1]
    y_0 = hex[2]
    dist = (x - x_0)**2 + (y - y_0)**2
    if (dist < min_dist):
      min_ind = ind
      min_dist = dist
  inds.append(min_ind)

max_ind = numpy.max(inds)
min_ind = numpy.min(inds)

# Join road index to graph
input_file = r"road_index (1).csv"
with open(input_file, 'r') as file:
  csvreader = csv.reader(file)
  next(csvreader, None)
  for row in csvreader:
    row = row[0:]
    row[0] = int(row[0])
    row[1] = float(row[1])
    if (row[0] >= min_ind and row[0] <= max_ind): # Heuristic hack
      for i in range(len(arr)):
        if (inds[i] == row[0]):
          arr[i].append(row[1])
for i in range(len(arr)):
  if len(arr[i]) < 16:
    arr[i].append(1)

# Join routes to graph
input_file = r"routes_hex20m.csv"
with open(input_file, 'r') as file:
  csvreader = csv.reader(file)
  next(csvreader, None)
  for row in csvreader:
    row = row[0:]
    row[0] = float(row[0])
    row[1] = int(row[1])
    if (row[1] >= min_ind and row[1] <= max_ind): # Heuristic hack
      for i in range(len(arr)):
        if (inds[i] == row[1]):
          arr[i].append(row[0])
for i in range(len(arr)):
  if len(arr[i]) < 17:
    arr[i].append(0)

# Join median speed to graph
input_file = r"speed_median_hex20m_hackaton (1).csv"
with open(input_file, 'r') as file:
  csvreader = csv.reader(file)
  next(csvreader, None)
  for row in csvreader:
    row = row[0:]
    row[0] = float(row[0])
    row[1] = int(row[1])
    if (row[1] >= min_ind and row[1] <= max_ind): # Heuristic hack
      for i in range(len(arr)):
        if (inds[i] == row[1]):
          arr[i].append(row[0])
for i in range(len(arr)):
  if len(arr[i]) < 18:
    arr[i].append(0)

file = numpy.asarray(arr)
df = pandas.DataFrame(file)
df.to_csv('graph.csv')