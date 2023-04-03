# Compiler and linker options
CXX = mpicxx
CXXFLAGS = -fopenmp -Wextra -pedantic -Wall -O3
MPI_LFLAGS =
LDFLAGS = -pthread

# Source files and object files
SRCS = main.cpp
OBJS = $(patsubst %.cpp,%.o,$(SRCS))

# Libraries and header file paths
LIBS = 
INCLUDES = 

# Target executable
TARGET = mpi-func

# Phony targets
.PHONY: all clean

# Build rules
all: clean $(TARGET)

$(TARGET): $(OBJS)
	$(CXX) $(MPI_LFLAGS) $(OBJS) $(LIBS) -o $(TARGET) $(LDFLAGS)

%.o: %.cpp
	$(CXX) $(CXXFLAGS) $(INCLUDES) -c $< -o $@

clean:
	rm -f $(OBJS) $(TARGET)

.DEFAULT_GOAL := all