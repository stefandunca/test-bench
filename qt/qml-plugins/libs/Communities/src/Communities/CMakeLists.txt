target_sources(${PROJECT_NAME}
    PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/CommunitiesModule.h

    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/CommunitiesPlugin.h
        ${CMAKE_CURRENT_SOURCE_DIR}/CommunitiesPlugin.cpp
        ${CMAKE_CURRENT_SOURCE_DIR}/CommunitiesModule.cpp

        ${CMAKE_CURRENT_SOURCE_DIR}/CMakeLists.txt
)
