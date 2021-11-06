

target_include_directories(${PROJECT_NAME}
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/BinaryTreeModel

    PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}
)

add_subdirectory(BinaryTreeModel)
