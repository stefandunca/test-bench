cmake_minimum_required(VERSION 3.19)

project(Chat)

set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)

find_package(Qt5 COMPONENTS Quick Qml REQUIRED)

add_library(${PROJECT_NAME} STATIC "")

add_subdirectory(src)
add_subdirectory(qml)

target_link_libraries(${PROJECT_NAME}
    PUBLIC
        Communities

    PRIVATE
        Qt5::Quick
        Qt5::Qml
)

# install the shared library
install(TARGETS ${PROJECT_NAME}
    DESTINATION
        ${libDir}
    COMPONENT
        libs
)
