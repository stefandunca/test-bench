target_include_directories(${PROJECT_NAME}
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/Communities

    PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}
)

add_subdirectory(Communities)