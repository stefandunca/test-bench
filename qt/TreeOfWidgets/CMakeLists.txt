# global minimum version
cmake_minimum_required(VERSION 3.19)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

set(projectname TreeOfWidgets)
project(${projectname})

add_subdirectory(apps)
add_subdirectory(libs)


set (incDir ${CMAKE_INSTALL_PREFIX}/inc)
set (binDir ${CMAKE_INSTALL_PREFIX}/bin)
set (libDir ${CMAKE_INSTALL_PREFIX}/lib${LIB_SUFFIX})