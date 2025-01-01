# Compiler and flags
CXX = g++
CXXFLAGS = -std=c++23 -Wall -Wextra -Wpedantic -Wshadow -Wunused -Wformat=2
LDFLAGS = -lsqlite3

# Directories
SRC_DIR = src/migration-manager
INCLUDE_DIR = src
BIN_DIR = bin
OBJ_DIR = obj

# Files
SOURCES = $(SRC_DIR)/main.cpp \
          $(SRC_DIR)/migrations.cpp
OBJECTS = $(addprefix $(OBJ_DIR)/, $(notdir $(SOURCES:.cpp=.o)))
TARGET = $(BIN_DIR)/doorkeeper

# Build rules
all: $(TARGET)

$(TARGET): $(OBJECTS)
	@mkdir -p $(BIN_DIR)
	$(CXX) $(CXXFLAGS) -o $@ $^ $(LDFLAGS)

$(OBJ_DIR)/%.o: $(SRC_DIR)/%.cpp
	@mkdir -p $(OBJ_DIR)
	$(CXX) $(CXXFLAGS) -I$(INCLUDE_DIR) -c $< -o $@

clean:
	rm -rf $(OBJ_DIR) $(BIN_DIR)

.PHONY: all clean
