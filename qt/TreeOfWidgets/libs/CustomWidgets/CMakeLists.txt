project(CustomWidgets)

set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)

find_package(Qt5 COMPONENTS Core Gui Widgets)

add_library(${PROJECT_NAME} STATIC "")

add_subdirectory(src)

target_link_libraries(${PROJECT_NAME}
    BinaryTreeModel

    Qt5::Core
    Qt5::Gui
    Qt5::Widgets
)


# install the shared library
install(TARGETS ${PROJECT_NAME}
    DESTINATION
        ${libDir}
    COMPONENT
        libs
)
