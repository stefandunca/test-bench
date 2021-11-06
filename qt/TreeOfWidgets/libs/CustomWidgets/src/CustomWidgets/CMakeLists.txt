
target_sources(${PROJECT_NAME}
    PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/BinaryTreeView.h
        ${CMAKE_CURRENT_SOURCE_DIR}/BaseNodeWidget.h
        ${CMAKE_CURRENT_SOURCE_DIR}/StringNodeWidget.h
        ${CMAKE_CURRENT_SOURCE_DIR}/FloatNodeWidget.h

    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/BinaryTreeView.cpp
        ${CMAKE_CURRENT_SOURCE_DIR}/BaseNodeWidget.cpp
        ${CMAKE_CURRENT_SOURCE_DIR}/StringNodeWidget.cpp
        ${CMAKE_CURRENT_SOURCE_DIR}/FloatNodeWidget.cpp

        ${CMAKE_CURRENT_SOURCE_DIR}/CMakeLists.txt
)

