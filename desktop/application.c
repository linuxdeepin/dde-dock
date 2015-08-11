#include <gio/gio.h>

int content_type_can_be_executable(char* type)
{
    return g_content_type_can_be_executable(type);
}


int content_type_is(char* type, char* expected_type)
{
    return g_content_type_is_a(type, expected_type);
}

