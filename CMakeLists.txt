cmake_minimum_required(VERSION 3.30)
project(fblthp_collection_manager)

set(CMAKE_CXX_STANDARD 23)

add_library(sqlite3 STATIC lib/sqlite3.c) # Adjust the path to your sqlite3.c
target_include_directories(sqlite3 PUBLIC lib) # Path to sqlite3.h

# Migration manager
add_executable(doorkeeper src/migration-manager/main.cpp
        src/migration-manager/migrations.h
        src/migration-manager/migrations.cpp
        src/migration-manager/main.cpp)
target_link_libraries(doorkeeper LINK_PUBLIC sqlite3)