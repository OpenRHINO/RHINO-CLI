# MPI Template
```
.
├── Dockerfile
├── ldd.sh
├── README.md
└── src
    ├── main.cpp
    └── Makefile
```
## Dockerfile
Use multi-stage compilation to generate runtime image
## ldd.sh
Analyze dynamic dependence, remove soft connections and pack libs
## main.cpp
Main function with MPI basic constructs
## Makefile
Makefile template to build cpp functions
> Note: If you copy your MPI project to `/src`, just make sure that:
> 1. Use `-f` to sepcify relative path of the Makefile, e.g. `rhino build -f ./src/conf/linux.makefile`
> 2. Modify the name of the target file to `mpi-func`, e.g. `$(EXEC): $(OBJS) $(CXX) -o mpi-func`

## Notice

**For developers of this project:**

When updates are made to these templates, they must be packaged (serialized as a Go byte array) by running `make generate`. Then check the resultant file `./generate/zz_filesystem_generated.go` file.
