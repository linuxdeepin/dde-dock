cmake_minimum_required(VERSION 3.7)

set(VERSION 4.0)

project(dde-dock)

#set(CMAKE_VERBOSE_MAKEFILE ON)
set(CMAKE_CXX_STANDARD 14)
set(CMAKE_INCLUDE_CURRENT_DIR ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_CXX_FLAGS -Wall)

if (NOT (${CMAKE_BUILD_TYPE} MATCHES "Debug"))
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -O3")

    # generate qm
    execute_process(COMMAND bash "translate_generation.sh"
                    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})
endif ()

file(GLOB INTERFACES "interfaces/*.h")

add_subdirectory("frame")
add_subdirectory("plugins")

# Install settings
if (CMAKE_INSTALL_PREFIX_INITIALIZED_TO_DEFAULT)
    set(CMAKE_INSTALL_PREFIX /usr)
endif ()

## qm files
file(GLOB QM_FILES "translations/*.qm")
install(FILES ${QM_FILES}
        DESTINATION share/dde-dock/translations)

## dev files
install (FILES ${INTERFACES}
         DESTINATION include/dde-dock)

## services files
install(FILES frame/com.deepin.dde.Dock.service
     DESTINATION /usr/share/dbus-1/services)
