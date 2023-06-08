import geopandas
import numpy
import pandas

# Reading a GeoPackage file
input_file = r"hexes (1).gpkg"
data = geopandas.read_file(input_file)
data.head()

arr = []
for i in range(len(data)):
    id = data.id[i]
    geom = data.geometry[i]
    arr.append([int(id), geom.centroid.x, geom.centroid.y])

file = numpy.asarray(arr)
df = pandas.DataFrame(file)
df.to_csv('hexes.csv')
