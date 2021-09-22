cmake_minimum_required(VERSION 3.7)

set(VERSION 4.0)

configure_file(dde-dock.pc.in dde-dock.pc @ONLY)

project(dde-dock)

#set(CMAKE_VERBOSE_MAKEFILE ON)
set(CMAKE_CXX_STANDARD 14)
set(CMAKE_INCLUDE_CURRENT_DIR ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_CXX_FLAGS "-g -Wall")

if (${CMAKE_SYSTEM_PROCESSOR} MATCHES "mips64")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -O3 -ftree-vectorize -march=loongson3a -mhard-float -mno-micromips -mno-mips16 -flax-vector-conversions -mloongson-ext2 -mloongson-mmi -fPIE")
    set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -pie")
endif()

# generate a compile commands file as complete database for vim-YouCompleteMe or some other similar tools
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

include(GNUInstallDirs)

if (NOT (${CMAKE_BUILD_TYPE} MATCHES "Debug"))
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Ofast")

    # generate qm
    execute_process(COMMAND bash "translate_generation.sh"
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR})
endif ()

set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -DQT_DEBUG")

# Test architecture
execute_process(COMMAND dpkg-architecture -qDEB_BUILD_ARCH OUTPUT_VARIABLE ARCHITECTURE RESULT_VARIABLE EXIT_CODE)
if (${EXIT_CODE} EQUAL 0)
    string(STRIP ${ARCHITECTURE} ARCHITECTURE)

    if (${ARCHITECTURE} STREQUAL "sw_64")
        # add compiler flags -mieee for mathmatic
        set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -mieee")
        add_definitions(-DDISABLE_SHOW_ANIMATION)
    endif()
endif()

file(GLOB INTERFACES "interfaces/*.h")

#因为单元测试需要直接测试源代码，而主程序代码中include的单元使用了相对路径
#单元测试的CMakeLists和主程序的CMakeLists路径不同，编译单元测试时会提示找不到文件
#因此设置搜索路径
include_directories(
    frame/accessible
    frame/controller
    frame/dbus
    frame/dbus/sni
    frame/display
    frame/item
    frame/item/components
    frame/item/resources
    frame/util
    frame/window
    frame/window/components
    frame/xcb
    ../widgets
    ../interfaces
    )
aux_source_directory(frame/accessible ACCESSIBLE)
aux_source_directory(frame/controller CONTROLLER)
aux_source_directory(frame/dbus DBUS)
aux_source_directory(frame/dbus/sni SNI)
aux_source_directory(frame/display DISPLAY)
aux_source_directory(frame/item ITEM)
aux_source_directory(frame/item/components ITEMCOMPONENTS)
aux_source_directory(frame/item/resources RESOURCES)
aux_source_directory(frame/util UTIL)
aux_source_directory(frame/window WINDOW)
aux_source_directory(frame/window/components WINDOWCOMPONENTS)
aux_source_directory(frame/xcb XCB)

file(GLOB SRC_PATH
    ${ACCESSIBLE}
    ${CONTROLLER}
    ${DBUS}
    ${SNI}
    ${DISPLAY}
    ${ITEM}
    ${ITEMCOMPONENTS}
    ${UTIL}
    ${WINDOW}
    ${WINDOWCOMPONENTS}
    ${XCB}
    )

add_subdirectory("frame")
add_subdirectory("plugins")
add_subdirectory("tests")

# Install settings
if (CMAKE_INSTALL_PREFIX_INITIALIZED_TO_DEFAULT)
    set(CMAKE_INSTALL_PREFIX /usr)
endif ()

## qm files
file(GLOB QM_FILES "translations/*.qm")
install(FILES ${QM_FILES}
    DESTINATION share/dde-dock/translations)

## dev files
install(FILES ${INTERFACES}
    DESTINATION include/dde-dock)

install(FILES ${CMAKE_BINARY_DIR}/dde-dock.pc
    DESTINATION ${CMAKE_INSTALL_LIBDIR}/pkgconfig)

install(FILES "cmake/DdeDock/DdeDockConfig.cmake"
    DESTINATION ${CMAKE_INSTALL_LIBDIR}/cmake/DdeDock)

install(FILES gschema/com.deepin.dde.dock.module.gschema.xml
    DESTINATION share/glib-2.0/schemas)

# Address Sanitizer 内存错误检测工具,打开下面的编译选项可以看到调试信息，正常运行时不需要这些信息
#set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -g -fsanitize=address -O2")
#set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -g -fsanitize=address -O2")
