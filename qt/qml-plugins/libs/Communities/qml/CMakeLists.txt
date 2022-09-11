# The variable is cached in ${BUILD_DIR}/CMakeCache.txt. We need FORCE to
# change it there immediately. Also, add a comment to the cache file.
set(QML_IMPORT_PATH ${CMAKE_CURRENT_SOURCE_DIR};${QML_IMPORT_PATH} CACHE STRING "Extra qml import paths for QtCreator" FORCE)
list(REMOVE_DUPLICATES QML_IMPORT_PATH)

target_sources(${PROJECT_NAME}
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/CMakeLists.txt

        ${CMAKE_CURRENT_SOURCE_DIR}/Communities-qml.qrc
)
