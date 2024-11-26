cmake_minimum_required(VERSION 3.30)
project(fblthp_collection_manager)

set(CMAKE_CXX_STANDARD 23)

add_library(sqlite3 STATIC lib/sqlite3.c) # Adjust the path to your sqlite3.c
target_include_directories(sqlite3 PUBLIC lib) # Path to sqlite3.h

add_executable(fblthp_collection_manager main.cpp)
target_link_libraries(fblthp_collection_manager LINK_PUBLIC sqlite3)