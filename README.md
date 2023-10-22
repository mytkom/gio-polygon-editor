# gio-polygon-editor
Polygon Editor written in Go using Gio library (University Project)

# Installation
TODO

# Functional Keys
`N` (*N*ew) - enables polygon builder, clicking on canvas would create a vertex for new polygon, clicking on the first one created closes it and disables builder. 
`D` (*D*elete) - if polygon/edge/vertex is selected it removes it.

### When edge is selected:
`V` (*V*ertical) - sets a constraint on edge, so it have to be vertical
`H` (*H*orizontal) - sets a constraint on edge, so it have to be horizontal
`C` (*C*lear) - unsets constraints on edge

`A` (*A*dd vertex) - adds vertex in the middle of the selected edge

### Important note
Constraints cannot be set to any two neighbour edges. Setting constraint wouldn't work. Additionally, because of this - addition of vertex in the middle of the edge would disable constraint previously set on the edge.
