# Polygon Editor
Polygon Editor written in Go using Gio library (University Project)

# Installation
### Development
1. (Download and Install Golang)[https://go.dev/doc/install]
2. execute `go install gioui.org/cmd/gogio@latest`.
3. execute `go run .`, in workdir.

If it is not working after these steps, follow (official Gio installation guide)[https://gioui.org/doc/install], you may be missing some dependencies, needed for your OS.

# Functional Keys
- `N` (*N*ew) - enables polygon builder, clicking on canvas would create a vertex for new polygon, clicking on the first one created closes it and disables builder. 
- `D` (*D*elete) - if polygon/edge/vertex is selected it removes it.

### Offseted Polygon
- `O` (*O*ffset) - toggles offsetted polygon
- `+` - make offset width bigger
- `-` - make offset width smaller

### When edge is selected:
- `V` (*V*ertical) - sets a constraint on edge, so it have to be vertical
- `H` (*H*orizontal) - sets a constraint on edge, so it have to be horizontal
- `C` (*C*lear) - unsets constraints on edge

`A` (*A*dd vertex) - adds vertex in the middle of the selected edge

### Saving/Loading one scene
- `S` (*S*ave) - saves current scene to "scene.json" file
- `L` (*L*oad) - loads scene from "scene.jsin" file

### Important note
The same constraints cannot be set to any two neighbour edges. Setting constraint wouldn't work. Additionally, because of this - addition of vertex in the middle of the edge would disable constraint previously set on the edge.
