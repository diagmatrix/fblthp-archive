# Compiler and flags
CXX = g++
CXXFLAGS = -std=c++23 -Wall -Wextra -Wpedantic -Wshadow -Wunused -Wnull-dereference -Wformat=2
LDFLAGS = -lsqlite3 -lcurl

# Directories
SRC_DIR_DOORKEEPER = src/migration-manager
SRC_DIR_FBLTHP = src/fblthp
INCLUDE_DIR = lib
BIN_DIR = bin
OBJ_DIR = obj
OBJ_DIR_DOORKEEPER = $(OBJ_DIR)/doorkeeper
OBJ_DIR_FBLTHP = $(OBJ_DIR)/fblthp

# Files
SOURCES_DOORKEEPER = $(SRC_DIR_DOORKEEPER)/main.cpp \
                     $(SRC_DIR_DOORKEEPER)/migrations.cpp
SOURCES_FBLTHP = $(SRC_DIR_FBLTHP)/main.cpp
OBJECTS_DOORKEEPER = $(addprefix $(OBJ_DIR_DOORKEEPER)/, $(notdir $(SOURCES_DOORKEEPER:.cpp=.o)))
OBJECTS_FBLTHP = $(addprefix $(OBJ_DIR_FBLTHP)/, $(notdir $(SOURCES_FBLTHP:.cpp=.o)))
TARGET_DOORKEEPER = $(BIN_DIR)/doorkeeper
TARGET_FBLTHP = $(BIN_DIR)/fblthp
JSON_URL = https://raw.githubusercontent.com/nlohmann/json/refs/tags/v3.11.3/single_include/nlohmann/json.hpp
JSON_HEADER = $(INCLUDE_DIR)/nlohmann/json.hpp

# Build rules
all: doorkeeper fblthp

doorkeeper: $(TARGET_DOORKEEPER)

fblthp: fetch-json $(TARGET_FBLTHP)

$(TARGET_DOORKEEPER): $(OBJECTS_DOORKEEPER)
	@mkdir -p $(BIN_DIR)
	$(CXX) $(CXXFLAGS) -o $@ $^ $(LDFLAGS)

$(TARGET_FBLTHP): $(OBJECTS_FBLTHP)
	@mkdir -p $(BIN_DIR)
	$(CXX) $(CXXFLAGS) -o $@ $^ $(LDFLAGS)

$(OBJ_DIR_DOORKEEPER)/%.o: $(SRC_DIR_DOORKEEPER)/%.cpp
	@mkdir -p $(OBJ_DIR_DOORKEEPER)
	$(CXX) $(CXXFLAGS) -I$(INCLUDE_DIR) -c $< -o $@

$(OBJ_DIR_FBLTHP)/%.o: $(SRC_DIR_FBLTHP)/%.cpp
	@mkdir -p $(OBJ_DIR_FBLTHP)
	$(CXX) $(CXXFLAGS) -I$(INCLUDE_DIR) -c $< -o $@

$(JSON_HEADER):
	@mkdir -p $(INCLUDE_DIR)/nlohmann
	curl -L $(JSON_URL) -o $(JSON_HEADER)

fetch-json: $(JSON_HEADER)

clean:
	rm -rf $(OBJ_DIR) $(BIN_DIR) $(INCLUDE_DIR)/nlohmann

.PHONY: all clean fetch-json