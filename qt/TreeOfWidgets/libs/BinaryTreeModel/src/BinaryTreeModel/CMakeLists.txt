
target_sources(${PROJECT_NAME}
    PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/BinaryTreeModel.h
        ${CMAKE_CURRENT_SOURCE_DIR}/BinaryNode.h

    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/BinaryTreeModel.cpp
        ${CMAKE_CURRENT_SOURCE_DIR}/BinaryNode.cpp

        ${CMAKE_CURRENT_SOURCE_DIR}/CMakeLists.txt
)

